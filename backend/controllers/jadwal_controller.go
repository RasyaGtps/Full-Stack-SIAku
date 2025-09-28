package controllers

import (
	"SIAku/config"
	"SIAku/models"
	"SIAku/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type JadwalController struct{}

func NewJadwalController() *JadwalController {
	return &JadwalController{}
}

func (jc *JadwalController) GetMyJadwal(c *gin.Context) {
	userID, _ := c.Get("user_id")
	semester := c.DefaultQuery("semester", "")
	tahunAjaran := c.DefaultQuery("tahun_ajaran", "")
	hari := c.DefaultQuery("hari", "")

	var krs []models.KRS
	krsQuery := config.DB.Where("mahasiswa_id = ?", userID)

	if semester != "" {
		krsQuery = krsQuery.Where("semester = ?", semester)
	}
	if tahunAjaran != "" {
		krsQuery = krsQuery.Where("tahun_ajaran = ?", tahunAjaran)
	}

	if err := krsQuery.Find(&krs).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch KRS")
		return
	}

	var courseIDs []uint
	for _, k := range krs {
		courseIDs = append(courseIDs, k.CourseID)
	}

	if len(courseIDs) == 0 {
		utils.SuccessResponse(c, []models.JadwalResponse{})
		return
	}

	var jadwal []models.Jadwal
	jadwalQuery := config.DB.Preload("Course").Where("course_id IN ?", courseIDs)

	if hari != "" {
		jadwalQuery = jadwalQuery.Where("LOWER(hari) = LOWER(?)", hari)
	}

	if err := jadwalQuery.Find(&jadwal).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch jadwal")
		return
	}

	var responses []models.JadwalResponse
	for _, j := range jadwal {
		responses = append(responses, models.JadwalResponse{
			ID:          j.ID,
			CourseName:  j.Course.Name,
			CourseCode:  j.Course.Code,
			Credits:     j.Course.Credits,
			Hari:        j.Hari,
			JamMulai:    j.JamMulai,
			JamSelesai:  j.JamSelesai,
			Ruangan:     j.Ruangan,
			Dosen:       j.Dosen,
			TipeKelas:   j.TipeKelas,
			Semester:    j.Semester,
			TahunAjaran: j.TahunAjaran,
			CreatedAt:   j.CreatedAt,
		})
	}

	utils.SuccessResponse(c, responses)
}

func (jc *JadwalController) GetJadwalByHari(c *gin.Context) {
	userID, _ := c.Get("user_id")
	hari := strings.ToLower(c.Param("hari"))

	validHari := map[string]bool{
		"senin": true, "selasa": true, "rabu": true,
		"kamis": true, "jumat": true, "sabtu": true, "minggu": true,
	}

	if !validHari[hari] {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid day. Use: senin, selasa, rabu, kamis, jumat, sabtu, minggu")
		return
	}

	var krs []models.KRS
	if err := config.DB.Where("mahasiswa_id = ?", userID).Find(&krs).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch KRS")
		return
	}

	var courseIDs []uint
	for _, k := range krs {
		courseIDs = append(courseIDs, k.CourseID)
	}

	if len(courseIDs) == 0 {
		utils.SuccessResponse(c, []models.JadwalResponse{})
		return
	}

	var jadwal []models.Jadwal
	if err := config.DB.Preload("Course").
		Where("course_id IN ? AND LOWER(hari) = ?", courseIDs, hari).
		Order("jam_mulai ASC").
		Find(&jadwal).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch jadwal")
		return
	}

	var responses []models.JadwalResponse
	for _, j := range jadwal {
		responses = append(responses, models.JadwalResponse{
			ID:          j.ID,
			CourseName:  j.Course.Name,
			CourseCode:  j.Course.Code,
			Credits:     j.Course.Credits,
			Hari:        j.Hari,
			JamMulai:    j.JamMulai,
			JamSelesai:  j.JamSelesai,
			Ruangan:     j.Ruangan,
			Dosen:       j.Dosen,
			TipeKelas:   j.TipeKelas,
			Semester:    j.Semester,
			TahunAjaran: j.TahunAjaran,
			CreatedAt:   j.CreatedAt,
		})
	}

	utils.SuccessResponse(c, responses)
}

func (jc *JadwalController) GetJadwalMingguIni(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var krs []models.KRS
	if err := config.DB.Where("mahasiswa_id = ?", userID).Find(&krs).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch KRS")
		return
	}

	var courseIDs []uint
	for _, k := range krs {
		courseIDs = append(courseIDs, k.CourseID)
	}

	if len(courseIDs) == 0 {
		utils.SuccessResponse(c, map[string][]models.JadwalResponse{})
		return
	}

	var jadwal []models.Jadwal
	if err := config.DB.Preload("Course").
		Where("course_id IN ?", courseIDs).
		Order("CASE WHEN LOWER(hari) = 'senin' THEN 1 WHEN LOWER(hari) = 'selasa' THEN 2 WHEN LOWER(hari) = 'rabu' THEN 3 WHEN LOWER(hari) = 'kamis' THEN 4 WHEN LOWER(hari) = 'jumat' THEN 5 WHEN LOWER(hari) = 'sabtu' THEN 6 WHEN LOWER(hari) = 'minggu' THEN 7 END, jam_mulai ASC").
		Find(&jadwal).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch jadwal")
		return
	}

	jadwalMinggu := make(map[string][]models.JadwalResponse)

	for _, j := range jadwal {
		hari := strings.Title(strings.ToLower(j.Hari))

		response := models.JadwalResponse{
			ID:          j.ID,
			CourseName:  j.Course.Name,
			CourseCode:  j.Course.Code,
			Credits:     j.Course.Credits,
			Hari:        j.Hari,
			JamMulai:    j.JamMulai,
			JamSelesai:  j.JamSelesai,
			Ruangan:     j.Ruangan,
			Dosen:       j.Dosen,
			TipeKelas:   j.TipeKelas,
			Semester:    j.Semester,
			TahunAjaran: j.TahunAjaran,
			CreatedAt:   j.CreatedAt,
		}

		jadwalMinggu[hari] = append(jadwalMinggu[hari], response)
	}

	utils.SuccessResponse(c, jadwalMinggu)
}
