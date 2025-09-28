package controllers

import (
	"SIAku/config"
	"SIAku/models"
	"SIAku/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type KRSController struct{}

func NewKRSController() *KRSController {
	return &KRSController{}
}

func (kc *KRSController) GetMyKRS(c *gin.Context) {
	userID, _ := c.Get("user_id")
	semester := c.DefaultQuery("semester", "")
	tahunAjaran := c.DefaultQuery("tahun_ajaran", "")

	var krs []models.KRS
	query := config.DB.Preload("Course").Where("mahasiswa_id = ?", userID)

	if semester != "" {
		query = query.Where("semester = ?", semester)
	}
	if tahunAjaran != "" {
		query = query.Where("tahun_ajaran = ?", tahunAjaran)
	}

	if err := query.Find(&krs).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch KRS")
		return
	}

	var responses []models.KRSResponse
	for _, k := range krs {
		responses = append(responses, models.KRSResponse{
			ID:          k.ID,
			CourseID:    k.CourseID,
			CourseName:  k.Course.Name,
			CourseCode:  k.Course.Code,
			Credits:     k.Course.Credits,
			Semester:    k.Semester,
			TahunAjaran: k.TahunAjaran,
			Status:      k.Status,
			CreatedAt:   k.CreatedAt,
		})
	}

	utils.SuccessResponse(c, responses)
}

func (kc *KRSController) AddCourseToKRS(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req models.KRSRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	var course models.Course
	if err := config.DB.Where("id = ?", req.CourseID).First(&course).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Course not found")
		return
	}

	var existingKRS models.KRS
	if err := config.DB.Where("mahasiswa_id = ? AND course_id = ? AND semester = ? AND tahun_ajaran = ?",
		userID, req.CourseID, req.Semester, req.TahunAjaran).First(&existingKRS).Error; err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Course already added to KRS")
		return
	}

	krs := models.KRS{
		MahasiswaID: userID.(uint),
		CourseID:    req.CourseID,
		Semester:    req.Semester,
		TahunAjaran: req.TahunAjaran,
		Status:      "diambil",
	}

	if err := config.DB.Create(&krs).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to add course to KRS")
		return
	}

	response := models.KRSResponse{
		ID:          krs.ID,
		CourseID:    course.ID,
		CourseName:  course.Name,
		CourseCode:  course.Code,
		Credits:     course.Credits,
		Semester:    krs.Semester,
		TahunAjaran: krs.TahunAjaran,
		Status:      krs.Status,
		CreatedAt:   krs.CreatedAt,
	}

	utils.CreatedResponse(c, response)
}

func (kc *KRSController) RemoveCourseFromKRS(c *gin.Context) {
	userID, _ := c.Get("user_id")
	krsID := c.Param("id")

	var krs models.KRS
	if err := config.DB.Where("id = ? AND mahasiswa_id = ?", krsID, userID).First(&krs).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "KRS entry not found")
		return
	}

	if err := config.DB.Delete(&krs).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to remove course from KRS")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Course removed from KRS successfully",
	})
}

func (kc *KRSController) GetAvailableCourses(c *gin.Context) {
	userID, _ := c.Get("user_id")
	semester, _ := strconv.Atoi(c.DefaultQuery("semester", "1"))
	tahunAjaran := c.DefaultQuery("tahun_ajaran", "")

	var courses []models.Course
	query := config.DB.Where("semester <= ?", semester)

	if err := query.Find(&courses).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch available courses")
		return
	}

	var availableCourses []models.Course
	for _, course := range courses {
		var existingKRS models.KRS
		if err := config.DB.Where("mahasiswa_id = ? AND course_id = ? AND semester = ? AND tahun_ajaran = ?",
			userID, course.ID, semester, tahunAjaran).First(&existingKRS).Error; err != nil {
			availableCourses = append(availableCourses, course)
		}
	}

	utils.SuccessResponse(c, availableCourses)
}
