package main

import (
	"SIAku/config"
	"SIAku/middleware"
	"SIAku/models"
	"SIAku/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Old migration
	if err := db.AutoMigrate(&models.Mahasiswa{}, &models.Course{}, &models.KRS{}, &models.Nilai{}, &models.Jadwal{}, &models.Dosen{}, &models.Absensi{}, &models.Materi{}, &models.Kajur{}, &models.Rektor{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Users table migration
	if err := db.AutoMigrate(&models.Users{}); err != nil {
		log.Fatalf("Users table migration failed: %v", err)
	}

	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	// Setup routes
	routes.SetupRoutes(r)

	port := config.AppConfig.ServerPort
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running on :%s", port)
	log.Printf("📋 Environment: %s", gin.Mode())

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
