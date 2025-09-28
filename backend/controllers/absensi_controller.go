package controllers

import (
	"SIAku/config"
	"SIAku/models"
	"SIAku/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AbsensiController struct{}

func NewAbsensiController() *AbsensiController {
	return &AbsensiController{}
}

// Input Absensi per Pertemuan
func (ac *AbsensiController) InputAbsensiPertemuan(c *gin.Context) {
	dosenID, _ := c.Get("user_id")
	var req models.AbsensiPertemuanRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// Verifikasi dosen mengajar mata kuliah ini
	var course models.Course
	if err := config.DB.Where("id = ? AND dosen_id = ?", req.CourseID, dosenID).First(&course).Error; err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "You are not authorized to input attendance for this course")
		return
	}

	// Parse tanggal
	tanggal, err := time.Parse("2006-01-02", req.Tanggal)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD")
		return
	}

	var successCount int
	var errors []string

	// Process setiap absensi mahasiswa
	for _, absensiInput := range req.Absensi {
		// Verifikasi mahasiswa terdaftar di mata kuliah
		var krs models.KRS
		if err := config.DB.Where("course_id = ? AND mahasiswa_id = ? AND approval_status = 'approved'",
			req.CourseID, absensiInput.MahasiswaID).First(&krs).Error; err != nil {
			errors = append(errors, "Student ID "+string(rune(absensiInput.MahasiswaID))+" not enrolled in this course")
			continue
		}

		// Cek apakah absensi sudah ada untuk pertemuan ini
		var absensi models.Absensi
		if err := config.DB.Where("course_id = ? AND mahasiswa_id = ? AND pertemuan = ?",
			req.CourseID, absensiInput.MahasiswaID, req.Pertemuan).First(&absensi).Error; err != nil {
			// Buat absensi baru
			absensi = models.Absensi{
				CourseID:    req.CourseID,
				MahasiswaID: absensiInput.MahasiswaID,
				Pertemuan:   req.Pertemuan,
				Tanggal:     tanggal,
				Status:      absensiInput.Status,
				Keterangan:  absensiInput.Keterangan,
			}
			if err := config.DB.Create(&absensi).Error; err != nil {
				errors = append(errors, "Failed to create attendance for student ID "+string(rune(absensiInput.MahasiswaID)))
				continue
			}
		} else {
			// Update absensi yang ada
			absensi.Status = absensiInput.Status
			absensi.Keterangan = absensiInput.Keterangan
			absensi.Tanggal = tanggal

			if err := config.DB.Save(&absensi).Error; err != nil {
				errors = append(errors, "Failed to update attendance for student ID "+string(rune(absensiInput.MahasiswaID)))
				continue
			}
		}
		successCount++
	}

	response := gin.H{
		"message":       "Attendance processing completed",
		"success_count": successCount,
		"total_count":   len(req.Absensi),
		"course": gin.H{
			"id":   course.ID,
			"name": course.Name,
			"code": course.Code,
		},
		"pertemuan": req.Pertemuan,
		"tanggal":   req.Tanggal,
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	utils.SuccessResponse(c, response)
}

// Get Absensi by Course and Pertemuan
func (ac *AbsensiController) GetAbsensiByPertemuan(c *gin.Context) {
	dosenID, _ := c.Get("user_id")
	courseID := c.Param("courseId")
	pertemuan := c.Query("pertemuan")

	// Verifikasi dosen mengajar mata kuliah ini
	var course models.Course
	if err := config.DB.Where("id = ? AND dosen_id = ?", courseID, dosenID).First(&course).Error; err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "You are not authorized to view attendance for this course")
		return
	}

	query := config.DB.Preload("Mahasiswa").Preload("Course").Where("course_id = ?", courseID)

	if pertemuan != "" {
		query = query.Where("pertemuan = ?", pertemuan)
	}

	var absensiList []models.Absensi
	if err := query.Find(&absensiList).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch attendance")
		return
	}

	var responses []models.AbsensiResponse
	for _, absensi := range absensiList {
		responses = append(responses, models.AbsensiResponse{
			ID:          absensi.ID,
			CourseID:    absensi.CourseID,
			CourseName:  absensi.Course.Name,
			CourseCode:  absensi.Course.Code,
			MahasiswaID: absensi.MahasiswaID,
			NIM:         absensi.Mahasiswa.NIM,
			NamaMhs:     absensi.Mahasiswa.Nama,
			Pertemuan:   absensi.Pertemuan,
			Tanggal:     absensi.Tanggal,
			Status:      absensi.Status,
			Keterangan:  absensi.Keterangan,
			CreatedAt:   absensi.CreatedAt,
		})
	}

	utils.SuccessResponse(c, gin.H{
		"course": gin.H{
			"id":   course.ID,
			"name": course.Name,
			"code": course.Code,
		},
		"total_records": len(responses),
		"absensi":       responses,
	})
}

// Get Rekap Absensi by Course
func (ac *AbsensiController) GetRekapAbsensi(c *gin.Context) {
	dosenID, _ := c.Get("user_id")
	courseID := c.Param("courseId")

	// Verifikasi dosen mengajar mata kuliah ini
	var course models.Course
	if err := config.DB.Where("id = ? AND dosen_id = ?", courseID, dosenID).First(&course).Error; err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "You are not authorized to view attendance for this course")
		return
	}

	// Ambil semua mahasiswa yang terdaftar di mata kuliah
	var krsList []models.KRS
	if err := config.DB.Preload("Mahasiswa").Where("course_id = ? AND approval_status = 'approved'", courseID).Find(&krsList).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch enrolled students")
		return
	}

	var rekapResponses []models.RekapAbsensiResponse

	for _, krs := range krsList {
		// Hitung statistik absensi per mahasiswa
		var absensiStats struct {
			TotalHadir int
			TotalIzin  int
			TotalSakit int
			TotalAlfa  int
			Total      int
		}

		config.DB.Model(&models.Absensi{}).
			Select("COUNT(CASE WHEN status = 'hadir' THEN 1 END) as total_hadir, COUNT(CASE WHEN status = 'izin' THEN 1 END) as total_izin, COUNT(CASE WHEN status = 'sakit' THEN 1 END) as total_sakit, COUNT(CASE WHEN status = 'alfa' THEN 1 END) as total_alfa, COUNT(*) as total").
			Where("course_id = ? AND mahasiswa_id = ?", courseID, krs.MahasiswaID).
			Scan(&absensiStats)

		persentase := 0.0
		if absensiStats.Total > 0 {
			persentase = float64(absensiStats.TotalHadir) / float64(absensiStats.Total) * 100
		}

		rekapResponses = append(rekapResponses, models.RekapAbsensiResponse{
			MahasiswaID:         krs.MahasiswaID,
			NIM:                 krs.Mahasiswa.NIM,
			NamaMahasiswa:       krs.Mahasiswa.Nama,
			TotalHadir:          absensiStats.TotalHadir,
			TotalIzin:           absensiStats.TotalIzin,
			TotalSakit:          absensiStats.TotalSakit,
			TotalAlfa:           absensiStats.TotalAlfa,
			TotalPertemuan:      absensiStats.Total,
			PersentaseKehadiran: persentase,
		})
	}

	utils.SuccessResponse(c, gin.H{
		"course": gin.H{
			"id":   course.ID,
			"name": course.Name,
			"code": course.Code,
		},
		"total_students": len(rekapResponses),
		"rekap_absensi":  rekapResponses,
	})
}
