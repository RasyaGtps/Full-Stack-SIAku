package main

import (
	"SIAku/config"
	"SIAku/models"
	"log"

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

	if err := db.AutoMigrate(&models.Mahasiswa{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "SIAku API running ✅",
		})
	})

	r.GET("/mahasiswa", func(c *gin.Context) {
		var mhs []models.Mahasiswa
		if err := db.Find(&mhs).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, mhs)
	})

	r.POST("/mahasiswa", func(c *gin.Context) {
		var input models.Mahasiswa
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&input).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, input)
	})

	port := config.AppConfig.ServerPort
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
