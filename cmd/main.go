package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"takehome-go/internal/config"
	"takehome-go/internal/database"
	"takehome-go/internal/http"
)

func main() {
	cfg := config.Load()
	db := database.Open(cfg.DSN)
	r := gin.Default()

	http.SetupRoutes(r, db)

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(fmt.Sprintf(":%s", cfg.Port))
}