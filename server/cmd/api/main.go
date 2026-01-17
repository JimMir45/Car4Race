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

	// 初始化服务层
	userService := service.NewUserService(userRepo, cfg.JWTSecret)

	// 初始化处理器
	userHandler := handler.NewUserHandler(userService)

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
