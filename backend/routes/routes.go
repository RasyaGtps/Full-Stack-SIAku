package routes

import (
	"SIAku/controllers"
	"SIAku/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	authController := controllers.NewAuthController()
	mahasiswaController := controllers.NewMahasiswaController()
	courseController := controllers.NewCourseController()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "SIAku API running ✅",
		})
	})

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		protected := api.Group("/")
		protected.Use(middleware.ValidateJWT())
		{
			protected.GET("/profile", authController.GetProfile)

			mahasiswa := protected.Group("/mahasiswa")
			{
				mahasiswa.GET("", mahasiswaController.GetAllMahasiswa)
				mahasiswa.GET("/:id", mahasiswaController.GetMahasiswaByID)
				mahasiswa.PUT("/:id", mahasiswaController.UpdateMahasiswa)
				mahasiswa.DELETE("/:id", mahasiswaController.DeleteMahasiswa)
			}

			courses := protected.Group("/courses")
			{
				courses.GET("", courseController.GetAllCourses)
				courses.GET("/:id", courseController.GetCourseByID)
				courses.POST("", courseController.CreateCourse)
				courses.PUT("/:id", courseController.UpdateCourse)
				courses.DELETE("/:id", courseController.DeleteCourse)
				courses.POST("/:id/enroll", courseController.EnrollCourse)
				courses.DELETE("/:id/enroll", courseController.UnenrollCourse)
			}
		}
	}
}
