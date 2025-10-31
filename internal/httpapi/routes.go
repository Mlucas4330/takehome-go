package httpapi

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mlucas4330/takehome-go/internal/docs"
	"github.com/mlucas4330/takehome-go/internal/handlers"
	"github.com/mlucas4330/takehome-go/internal/repositories"
	"github.com/mlucas4330/takehome-go/internal/services/collaborator"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	colRepo := repositories.NewCollaboratorRepository(db)
	deptRepo := repositories.NewDepartamentRepository(db)

	colSvc := collaborator.NewCollaboratorService(colRepo, deptRepo)
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

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
