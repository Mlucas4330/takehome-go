package httpapi

import (
	"github.com/mlucas4330/takehome-go/internal/handlers"
	"github.com/mlucas4330/takehome-go/internal/repositories"
	"github.com/mlucas4330/takehome-go/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	colRepo := repositories.NewCollaboratorRepository(db)
	deptRepo := repositories.NewDepartamentRepository(db)

	colSvc := services.NewCollaboratorService(colRepo, deptRepo)
	colH := handlers.NewCollaboratorHandler(colSvc)

	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.POST("/colaboradores", colH.Create)
		}
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
