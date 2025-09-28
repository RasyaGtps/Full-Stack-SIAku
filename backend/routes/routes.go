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
	krsController := controllers.NewKRSController()
	nilaiController := controllers.NewNilaiController()
	jadwalController := controllers.NewJadwalController()
	dosenController := controllers.NewDosenController()
	absensiController := controllers.NewAbsensiController()
	materiController := controllers.NewMateriController()

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

			krs := protected.Group("/krs")
			{
				krs.GET("", krsController.GetMyKRS)
				krs.POST("", krsController.AddCourseToKRS)
				krs.DELETE("/:id", krsController.RemoveCourseFromKRS)
				krs.GET("/available-courses", krsController.GetAvailableCourses)
			}

			nilai := protected.Group("/nilai")
			{
				nilai.GET("", nilaiController.GetMyNilai)
				nilai.GET("/transkrip", nilaiController.GetTranskrip)
				nilai.GET("/statistik", nilaiController.GetStatistikNilai)
			}

			jadwal := protected.Group("/jadwal")
			{
				jadwal.GET("", jadwalController.GetMyJadwal)
				jadwal.GET("/hari/:hari", jadwalController.GetJadwalByHari)
				jadwal.GET("/minggu-ini", jadwalController.GetJadwalMingguIni)
			}

			// Dosen endpoints
			dosen := protected.Group("/dosen")
			{
				// Input nilai mahasiswa
				dosen.POST("/courses/:courseId/students/:mahasiswaId/nilai", dosenController.InputNilai)

				// Lihat daftar mahasiswa di kelas
				dosen.GET("/courses/:courseId/students", dosenController.GetMahasiswaInClass)

				// Approve/Reject KRS
				dosen.GET("/krs/pending", dosenController.GetPendingKRS)
				dosen.PUT("/krs/:krsId/approval", dosenController.ProcessKRSApproval)

				// Get my courses
				dosen.GET("/courses", materiController.GetMyCourses)
			}

			// Absensi endpoints
			absensi := protected.Group("/absensi")
			{
				absensi.POST("/input", absensiController.InputAbsensiPertemuan)
				absensi.GET("/courses/:courseId", absensiController.GetAbsensiByPertemuan)
				absensi.GET("/courses/:courseId/rekap", absensiController.GetRekapAbsensi)
			}

			// Materi endpoints
			materi := protected.Group("/materi")
			{
				materi.POST("", materiController.CreateMateri)
				materi.GET("/courses/:courseId", materiController.GetMateriByCourse)
				materi.PUT("/:materiId", materiController.UpdateMateri)
				materi.DELETE("/:materiId", materiController.DeleteMateri)
			}
		}
	}
}
