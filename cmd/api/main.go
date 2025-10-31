package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/mlucas4330/takehome-go/internal/httpapi"
	"github.com/mlucas4330/takehome-go/internal/infrastructure/config"
	"github.com/mlucas4330/takehome-go/internal/infrastructure/database"
)

func main() {
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s", cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort, cfg.DBSSLMode, cfg.DBTimezone)

	db := database.Open(dsn)
	r := gin.Default()

	httpapi.SetupRoutes(r, db)

	r.Run(fmt.Sprintf(":%d", cfg.Port))
}
