package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/mlucas4330/takehome-go/internal/database"
	"github.com/mlucas4330/takehome-go/internal/handlers"
	"github.com/mlucas4330/takehome-go/internal/repositories"
	"github.com/mlucas4330/takehome-go/internal/services"

	_ "github.com/mlucas4330/takehome-go/docs"
)

// @title API de Colaboradores e Departamentos
// @version 1.0
// @description API REST para gerenciar colaboradores e departamentos
// @host localhost:8080
// @BasePath /
func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	colabRepo := repositories.NewColaboradorRepository(db)
	deptRepo := repositories.NewDepartamentoRepository(db)

	colabService := services.NewColaboradorService(colabRepo, deptRepo)
	deptService := services.NewDepartamentoService(deptRepo, colabRepo)

	colabHandler := handlers.NewColaboradorHandler(colabService)
	deptHandler := handlers.NewDepartamentoHandler(deptService)

	router := gin.Default()

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		colaboradores := v1.Group("/colaboradores")
		{
			colaboradores.POST("", colabHandler.Create)
			colaboradores.GET("/:id", colabHandler.GetByID)
			colaboradores.PUT("/:id", colabHandler.Update)
			colaboradores.DELETE("/:id", colabHandler.Delete)
			colaboradores.POST("/listar", colabHandler.List)
		}

		departamentos := v1.Group("/departamentos")
		{
			departamentos.POST("", deptHandler.Create)
			departamentos.GET("/:id", deptHandler.GetByID)
			departamentos.PUT("/:id", deptHandler.Update)
			departamentos.DELETE("/:id", deptHandler.Delete)
			departamentos.POST("/listar", deptHandler.List)
		}

		gerentes := v1.Group("/gerentes")
		{
			gerentes.GET("/:id/colaboradores", deptHandler.GetGerenteColaboradores)
		}
	}

	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
