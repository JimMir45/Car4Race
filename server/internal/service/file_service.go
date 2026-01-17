package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"car4race/internal/config"
	"car4race/internal/model"
	"car4race/internal/repository"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FileService struct {
	courseRepo  *repository.CourseRepository
	minioClient *minio.Client
	bucket      string
}

func NewFileService(courseRepo *repository.CourseRepository, cfg *config.Config) (*FileService, error) {
	// 初始化 MinIO 客户端
	minioClient, err := minio.New(cfg.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOSecretKey, ""),
		Secure: cfg.MinIOUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %v", err)
	}

	// 确保 bucket 存在
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, cfg.MinIOBucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %v", err)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, cfg.MinIOBucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %v", err)
		}
	}

	return &FileService{
		courseRepo:  courseRepo,
		minioClient: minioClient,
		bucket:      cfg.MinIOBucket,
	}, nil
}

// UploadCourseFile 上传课程文件到 MinIO
func (s *FileService) UploadCourseFile(courseID uint, fileType string, file *multipart.FileHeader) (*model.CourseFile, error) {
	// 验证课程是否存在
	course, err := s.courseRepo.GetCourseByID(courseID)
	if err != nil {
		return nil, fmt.Errorf("课程不存在")
	}

	// 生成对象路径: courses/{course_id}/{file_type}/{timestamp}_{filename}
	ext := filepath.Ext(file.Filename)
	baseName := strings.TrimSuffix(file.Filename, ext)
	timestamp := time.Now().UnixNano() / 1e6
	objectName := fmt.Sprintf("courses/%d/%s/%s_%d%s", courseID, fileType, baseName, timestamp, ext)

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer src.Close()

	// 设置 Content-Type
	contentType := "application/octet-stream"
	if ext == ".md" {
		contentType = "text/markdown; charset=utf-8"
	} else if ext == ".zip" {
		contentType = "application/zip"
	}

	// 上传到 MinIO
	ctx := context.Background()
	_, err = s.minioClient.PutObject(ctx, s.bucket, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, fmt.Errorf("上传文件失败: %v", err)
	}

	// 获取当前文件的最大排序号
	maxSort, _ := s.courseRepo.GetMaxCourseFileSort(courseID)

	// 创建文件记录
	courseFile := &model.CourseFile{
		CourseID: courseID,
		FileType: fileType, // intro | resource
		FileName: file.Filename,
		FilePath: objectName, // 存储 MinIO 对象路径
		FileSize: file.Size,
		Sort:     maxSort + 1,
	}

	if err := s.courseRepo.CreateCourseFile(courseFile); err != nil {
		// 如果数据库操作失败，删除已上传的文件
		s.minioClient.RemoveObject(ctx, s.bucket, objectName, minio.RemoveObjectOptions{})
		return nil, fmt.Errorf("保存记录失败: %v", err)
	}

	// 如果是 intro 类型，更新课程的 IntroPath
	if fileType == "intro" {
		course.IntroPath = objectName
		s.courseRepo.UpdateCourse(course)
	}

	return courseFile, nil
}

// DeleteCourseFile 删除课程文件
func (s *FileService) DeleteCourseFile(fileID uint) error {
	// 获取文件记录
	file, err := s.courseRepo.GetCourseFileByID(fileID)
	if err != nil {
		return fmt.Errorf("文件不存在")
	}

	// 从 MinIO 删除文件
	ctx := context.Background()
	if err := s.minioClient.RemoveObject(ctx, s.bucket, file.FilePath, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}

	// 删除数据库记录
	if err := s.courseRepo.DeleteCourseFile(fileID); err != nil {
		return fmt.Errorf("删除记录失败: %v", err)
	}

	return nil
}

// GetCourseFiles 获取课程文件列表
func (s *FileService) GetCourseFiles(courseID uint) ([]model.CourseFile, error) {
	return s.courseRepo.GetCourseFiles(courseID)
}

// GetCourseIntroContent 获取课程介绍 Markdown 内容
func (s *FileService) GetCourseIntroContent(courseID uint) (string, error) {
	course, err := s.courseRepo.GetCourseByID(courseID)
	if err != nil {
		return "", fmt.Errorf("课程不存在")
	}

	if course.IntroPath == "" {
		return "", nil
	}

	// 从 MinIO 读取文件内容
	ctx := context.Background()
	obj, err := s.minioClient.GetObject(ctx, s.bucket, course.IntroPath, minio.GetObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("获取文件失败: %v", err)
	}
	defer obj.Close()

	content, err := io.ReadAll(obj)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}

	return string(content), nil
}

// GetFileObject 获取文件对象（用于下载）
func (s *FileService) GetFileObject(fileID uint) (io.ReadCloser, string, int64, error) {
	file, err := s.courseRepo.GetCourseFileByID(fileID)
	if err != nil {
		return nil, "", 0, fmt.Errorf("文件不存在")
	}

	// 从 MinIO 获取文件
	ctx := context.Background()
	obj, err := s.minioClient.GetObject(ctx, s.bucket, file.FilePath, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", 0, fmt.Errorf("获取文件失败: %v", err)
	}

	return obj, file.FileName, file.FileSize, nil
}

// GetPresignedURL 获取预签名 URL（用于前端直接下载）
func (s *FileService) GetPresignedURL(fileID uint, expiry time.Duration) (string, error) {
	file, err := s.courseRepo.GetCourseFileByID(fileID)
	if err != nil {
		return "", fmt.Errorf("文件不存在")
	}

	ctx := context.Background()
	url, err := s.minioClient.PresignedGetObject(ctx, s.bucket, file.FilePath, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("生成下载链接失败: %v", err)
	}

	return url.String(), nil
}
