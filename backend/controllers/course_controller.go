package controllers

import (
	"SIAku/config"
	"SIAku/models"
	"SIAku/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CourseController struct{}

func NewCourseController() *CourseController {
	return &CourseController{}
}

func (cc *CourseController) GetAllCourses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var courses []models.Course
	var total int64

	config.DB.Model(&models.Course{}).Count(&total)

	if err := config.DB.Preload("Mahasiswas").Limit(limit).Offset(offset).Find(&courses).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch courses")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    courses,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (cc *CourseController) GetCourseByID(c *gin.Context) {
	id := c.Param("id")

	var course models.Course
	if err := config.DB.Preload("Mahasiswas").Where("id = ?", id).First(&course).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Course not found")
		return
	}

	utils.SuccessResponse(c, course)
}

func (cc *CourseController) CreateCourse(c *gin.Context) {
	var req models.CourseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Validasi input
	if err := utils.ValidateStruct(req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// Create course
	course := models.Course{
		Code:    req.Code,
		Name:    req.Name,
		Credits: req.Credits,
	}

	if err := config.DB.Create(&course).Error; err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "courses_code_key"` {
			utils.ErrorResponse(c, http.StatusConflict, "Course code already exists")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create course")
		return
	}

	utils.CreatedResponse(c, course)
}

// UpdateCourse - Update course (require auth)
func (cc *CourseController) UpdateCourse(c *gin.Context) {
	id := c.Param("id")

	var req models.CourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	var course models.Course
	if err := config.DB.Where("id = ?", id).First(&course).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Course not found")
		return
	}

	// Update course
	course.Code = req.Code
	course.Name = req.Name
	course.Credits = req.Credits

	if err := config.DB.Save(&course).Error; err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "courses_code_key"` {
			utils.ErrorResponse(c, http.StatusConflict, "Course code already exists")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update course")
		return
	}

	utils.SuccessResponse(c, course)
}

// DeleteCourse - Delete course (require auth)
func (cc *CourseController) DeleteCourse(c *gin.Context) {
	id := c.Param("id")

	var course models.Course
	if err := config.DB.Where("id = ?", id).First(&course).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Course not found")
		return
	}

	if err := config.DB.Delete(&course).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete course")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Course deleted successfully",
	})
}

// EnrollCourse - Enroll mahasiswa ke course (require auth)
func (cc *CourseController) EnrollCourse(c *gin.Context) {
	courseID := c.Param("id")
	userID, _ := c.Get("user_id")

	var course models.Course
	if err := config.DB.Where("id = ?", courseID).First(&course).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Course not found")
		return
	}

	var mahasiswa models.Mahasiswa
	if err := config.DB.Where("id = ?", userID).First(&mahasiswa).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Mahasiswa not found")
		return
	}

	// Check if already enrolled
	if err := config.DB.Model(&mahasiswa).Association("Courses").Find(&course, "id = ?", courseID); err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Already enrolled in this course")
		return
	}

	// Enroll mahasiswa
	if err := config.DB.Model(&mahasiswa).Association("Courses").Append(&course); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to enroll in course")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Successfully enrolled in course",
	})
}

// UnenrollCourse - Unenroll mahasiswa dari course (require auth)
func (cc *CourseController) UnenrollCourse(c *gin.Context) {
	courseID := c.Param("id")
	userID, _ := c.Get("user_id")

	var course models.Course
	if err := config.DB.Where("id = ?", courseID).First(&course).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Course not found")
		return
	}

	var mahasiswa models.Mahasiswa
	if err := config.DB.Where("id = ?", userID).First(&mahasiswa).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Mahasiswa not found")
		return
	}

	// Unenroll mahasiswa
	if err := config.DB.Model(&mahasiswa).Association("Courses").Delete(&course); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to unenroll from course")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Successfully unenrolled from course",
	})
}
