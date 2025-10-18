package controllers

import (
	"SIAku/config"
	"SIAku/middleware"
	"SIAku/models"
	"SIAku/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

// Register - Universal registration for all roles
func (ac *AuthController) Register(c *gin.Context) {
	var req models.UserRegistrationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Start database transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create user record
	user := models.Users{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
		Status:   "aktif",
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "username") || strings.Contains(errorMsg, "uni_users_username") {
			utils.ErrorResponse(c, http.StatusConflict, "Username sudah terdaftar, silakan gunakan username yang lain")
			return
		}
		if strings.Contains(errorMsg, "email") || strings.Contains(errorMsg, "uni_users_email") {
			utils.ErrorResponse(c, http.StatusConflict, "Email sudah terdaftar, silakan gunakan email yang lain")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat akun, silakan coba lagi")
		return
	}

	// Create role-specific record
	switch req.Role {
	case "mahasiswa":
		if req.NIM == "" || req.Jurusan == "" {
			tx.Rollback()
			utils.ErrorResponse(c, http.StatusBadRequest, "NIM and Jurusan are required for mahasiswa")
			return
		}

		semester := req.Semester
		if semester == 0 {
			semester = 1
		}

		statusAkademik := req.StatusAkademik
		if statusAkademik == "" {
			statusAkademik = "aktif"
		}

		mahasiswa := models.Mahasiswa{
			UserID:         &user.ID,
			NIM:            req.NIM,
			Nama:           req.Nama,
			Jurusan:        req.Jurusan,
			PhoneNumber:    req.PhoneNumber,
			StatusAkademik: statusAkademik,
			Semester:       semester,
			IPK:            0.00,
		}

		if err := tx.Create(&mahasiswa).Error; err != nil {
			tx.Rollback()
			errorMsg := err.Error()
			if strings.Contains(errorMsg, "nim") {
				utils.ErrorResponse(c, http.StatusConflict, "NIM sudah terdaftar, silakan gunakan NIM yang lain")
				return
			}
			utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat data mahasiswa")
			return
		}

	case "dosen":
		if req.NIDN == "" || req.Jurusan == "" {
			tx.Rollback()
			utils.ErrorResponse(c, http.StatusBadRequest, "NIDN and Jurusan are required for dosen")
			return
		}

		dosen := models.Dosen{
			NIDN:        req.NIDN,
			Nama:        req.Nama,
			Email:       req.Email,
			PhoneNumber: req.PhoneNumber,
			Jurusan:     req.Jurusan,
			Status:      "aktif",
		}

		if err := tx.Create(&dosen).Error; err != nil {
			tx.Rollback()
			errorMsg := err.Error()
			if strings.Contains(errorMsg, "nidn") {
				utils.ErrorResponse(c, http.StatusConflict, "NIDN sudah terdaftar, silakan gunakan NIDN yang lain")
				return
			}
			utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat data dosen")
			return
		}

	case "kajur":
		if req.NIDN == "" || req.Jurusan == "" {
			tx.Rollback()
			utils.ErrorResponse(c, http.StatusBadRequest, "NIDN and Jurusan are required for kajur")
			return
		}

		kajur := models.Kajur{
			NIDN:    req.NIDN,
			Nama:    req.Nama,
			Email:   req.Email,
			Jurusan: req.Jurusan,
			Status:  "aktif",
		}

		if err := tx.Create(&kajur).Error; err != nil {
			tx.Rollback()
			errorMsg := err.Error()
			if strings.Contains(errorMsg, "nidn") {
				utils.ErrorResponse(c, http.StatusConflict, "NIDN sudah terdaftar sebagai kajur, silakan gunakan NIDN yang lain")
				return
			}
			utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat data kajur")
			return
		}

	case "rektor":
		if req.NIDN == "" {
			tx.Rollback()
			utils.ErrorResponse(c, http.StatusBadRequest, "NIDN is required for rektor")
			return
		}

		rektor := models.Rektor{
			NIDN:   req.NIDN,
			Nama:   req.Nama,
			Email:  req.Email,
			Status: "aktif",
		}

		if err := tx.Create(&rektor).Error; err != nil {
			tx.Rollback()
			errorMsg := err.Error()
			if strings.Contains(errorMsg, "nidn") {
				utils.ErrorResponse(c, http.StatusConflict, "NIDN sudah terdaftar sebagai rektor, silakan gunakan NIDN yang lain")
				return
			}
			utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat data rektor")
			return
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to complete registration")
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateJWT(user.ID, user.Username)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Get complete user data with details
	response := ac.buildUserResponse(user)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Registration successful",
		"data":    response,
		"token":   token,
	})
}

// Login - Universal login for all roles
func (ac *AuthController) Login(c *gin.Context) {
	var req models.UserLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	var user models.Users

	// Find user by username or email
	err := config.DB.Where("username = ? OR email = ?", req.Identifier, req.Identifier).First(&user).Error

	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Username atau email tidak ditemukan")
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateJWT(user.ID, user.Username)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Get complete user data with details
	response := ac.buildUserResponse(user)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"data":    response,
		"token":   token,
	})
}

// GetProfile - Get user profile with role details
func (ac *AuthController) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var user models.Users
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	response := ac.buildUserResponse(user)
	utils.SuccessResponse(c, response)
}

// Helper function to build complete user response with details
func (ac *AuthController) buildUserResponse(user models.Users) models.UserResponse {
	response := models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// Load role-specific data from respective tables
	switch user.Role {
	case "mahasiswa":
		var mahasiswa models.Mahasiswa
		if err := config.DB.Where("user_id = ?", user.ID).First(&mahasiswa).Error; err == nil {
			response.RoleData = mahasiswa
		}

	case "dosen":
		var dosen models.Dosen
		if err := config.DB.Where("email = ?", user.Email).First(&dosen).Error; err == nil {
			response.RoleData = dosen
		}

	case "kajur":
		var kajur models.Kajur
		if err := config.DB.Where("email = ?", user.Email).First(&kajur).Error; err == nil {
			response.RoleData = kajur
		}

	case "rektor":
		var rektor models.Rektor
		if err := config.DB.Where("email = ?", user.Email).First(&rektor).Error; err == nil {
			response.RoleData = rektor
		}
	}

	return response
}
