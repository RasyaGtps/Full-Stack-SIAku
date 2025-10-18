package controllers

import (
	"SIAku/config"
	"SIAku/models"
	"SIAku/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MahasiswaController struct{}

func NewMahasiswaController() *MahasiswaController {
	return &MahasiswaController{}
}

func (mc *MahasiswaController) GetAllMahasiswa(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var mahasiswas []models.Mahasiswa
	var total int64

	config.DB.Model(&models.Mahasiswa{}).Count(&total)

	if err := config.DB.Preload("Courses").Limit(limit).Offset(offset).Find(&mahasiswas).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch mahasiswa")
		return
	}

	var responses []models.MahasiswaResponse
	for _, mhs := range mahasiswas {
		responses = append(responses, models.MahasiswaResponse{
			ID:        mhs.ID,
			NIM:       mhs.NIM,
			Nama:      mhs.Nama,
			Jurusan:   mhs.Jurusan,
			Courses:   mhs.Courses,
			CreatedAt: mhs.CreatedAt,
			UpdatedAt: mhs.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    responses,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (mc *MahasiswaController) GetMahasiswaByID(c *gin.Context) {
	id := c.Param("id")

	var mahasiswa models.Mahasiswa
	if err := config.DB.Preload("Courses").Where("id = ?", id).First(&mahasiswa).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Mahasiswa not found")
		return
	}

	response := models.MahasiswaResponse{
		ID:        mahasiswa.ID,
		NIM:       mahasiswa.NIM,
		Nama:      mahasiswa.Nama,
		Jurusan:   mahasiswa.Jurusan,
		Courses:   mahasiswa.Courses,
		CreatedAt: mahasiswa.CreatedAt,
		UpdatedAt: mahasiswa.UpdatedAt,
	}

	utils.SuccessResponse(c, response)
}

func (mc *MahasiswaController) GetMahasiswaByNIM(c *gin.Context) {
	nim := c.Param("nim")

	var mahasiswa models.Mahasiswa
	if err := config.DB.Preload("Courses").Where("nim = ?", nim).First(&mahasiswa).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Mahasiswa dengan NIM "+nim+" tidak ditemukan")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":              mahasiswa.ID,
			"nim":             mahasiswa.NIM,
			"nama":            mahasiswa.Nama,
			"jurusan":         mahasiswa.Jurusan,
			"phone_number":    mahasiswa.PhoneNumber,
			"status_akademik": mahasiswa.StatusAkademik,
			"semester":        mahasiswa.Semester,
			"ipk":             mahasiswa.IPK,
			"total_courses":   len(mahasiswa.Courses),
		},
	})
}

func (mc *MahasiswaController) UpdateMahasiswa(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	if strconv.Itoa(int(userID.(uint))) != id {
		utils.ErrorResponse(c, http.StatusForbidden, "You can only update your own data")
		return
	}

	var req struct {
		Nama    string `json:"nama" validate:"omitempty,min=2,max=100"`
		Jurusan string `json:"jurusan" validate:"omitempty,min=2,max=100"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	var mahasiswa models.Mahasiswa
	if err := config.DB.Where("id = ?", id).First(&mahasiswa).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Mahasiswa not found")
		return
	}

	// Update only provided fields
	if req.Nama != "" {
		mahasiswa.Nama = req.Nama
	}
	if req.Jurusan != "" {
		mahasiswa.Jurusan = req.Jurusan
	}

	if err := config.DB.Save(&mahasiswa).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update mahasiswa")
		return
	}

	response := models.MahasiswaResponse{
		ID:        mahasiswa.ID,
		NIM:       mahasiswa.NIM,
		Nama:      mahasiswa.Nama,
		Jurusan:   mahasiswa.Jurusan,
		CreatedAt: mahasiswa.CreatedAt,
		UpdatedAt: mahasiswa.UpdatedAt,
	}

	utils.SuccessResponse(c, response)
}

func (mc *MahasiswaController) DeleteMahasiswa(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	if strconv.Itoa(int(userID.(uint))) != id {
		utils.ErrorResponse(c, http.StatusForbidden, "You can only delete your own data")
		return
	}

	var mahasiswa models.Mahasiswa
	if err := config.DB.Where("id = ?", id).First(&mahasiswa).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Mahasiswa not found")
		return
	}

	if err := config.DB.Delete(&mahasiswa).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete mahasiswa")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Mahasiswa deleted successfully",
	})
}
