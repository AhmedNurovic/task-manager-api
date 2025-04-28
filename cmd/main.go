package main

import (
	"os"
	"os/signal"
	"syscall"
	"net/http"

	"github.com/ahmednurovic/task-manager-api/internal/config"
	"github.com/ahmednurovic/task-manager-api/internal/handler"
	"github.com/ahmednurovic/task-manager-api/internal/middleware"
	"github.com/ahmednurovic/task-manager-api/internal/repository"
	"github.com/ahmednurovic/task-manager-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	_ "github.com/ahmednurovic/task-manager-api/docs"
	swaggerFiles "github.com/swaggo/files" 
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Task Manager API
// @version 1.0
// @description This is a task management API with JWT authentication
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@taskmanager.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
		os.Exit(1)
	}

	logger.Info("Config loaded successfully", zap.String("db_url", cfg.DBURL), zap.String("port", cfg.Port))

	db, err := sqlx.Connect("postgres", cfg.DBURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	logger.Info("Successfully connected to the database")

	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	taskService := service.NewTaskService(taskRepo)
	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)

	taskHandler := handler.NewTaskHandler(taskService)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.ZapLogger(logger))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", handler.Register(authService))
			auth.POST("/login", handler.Login(authService))
		}

		tasks := api.Group("/tasks").Use(authMiddleware)
		{
			tasks.POST("", taskHandler.CreateTask)
			tasks.GET("", taskHandler.GetTasks)
			tasks.PUT("/:id", taskHandler.UpdateTask)
			tasks.DELETE("/:id", taskHandler.DeleteTask)
		}
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")
}
