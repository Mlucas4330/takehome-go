package http

import (
	"takehome-go/internal/domain/handler"
	"takehome-go/internal/domain/repository"
	"takehome-go/internal/domain/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	colabRepo := repository.NewColaboradorRepository(db)
	depRepo := repository.NewDepartamentoRepository(db)

	colabSvc := service.NewColaboradorService(colabRepo, depRepo)
	depSvc := service.NewDepartamentoService(depRepo, colabRepo)
	gerenteSvc := service.NewGerenteService(depRepo, colabRepo)

	colabH := handler.NewColaboradorHandler(colabSvc, depSvc)
	depH := handler.NewDepartamentoHandler(depSvc)
	gerH := handler.NewGerenteHandler(gerenteSvc)

	api := r.Group("/api/v1")
	{
		api.POST("/colaboradores", colabH.Create)
		api.GET("/colaboradores/:id", colabH.Get)
		api.PUT("/colaboradores/:id", colabH.Update)
		api.DELETE("/colaboradores/:id", colabH.Delete)
		api.POST("/colaboradores/listar", colabH.List)

		api.POST("/departamentos", depH.Create)
		api.GET("/departamentos/:id", depH.Get)
		api.PUT("/departamentos/:id", depH.Update)
		api.DELETE("/departamentos/:id", depH.Delete)
		api.POST("/departamentos/listar", depH.List)

		api.GET("/gerentes/:id/colaboradores", gerH.ListSubordinates)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
