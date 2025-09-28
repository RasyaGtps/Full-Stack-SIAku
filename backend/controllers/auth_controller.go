package controllers

import (
	"SIAku/config"
	"SIAku/middleware"
	"SIAku/models"
	"SIAku/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (ac *AuthController) Register(c *gin.Context) {
	var req models.MahasiswaRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	mahasiswa := models.Mahasiswa{
		NIM:      req.NIM,
		Nama:     req.Nama,
		Jurusan:  req.Jurusan,
		Password: string(hashedPassword),
	}

	if err := config.DB.Create(&mahasiswa).Error; err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "mahasiswas_nim_key"` {
			utils.ErrorResponse(c, http.StatusConflict, "NIM already exists")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create mahasiswa")
		return
	}

	token, err := middleware.GenerateJWT(mahasiswa.ID, mahasiswa.NIM)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
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

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Registration successful",
		"data":    response,
		"token":   token,
	})
}

func (ac *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	var mahasiswa models.Mahasiswa
	if err := config.DB.Where("nim = ?", req.NIM).First(&mahasiswa).Error; err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(mahasiswa.Password), []byte(req.Password)); err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := middleware.GenerateJWT(mahasiswa.ID, mahasiswa.NIM)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"data":    response,
		"token":   token,
	})
}

func (ac *AuthController) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var mahasiswa models.Mahasiswa
	if err := config.DB.Preload("Courses").Where("id = ?", userID).First(&mahasiswa).Error; err != nil {
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
