package controllers

import (
	"SIAku/config"
	"SIAku/models"
	"SIAku/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type DosenController struct{}

func NewDosenController() *DosenController {
	return &DosenController{}
}

// Input Nilai Mahasiswa per Mata Kuliah
func (dc *DosenController) InputNilai(c *gin.Context) {
	dosenID, _ := c.Get("user_id")
	courseID := c.Param("courseId")
	mahasiswaID := c.Param("mahasiswaId")

	var req struct {
		NilaiTugas float64 `json:"nilai_tugas" validate:"min=0,max=100"`
		NilaiUTS   float64 `json:"nilai_uts" validate:"min=0,max=100"`
		NilaiUAS   float64 `json:"nilai_uas" validate:"min=0,max=100"`
	}

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
	if err := config.DB.Where("id = ? AND dosen_id = ?", courseID, dosenID).First(&course).Error; err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "You are not authorized to input grades for this course")
		return
	}

	// Verifikasi mahasiswa terdaftar di mata kuliah
	var krs models.KRS
	if err := config.DB.Where("course_id = ? AND mahasiswa_id = ? AND approval_status = 'approved'", courseID, mahasiswaID).First(&krs).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Student not enrolled in this course")
		return
	}

	// Hitung nilai akhir (30% tugas, 35% UTS, 35% UAS)
	nilaiAkhir := (req.NilaiTugas * 0.3) + (req.NilaiUTS * 0.35) + (req.NilaiUAS * 0.35)

	// Tentukan grade huruf dan poin
	gradeHuruf, gradePoint := calculateGrade(nilaiAkhir)

	// Cek apakah nilai sudah ada
	var nilai models.Nilai
	if err := config.DB.Where("mahasiswa_id = ? AND course_id = ?", mahasiswaID, courseID).First(&nilai).Error; err != nil {
		// Buat nilai baru
		nilai = models.Nilai{
			MahasiswaID: uint(parseUint(mahasiswaID)),
			CourseID:    uint(parseUint(courseID)),
			Semester:    krs.Semester,
			TahunAjaran: krs.TahunAjaran,
			NilaiTugas:  req.NilaiTugas,
			NilaiUTS:    req.NilaiUTS,
			NilaiUAS:    req.NilaiUAS,
			NilaiAkhir:  nilaiAkhir,
			GradeHuruf:  gradeHuruf,
			GradePoint:  gradePoint,
			Status:      "sudah_dinilai",
		}
		if err := config.DB.Create(&nilai).Error; err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create grade")
			return
		}
	} else {
		// Update nilai yang ada
		nilai.NilaiTugas = req.NilaiTugas
		nilai.NilaiUTS = req.NilaiUTS
		nilai.NilaiUAS = req.NilaiUAS
		nilai.NilaiAkhir = nilaiAkhir
		nilai.GradeHuruf = gradeHuruf
		nilai.GradePoint = gradePoint
		nilai.Status = "sudah_dinilai"

		if err := config.DB.Save(&nilai).Error; err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update grade")
			return
		}
	}

	// Notification removed - WhatsApp integration disabled

	utils.SuccessResponse(c, gin.H{
		"message": "Grade successfully inputted",
		"nilai": gin.H{
			"nilai_tugas": req.NilaiTugas,
			"nilai_uts":   req.NilaiUTS,
			"nilai_uas":   req.NilaiUAS,
			"nilai_akhir": nilaiAkhir,
			"grade_huruf": gradeHuruf,
			"grade_point": gradePoint,
		},
	})
}

// Lihat Daftar Mahasiswa di Kelas
func (dc *DosenController) GetMahasiswaInClass(c *gin.Context) {
	dosenID, _ := c.Get("user_id")
	courseID := c.Param("courseId")

	// Verifikasi dosen mengajar mata kuliah ini
	var course models.Course
	if err := config.DB.Where("id = ? AND dosen_id = ?", courseID, dosenID).First(&course).Error; err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "You are not authorized to view this class")
		return
	}

	// Ambil daftar mahasiswa yang KRS-nya sudah approved
	var krsList []models.KRS
	if err := config.DB.Preload("Mahasiswa").Where("course_id = ? AND approval_status = 'approved'", courseID).Find(&krsList).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch students")
		return
	}

	var students []gin.H
	for _, krs := range krsList {
		students = append(students, gin.H{
			"id":              krs.Mahasiswa.ID,
			"nim":             krs.Mahasiswa.NIM,
			"nama":            krs.Mahasiswa.Nama,
			"jurusan":         krs.Mahasiswa.Jurusan,
			"semester":        krs.Mahasiswa.Semester,
			"status_akademik": krs.Mahasiswa.StatusAkademik,
			"enrolled_at":     krs.CreatedAt,
		})
	}

	utils.SuccessResponse(c, gin.H{
		"course": gin.H{
			"id":      course.ID,
			"code":    course.Code,
			"name":    course.Name,
			"credits": course.Credits,
		},
		"total_students": len(students),
		"students":       students,
	})
}

