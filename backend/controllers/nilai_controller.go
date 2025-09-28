package controllers

import (
	"SIAku/config"
	"SIAku/models"
	"SIAku/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NilaiController struct{}

func NewNilaiController() *NilaiController {
	return &NilaiController{}
}

func (nc *NilaiController) GetMyNilai(c *gin.Context) {
	userID, _ := c.Get("user_id")
	semester := c.DefaultQuery("semester", "")
	tahunAjaran := c.DefaultQuery("tahun_ajaran", "")

	var nilai []models.Nilai
	query := config.DB.Preload("Course").Where("mahasiswa_id = ?", userID)
	
	if semester != "" {
		query = query.Where("semester = ?", semester)
	}
	if tahunAjaran != "" {
		query = query.Where("tahun_ajaran = ?", tahunAjaran)
	}

	if err := query.Find(&nilai).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch nilai")
		return
	}

	var responses []models.NilaiResponse
	for _, n := range nilai {
		responses = append(responses, models.NilaiResponse{
			ID:          n.ID,
			CourseName:  n.Course.Name,
			CourseCode:  n.Course.Code,
			Credits:     n.Course.Credits,
			Semester:    n.Semester,
			TahunAjaran: n.TahunAjaran,
			NilaiTugas:  n.NilaiTugas,
			NilaiUTS:    n.NilaiUTS,
			NilaiUAS:    n.NilaiUAS,
			NilaiAkhir:  n.NilaiAkhir,
			GradeHuruf:  n.GradeHuruf,
			GradePoint:  n.GradePoint,
			Status:      n.Status,
			CreatedAt:   n.CreatedAt,
		})
	}

	utils.SuccessResponse(c, responses)
}

func (nc *NilaiController) GetTranskrip(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var mahasiswa models.Mahasiswa
	if err := config.DB.Where("id = ?", userID).First(&mahasiswa).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Mahasiswa not found")
		return
	}

	var nilaiList []models.Nilai
	if err := config.DB.Preload("Course").Where("mahasiswa_id = ?", userID).Find(&nilaiList).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch transkrip")
		return
	}

	mahasiswaResponse := models.MahasiswaResponse{
		ID:        mahasiswa.ID,
		NIM:       mahasiswa.NIM,
		Nama:      mahasiswa.Nama,
		Jurusan:   mahasiswa.Jurusan,
		CreatedAt: mahasiswa.CreatedAt,
		UpdatedAt: mahasiswa.UpdatedAt,
	}

	var riwayatNilai []models.NilaiResponse
	totalSKS := 0
	totalSKSLulus := 0
	totalPoin := 0.0

	for _, nilai := range nilaiList {
		totalSKS += nilai.Course.Credits
		if nilai.GradeHuruf != "E" {
			totalSKSLulus += nilai.Course.Credits
		}
		totalPoin += nilai.GradePoint * float64(nilai.Course.Credits)
		
		nilaiResp := models.NilaiResponse{
			ID:          nilai.ID,
			CourseName:  nilai.Course.Name,
			CourseCode:  nilai.Course.Code,
			Credits:     nilai.Course.Credits,
			Semester:    nilai.Semester,
			TahunAjaran: nilai.TahunAjaran,
			NilaiTugas:  nilai.NilaiTugas,
			NilaiUTS:    nilai.NilaiUTS,
			NilaiUAS:    nilai.NilaiUAS,
			NilaiAkhir:  nilai.NilaiAkhir,
			GradeHuruf:  nilai.GradeHuruf,
			GradePoint:  nilai.GradePoint,
			Status:      nilai.Status,
			CreatedAt:   nilai.CreatedAt,
		}
		
		riwayatNilai = append(riwayatNilai, nilaiResp)
	}

	ipk := 0.0
	if totalSKS > 0 {
		ipk = totalPoin / float64(totalSKS)
	}

	statusKelulusan := "Aktif"
	if mahasiswa.StatusAkademik == "lulus" {
		statusKelulusan = "Lulus"
	} else if mahasiswa.StatusAkademik == "drop_out" {
		statusKelulusan = "Drop Out"
	}

	transkrip := models.TranskripResponse{
		Mahasiswa:       mahasiswaResponse,
		TotalSKS:        totalSKS,
		TotalSKSLulus:   totalSKSLulus,
		IPKKumulatif:    ipk,
		RiwayatNilai:    riwayatNilai,
		StatusKelulusan: statusKelulusan,
	}

	utils.SuccessResponse(c, transkrip)
}

func (nc *NilaiController) GetStatistikNilai(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var nilaiList []models.Nilai
	if err := config.DB.Preload("Course").Where("mahasiswa_id = ?", userID).Find(&nilaiList).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch statistik nilai")
		return
	}

	statistik := map[string]interface{}{
		"total_mata_kuliah": len(nilaiList),
		"grade_distribution": map[string]int{
			"A":  0,
			"AB": 0,
			"B":  0,
			"BC": 0,
			"C":  0,
			"D":  0,
			"E":  0,
		},
		"rata_rata_nilai": 0.0,
		"total_poin":     0.0,
		"total_sks":      0,
	}

	totalNilai := 0.0
	totalPoin := 0.0
	totalSKS := 0
	gradeCount := map[string]int{
		"A": 0, "AB": 0, "B": 0, "BC": 0, "C": 0, "D": 0, "E": 0,
	}

	for _, nilai := range nilaiList {
		totalNilai += nilai.NilaiAkhir
		totalPoin += nilai.GradePoint * float64(nilai.Course.Credits)
		totalSKS += nilai.Course.Credits
		gradeCount[nilai.GradeHuruf]++
	}

	if len(nilaiList) > 0 {
		statistik["rata_rata_nilai"] = totalNilai / float64(len(nilaiList))
	}
	statistik["total_poin"] = totalPoin
	statistik["total_sks"] = totalSKS
	statistik["grade_distribution"] = gradeCount

	utils.SuccessResponse(c, statistik)
}