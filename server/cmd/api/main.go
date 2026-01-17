package main

import (
	"log"
	"os"

	"car4race/internal/config"
	"car4race/internal/handler"
	"car4race/internal/middleware"
	"car4race/internal/repository"
	"car4race/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	db, err := repository.InitDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	// 初始化仓储层
	userRepo := repository.NewUserRepository(db)
	contentRepo := repository.NewContentRepository(db)
	courseRepo := repository.NewCourseRepository(db)

	// 初始化服务层
	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	contentService := service.NewContentService(contentRepo)
	courseService := service.NewCourseService(courseRepo, userRepo)

	// 初始化处理器
	userHandler := handler.NewUserHandler(userService)
	contentHandler := handler.NewContentHandler(contentService)
	courseHandler := handler.NewCourseHandler(courseService)
	adminHandler := handler.NewAdminHandler(contentService, courseService)

	// 设置 Gin 模式
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	r := gin.Default()

	// 全局中间件
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API 路由组
	api := r.Group("/api/v1")
	{
		// 用户认证路由（无需登录）
		auth := api.Group("/auth")
		{
			auth.POST("/send-code", userHandler.SendCode)
			auth.POST("/login", userHandler.Login)
		}

		// 需要登录的路由
		protected := api.Group("")
		protected.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			protected.GET("/user/profile", userHandler.GetProfile)
			protected.PUT("/user/profile", userHandler.UpdateProfile)
		}

		// ========== 私域视频网站 (HPA) ==========
		hpa := api.Group("/hpa")
		{
			// 公开接口
			hpa.GET("/categories", contentHandler.GetCategories)
			hpa.GET("/notes", contentHandler.GetNotes)
			hpa.GET("/notes/:slug", contentHandler.GetNote)
			hpa.GET("/courses", courseHandler.GetCourses)
			hpa.GET("/courses/:slug", middleware.OptionalJWTAuth(cfg.JWTSecret), courseHandler.GetCourse)

			// 需要登录
			hpaAuth := hpa.Group("")
			hpaAuth.Use(middleware.JWTAuth(cfg.JWTSecret))
			{
				hpaAuth.GET("/history", contentHandler.GetBrowseHistory)
				hpaAuth.POST("/orders", courseHandler.CreateOrder)
				hpaAuth.GET("/orders", courseHandler.GetOrders)
				hpaAuth.POST("/redeem", courseHandler.RedeemCode)
				hpaAuth.POST("/download", courseHandler.CreateDownload)
				hpaAuth.GET("/download/:token", courseHandler.Download)
			}
		}

		// ========== 管理后台 ==========
		admin := api.Group("/admin")
		admin.Use(middleware.JWTAuth(cfg.JWTSecret))
		admin.Use(middleware.AdminAuth())
		{
			// 分类管理
			admin.POST("/categories", adminHandler.CreateCategory)
			admin.PUT("/categories/:id", adminHandler.UpdateCategory)
			admin.DELETE("/categories/:id", adminHandler.DeleteCategory)

			// 笔记管理
			admin.POST("/notes", adminHandler.CreateNote)
			admin.PUT("/notes/:id", adminHandler.UpdateNote)
			admin.DELETE("/notes/:id", adminHandler.DeleteNote)

			// 课程管理
			admin.GET("/courses", adminHandler.GetCourses)
			admin.POST("/courses", adminHandler.CreateCourse)
			admin.PUT("/courses/:id", adminHandler.UpdateCourse)
			admin.DELETE("/courses/:id", adminHandler.DeleteCourse)

			// 邀请码管理
			admin.GET("/invite-codes", adminHandler.GetInviteCodes)
			admin.POST("/invite-codes", adminHandler.CreateInviteCode)
		}
	}

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
	}
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
