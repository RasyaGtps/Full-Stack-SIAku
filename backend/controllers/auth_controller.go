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
		if err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"` {
			utils.ErrorResponse(c, http.StatusConflict, "Username already exists")
			return
		}
		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			utils.ErrorResponse(c, http.StatusConflict, "Email already exists")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
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
			NIM:            req.NIM,
			Nama:           req.Nama,
			Jurusan:        req.Jurusan,
			StatusAkademik: statusAkademik,
			Semester:       semester,
			IPK:            0.00,
		}

		if err := tx.Create(&mahasiswa).Error; err != nil {
			tx.Rollback()
			if err.Error() == `pq: duplicate key value violates unique constraint "mahasiswas_nim_key"` {
				utils.ErrorResponse(c, http.StatusConflict, "NIM already exists")
				return
			}
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create mahasiswa")
			return
		}

	case "dosen":
		if req.NIDN == "" || req.Jurusan == "" {
			tx.Rollback()
			utils.ErrorResponse(c, http.StatusBadRequest, "NIDN and Jurusan are required for dosen")
			return
		}

		dosen := models.Dosen{
			NIDN:    req.NIDN,
			Nama:    req.Nama,
			Email:   req.Email,
			Jurusan: req.Jurusan,
			Status:  "aktif",
		}

		if err := tx.Create(&dosen).Error; err != nil {
			tx.Rollback()
			if err.Error() == `pq: duplicate key value violates unique constraint "dosens_nidn_key"` {
				utils.ErrorResponse(c, http.StatusConflict, "NIDN already exists")
				return
			}
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create dosen")
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
			if err.Error() == `pq: duplicate key value violates unique constraint "kajurs_nidn_key"` {
				utils.ErrorResponse(c, http.StatusConflict, "NIDN already exists")
				return
			}
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create kajur")
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
			if err.Error() == `pq: duplicate key value violates unique constraint "rektors_nidn_key"` {
				utils.ErrorResponse(c, http.StatusConflict, "NIDN already exists")
				return
			}
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create rektor")
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

	// Find user by username or email first
	query := config.DB.Where("username = ? OR email = ?", req.Identifier, req.Identifier)

	// Also check in role tables for nim/nidn and find corresponding user
	var mahasiswa models.Mahasiswa
	if err := config.DB.Where("nim = ?", req.Identifier).First(&mahasiswa).Error; err == nil {
		// Find user by matching email from mahasiswa
		query = config.DB.Where("email = ?", mahasiswa.Nama) // Assume we can match by some field
	} else {
		var dosen models.Dosen
		if err := config.DB.Where("nidn = ?", req.Identifier).First(&dosen).Error; err == nil {
			query = config.DB.Where("email = ?", dosen.Email)
		} else {
			var kajur models.Kajur
			if err := config.DB.Where("nidn = ?", req.Identifier).First(&kajur).Error; err == nil {
				query = config.DB.Where("email = ?", kajur.Email)
			} else {
				var rektor models.Rektor
				if err := config.DB.Where("nidn = ?", req.Identifier).First(&rektor).Error; err == nil {
					query = config.DB.Where("email = ?", rektor.Email)
				}
			}
		}
	}

	if err := query.First(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
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
		if err := config.DB.Where("nama = ?", user.Username).Or("nim = ?", user.Username).First(&mahasiswa).Error; err == nil {
			response.RoleData = mahasiswa
		}

	case "dosen":
		var dosen models.Dosen
		if err := config.DB.Where("email = ?", user.Email).Or("nidn = ?", user.Username).First(&dosen).Error; err == nil {
			response.RoleData = dosen
		}

	case "kajur":
		var kajur models.Kajur
		if err := config.DB.Where("email = ?", user.Email).Or("nidn = ?", user.Username).First(&kajur).Error; err == nil {
			response.RoleData = kajur
		}

	case "rektor":
		var rektor models.Rektor
		if err := config.DB.Where("email = ?", user.Email).Or("nidn = ?", user.Username).First(&rektor).Error; err == nil {
			response.RoleData = rektor
		}
	}

	return response
}
