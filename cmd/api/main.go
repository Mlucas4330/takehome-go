package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "takehome-go/docs"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"takehome-go/internal/config"
	"takehome-go/internal/database"
	"takehome-go/internal/handler"
	"takehome-go/internal/repository"
	"takehome-go/internal/service"
)

// @title Takehome-go API
// @version 1.0
// @description API REST para gerenciar Colaboradores e Departamentos
// @host localhost:8080
// @BasePath /api/v1
func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load .env variables", zap.Error(err))
	}

	postgresDsn := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresHost, cfg.PostgresDb)

	db, err := database.Connect(postgresDsn)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	redisAddr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)

	cache := database.NewRedisCache(redisAddr)

	colaboradorRepo := repository.NewColaboradorRepository(db)
	departamentoRepo := repository.NewDepartamentoRepository(db)

	colaboradorSvc := service.NewColaboradorService(colaboradorRepo, departamentoRepo, cache, logger)
	departamentoSvc := service.NewDepartamentoService(departamentoRepo, colaboradorRepo, cache, logger)

	colaboradorHandler := handler.NewColaboradorHandler(colaboradorSvc, logger)
	departamentoHandler := handler.NewDepartamentoHandler(departamentoSvc, logger)

	router := setupRouter(colaboradorHandler, departamentoHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	logger.Info("Server started successfully", zap.String("port", cfg.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited gracefully")
}

func setupRouter(colaboradorHandler *handler.ColaboradorHandler, departamentoHandler *handler.DepartamentoHandler) *gin.Engine {
	router := gin.Default()

	router.Use(handler.PrometheusMiddleware())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		colaboradores := v1.Group("/colaboradores")
		{
			colaboradores.POST("", colaboradorHandler.Create)
			colaboradores.GET("/:id", colaboradorHandler.GetByID)
			colaboradores.PUT("/:id", colaboradorHandler.Update)
			colaboradores.DELETE("/:id", colaboradorHandler.Delete)
			colaboradores.POST("/listar", colaboradorHandler.List)
		}

		departamentos := v1.Group("/departamentos")
		{
			departamentos.POST("", departamentoHandler.Create)
			departamentos.GET("/:id", departamentoHandler.GetByID)
			departamentos.PUT("/:id", departamentoHandler.Update)
			departamentos.DELETE("/:id", departamentoHandler.Delete)
			departamentos.POST("/listar", departamentoHandler.List)
		}

		gerentes := v1.Group("/gerentes")
		{
			gerentes.GET("/:id/colaboradores", departamentoHandler.GetColaboradoresByGerente)
		}
	}

	return router
}