// Approve/Reject KRS Mahasiswa
func (dc *DosenController) ProcessKRSApproval(c *gin.Context) {
	dosenID, _ := c.Get("user_id")
	krsID := c.Param("krsId")

	var req models.KRSApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// Ambil KRS dan verifikasi dosen wali
	var krs models.KRS
	if err := config.DB.Preload("Mahasiswa").Where("id = ?", krsID).First(&krs).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "KRS not found")
		return
	}

	// Verifikasi dosen adalah dosen wali mahasiswa
	if krs.Mahasiswa.DosenWaliID == nil || *krs.Mahasiswa.DosenWaliID != dosenID.(uint) {
		utils.ErrorResponse(c, http.StatusForbidden, "You are not authorized to approve this KRS")
		return
	}

	// Update status approval
	now := time.Now()
	dosenIDUint := dosenID.(uint)
	if req.Action == "approve" {
		krs.ApprovalStatus = "approved"
		krs.Status = "diambil"
		krs.ApprovedBy = &dosenIDUint
		krs.ApprovedAt = &now
		krs.RejectionReason = ""
	} else {
		krs.ApprovalStatus = "rejected"
		krs.Status = "ditolak"
		krs.ApprovedBy = &dosenIDUint
		krs.ApprovedAt = &now
		krs.RejectionReason = req.RejectionReason
	}

	if err := config.DB.Save(&krs).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to process KRS approval")
		return
	}

	// Notification removed - WhatsApp integration disabled

	utils.SuccessResponse(c, gin.H{
		"message": "KRS " + req.Action + "d successfully",
		"krs": gin.H{
			"id":               krs.ID,
			"approval_status":  krs.ApprovalStatus,
			"status":           krs.Status,
			"approved_by":      krs.ApprovedBy,
			"approved_at":      krs.ApprovedAt,
			"rejection_reason": krs.RejectionReason,
		},
	})
}

// Get Pending KRS for Approval
func (dc *DosenController) GetPendingKRS(c *gin.Context) {
	dosenID, _ := c.Get("user_id")

	var pendingKRS []models.KRS
	if err := config.DB.Preload("Mahasiswa").Preload("Course").
		Joins("JOIN mahasiswas ON mahasiswas.id = krs.mahasiswa_id").
		Where("mahasiswas.dosen_wali_id = ? AND krs.approval_status = 'pending'", dosenID).
		Find(&pendingKRS).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch pending KRS")
		return
	}

	var responses []gin.H
	for _, krs := range pendingKRS {
		responses = append(responses, gin.H{
			"id": krs.ID,
			"mahasiswa": gin.H{
				"id":       krs.Mahasiswa.ID,
				"nim":      krs.Mahasiswa.NIM,
				"nama":     krs.Mahasiswa.Nama,
				"semester": krs.Mahasiswa.Semester,
			},
			"course": gin.H{
				"id":      krs.Course.ID,
				"code":    krs.Course.Code,
				"name":    krs.Course.Name,
				"credits": krs.Course.Credits,
			},
			"semester":     krs.Semester,
			"tahun_ajaran": krs.TahunAjaran,
			"created_at":   krs.CreatedAt,
		})
	}

	utils.SuccessResponse(c, gin.H{
		"total_pending": len(responses),
		"pending_krs":   responses,
	})
}

// Helper functions
func parseUint(s string) uint64 {
	val, _ := strconv.ParseUint(s, 10, 32)
	return val
}

func calculateGrade(nilai float64) (string, float64) {
	if nilai >= 85 {
		return "A", 4.0
	} else if nilai >= 80 {
		return "AB", 3.5
	} else if nilai >= 75 {
		return "B", 3.0
	} else if nilai >= 70 {
		return "BC", 2.5
	} else if nilai >= 65 {
		return "C", 2.0
	} else if nilai >= 50 {
		return "D", 1.0
	} else {
		return "E", 0.0
	}
}
