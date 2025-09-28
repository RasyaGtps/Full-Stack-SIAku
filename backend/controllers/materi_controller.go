package controllers

import (
	"SIAku/config"
	"SIAku/models"
	"SIAku/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MateriController struct{}

func NewMateriController() *MateriController {
	return &MateriController{}
}

// Upload/Create Materi Kuliah
func (mc *MateriController) CreateMateri(c *gin.Context) {
	dosenID, _ := c.Get("user_id")
	var req models.MateriRequest

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
		utils.ErrorResponse(c, http.StatusForbidden, "You are not authorized to upload material for this course")
		return
	}

	// Buat materi baru
	materi := models.Materi{
		CourseID:   req.CourseID,
		Judul:      req.Judul,
		Deskripsi:  req.Deskripsi,
		Pertemuan:  req.Pertemuan,
		TipeMateri: req.TipeMateri,
		URL:        req.URL,
		Status:     "aktif",
	}

	if err := config.DB.Create(&materi).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create material")
		return
	}

	// Load course info untuk response
	config.DB.Preload("Course").First(&materi, materi.ID)

	response := models.MateriResponse{
		ID:         materi.ID,
		CourseID:   materi.CourseID,
		CourseName: course.Name,
		CourseCode: course.Code,
		Judul:      materi.Judul,
		Deskripsi:  materi.Deskripsi,
		Pertemuan:  materi.Pertemuan,
		TipeMateri: materi.TipeMateri,
		FilePath:   materi.FilePath,
		FileSize:   materi.FileSize,
		URL:        materi.URL,
		Status:     materi.Status,
		CreatedAt:  materi.CreatedAt,
		UpdatedAt:  materi.UpdatedAt,
	}

	utils.CreatedResponse(c, response)
}

// Get Materi by Course
func (mc *MateriController) GetMateriByCourse(c *gin.Context) {
	dosenID, _ := c.Get("user_id")
	courseID := c.Param("courseId")
	pertemuan := c.Query("pertemuan")

	// Verifikasi dosen mengajar mata kuliah ini
	var course models.Course
	if err := config.DB.Where("id = ? AND dosen_id = ?", courseID, dosenID).First(&course).Error; err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "You are not authorized to view materials for this course")
		return
	}

	query := config.DB.Where("course_id = ? AND status = 'aktif'", courseID)

	if pertemuan != "" {
		query = query.Where("pertemuan = ?", pertemuan)
	}

	var materiList []models.Materi
	if err := query.Order("pertemuan ASC, created_at ASC").Find(&materiList).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch materials")
		return
	}

	var responses []models.MateriResponse
	for _, materi := range materiList {
		responses = append(responses, models.MateriResponse{
			ID:         materi.ID,
			CourseID:   materi.CourseID,
			CourseName: course.Name,
			CourseCode: course.Code,
			Judul:      materi.Judul,
			Deskripsi:  materi.Deskripsi,
			Pertemuan:  materi.Pertemuan,
			TipeMateri: materi.TipeMateri,
			FilePath:   materi.FilePath,
			FileSize:   materi.FileSize,
			URL:        materi.URL,
			Status:     materi.Status,
			CreatedAt:  materi.CreatedAt,
			UpdatedAt:  materi.UpdatedAt,
		})
	}

	utils.SuccessResponse(c, gin.H{
		"course": gin.H{
			"id":   course.ID,
			"name": course.Name,
			"code": course.Code,
		},
		"total_materials": len(responses),
		"materials":       responses,
	})
}

// Update Materi
func (mc *MateriController) UpdateMateri(c *gin.Context) {
	dosenID, _ := c.Get("user_id")
	materiID := c.Param("materiId")

	var req models.MateriRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// Ambil materi dan verifikasi kepemilikan
	var materi models.Materi
	if err := config.DB.Preload("Course").Where("id = ?", materiID).First(&materi).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Material not found")
		return
	}

	// Verifikasi dosen mengajar mata kuliah ini
	if materi.Course.DosenID == nil || *materi.Course.DosenID != dosenID.(uint) {
		utils.ErrorResponse(c, http.StatusForbidden, "You are not authorized to update this material")
		return
	}

	// Update materi
	materi.Judul = req.Judul
	materi.Deskripsi = req.Deskripsi
	materi.Pertemuan = req.Pertemuan
	materi.TipeMateri = req.TipeMateri
	materi.URL = req.URL

	if err := config.DB.Save(&materi).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update material")
		return
	}

	response := models.MateriResponse{
		ID:         materi.ID,
		CourseID:   materi.CourseID,
		CourseName: materi.Course.Name,
		CourseCode: materi.Course.Code,
		Judul:      materi.Judul,
		Deskripsi:  materi.Deskripsi,
		Pertemuan:  materi.Pertemuan,
		TipeMateri: materi.TipeMateri,
		FilePath:   materi.FilePath,
		FileSize:   materi.FileSize,
		URL:        materi.URL,
		Status:     materi.Status,
		CreatedAt:  materi.CreatedAt,
		UpdatedAt:  materi.UpdatedAt,
	}

	utils.SuccessResponse(c, response)
}

// Delete Materi (soft delete - set status to inactive)
func (mc *MateriController) DeleteMateri(c *gin.Context) {
	dosenID, _ := c.Get("user_id")
	materiID := c.Param("materiId")

	// Ambil materi dan verifikasi kepemilikan
	var materi models.Materi
	if err := config.DB.Preload("Course").Where("id = ?", materiID).First(&materi).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Material not found")
		return
	}

	// Verifikasi dosen mengajar mata kuliah ini
	if materi.Course.DosenID == nil || *materi.Course.DosenID != dosenID.(uint) {
		utils.ErrorResponse(c, http.StatusForbidden, "You are not authorized to delete this material")
		return
	}

	// Soft delete - set status to inactive
	materi.Status = "inactive"
	if err := config.DB.Save(&materi).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete material")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Material deleted successfully",
	})
}

// Get All My Courses (untuk dosen)
func (mc *MateriController) GetMyCourses(c *gin.Context) {
	dosenID, _ := c.Get("user_id")

	var courses []models.Course
	if err := config.DB.Where("dosen_id = ?", dosenID).Find(&courses).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch courses")
		return
	}

	var responses []gin.H
	for _, course := range courses {
		// Hitung jumlah mahasiswa terdaftar
		var studentCount int64
		config.DB.Model(&models.KRS{}).Where("course_id = ? AND approval_status = 'approved'", course.ID).Count(&studentCount)

		// Hitung jumlah materi
		var materialCount int64
		config.DB.Model(&models.Materi{}).Where("course_id = ? AND status = 'aktif'", course.ID).Count(&materialCount)

		responses = append(responses, gin.H{
			"id":             course.ID,
			"code":           course.Code,
			"name":           course.Name,
			"credits":        course.Credits,
			"semester":       course.Semester,
			"deskripsi":      course.Deskripsi,
			"student_count":  studentCount,
			"material_count": materialCount,
			"created_at":     course.CreatedAt,
		})
	}

	utils.SuccessResponse(c, gin.H{
		"total_courses": len(responses),
		"courses":       responses,
	})
}
