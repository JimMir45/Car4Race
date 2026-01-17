package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"car4race/internal/model"
	"car4race/internal/repository"
)

type FileService struct {
	courseRepo *repository.CourseRepository
	uploadPath string
}

func NewFileService(courseRepo *repository.CourseRepository, uploadPath string) *FileService {
	return &FileService{
		courseRepo: courseRepo,
		uploadPath: uploadPath,
	}
}

// UploadCourseFile 上传课程文件
func (s *FileService) UploadCourseFile(courseID uint, fileType string, file *multipart.FileHeader) (*model.CourseFile, error) {
	// 验证课程是否存在
	course, err := s.courseRepo.GetCourseByID(courseID)
	if err != nil {
		return nil, fmt.Errorf("课程不存在")
	}

	// 创建上传目录 uploads/courses/{course_id}/
	courseDir := filepath.Join(s.uploadPath, "courses", fmt.Sprintf("%d", courseID))
	if err := os.MkdirAll(courseDir, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %v", err)
	}

	// 生成文件名（保留原始扩展名）
	ext := filepath.Ext(file.Filename)
	baseName := strings.TrimSuffix(file.Filename, ext)
	timestamp := time.Now().UnixNano() / 1e6
	newFileName := fmt.Sprintf("%s_%d%s", baseName, timestamp, ext)
	filePath := filepath.Join(courseDir, newFileName)

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %v", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("保存文件失败: %v", err)
	}

	// 获取当前文件的最大排序号
	maxSort, _ := s.courseRepo.GetMaxCourseFileSort(courseID)

	// 创建文件记录
	courseFile := &model.CourseFile{
		CourseID: courseID,
		FileType: fileType, // intro | resource
		FileName: file.Filename,
		FilePath: filePath,
		FileSize: file.Size,
		Sort:     maxSort + 1,
	}

	if err := s.courseRepo.CreateCourseFile(courseFile); err != nil {
		// 如果数据库操作失败，删除已上传的文件
		os.Remove(filePath)
		return nil, fmt.Errorf("保存记录失败: %v", err)
	}

	// 如果是 intro 类型，更新课程的 IntroPath
	if fileType == "intro" {
		course.IntroPath = filePath
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

	// 删除物理文件
	if err := os.Remove(file.FilePath); err != nil && !os.IsNotExist(err) {
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

	content, err := os.ReadFile(course.IntroPath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}

	return string(content), nil
}

// GetFilePath 获取文件路径（用于下载）
func (s *FileService) GetFilePath(fileID uint) (string, string, error) {
	file, err := s.courseRepo.GetCourseFileByID(fileID)
	if err != nil {
		return "", "", fmt.Errorf("文件不存在")
	}

	return file.FilePath, file.FileName, nil
}
