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

type KajurController struct{}

func NewKajurController() *KajurController {
	return &KajurController{}
}

// Dashboard Kajur - Overview data jurusan
func (kc *KajurController) GetDashboard(c *gin.Context) {
	kajurID, _ := c.Get("user_id")

	// Ambil data kajur untuk mendapatkan jurusan
	var kajur models.Kajur
	if err := config.DB.Where("id = ?", kajurID).First(&kajur).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Kajur not found")
		return
	}

	// Statistik Mahasiswa
	var totalMahasiswa, mahasiswaAktif, mahasiswaCuti, mahasiswaDropOut int64
	config.DB.Model(&models.Mahasiswa{}).Where("jurusan = ?", kajur.Jurusan).Count(&totalMahasiswa)
	config.DB.Model(&models.Mahasiswa{}).Where("jurusan = ? AND status_akademik = 'aktif'", kajur.Jurusan).Count(&mahasiswaAktif)
	config.DB.Model(&models.Mahasiswa{}).Where("jurusan = ? AND status_akademik = 'cuti'", kajur.Jurusan).Count(&mahasiswaCuti)
	config.DB.Model(&models.Mahasiswa{}).Where("jurusan = ? AND status_akademik = 'drop_out'", kajur.Jurusan).Count(&mahasiswaDropOut)

	// Statistik Dosen
	var totalDosen, dosenAktif int64
	config.DB.Model(&models.Dosen{}).Where("jurusan = ?", kajur.Jurusan).Count(&totalDosen)
	config.DB.Model(&models.Dosen{}).Where("jurusan = ? AND status = 'aktif'", kajur.Jurusan).Count(&dosenAktif)

	// Statistik Mata Kuliah
	var totalMataKuliah, mataKuliahAktif int64
	config.DB.Table("courses").
		Joins("JOIN dosens ON courses.dosen_id = dosens.id").
		Where("dosens.jurusan = ?", kajur.Jurusan).
		Count(&totalMataKuliah)

	// Mata kuliah aktif (yang sedang ada jadwal di semester ini)
	config.DB.Table("courses").
		Joins("JOIN dosens ON courses.dosen_id = dosens.id").
		Joins("JOIN jadwals ON courses.id = jadwals.course_id").
		Where("dosens.jurusan = ? AND jadwals.tahun_ajaran = ?", kajur.Jurusan, getCurrentAcademicYear()).
		Distinct("courses.id").
		Count(&mataKuliahAktif)

	// IPK Rata-rata jurusan
	var avgIPK struct {
		Average float64
	}
	config.DB.Model(&models.Mahasiswa{}).
		Select("AVG(ipk) as average").
		Where("jurusan = ? AND status_akademik = 'aktif'", kajur.Jurusan).
		Scan(&avgIPK)

	// Tingkat kelulusan (mahasiswa lulus dibanding total alumni)
	var mahasiswaLulus int64
	config.DB.Model(&models.Mahasiswa{}).Where("jurusan = ? AND status_akademik = 'lulus'", kajur.Jurusan).Count(&mahasiswaLulus)

	tingkatKelulusan := 0.0
	totalAlumni := mahasiswaLulus + mahasiswaDropOut
	if totalAlumni > 0 {
		tingkatKelulusan = float64(mahasiswaLulus) / float64(totalAlumni) * 100
	}

	// Pending KRS yang perlu approval kajur
	var pendingKRS int64
	config.DB.Table("krs").
		Joins("JOIN mahasiswas ON krs.mahasiswa_id = mahasiswas.id").
		Where("mahasiswas.jurusan = ? AND krs.approval_status = 'pending'", kajur.Jurusan).
		Count(&pendingKRS)

	dashboard := models.KajurDashboardResponse{
		TotalMahasiswa:     int(totalMahasiswa),
		MahasiswaAktif:     int(mahasiswaAktif),
		MahasiswaCuti:      int(mahasiswaCuti),
		MahasiswaDropOut:   int(mahasiswaDropOut),
		TotalDosen:         int(totalDosen),
		DosenAktif:         int(dosenAktif),
		TotalMataKuliah:    int(totalMataKuliah),
		MataKuliahAktif:    int(mataKuliahAktif),
		IPKRataRata:        avgIPK.Average,
		TingkatKelulusan:   tingkatKelulusan,
		PendingKRSApproval: int(pendingKRS),
	}

	utils.SuccessResponse(c, dashboard)
}

// Lihat semua mahasiswa di jurusan
func (kc *KajurController) GetMahasiswaDiJurusan(c *gin.Context) {
	kajurID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	semester := c.Query("semester")
	statusAkademik := c.Query("status_akademik")
	search := c.Query("search")

	// Ambil data kajur untuk mendapatkan jurusan
	var kajur models.Kajur
	if err := config.DB.Where("id = ?", kajurID).First(&kajur).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Kajur not found")
		return
	}

	offset := (page - 1) * limit
	query := config.DB.Where("jurusan = ?", kajur.Jurusan)

	// Filter berdasarkan parameter
	if semester != "" {
		query = query.Where("semester = ?", semester)
	}
	if statusAkademik != "" {
		query = query.Where("status_akademik = ?", statusAkademik)
	}
	if search != "" {
		query = query.Where("nama ILIKE ? OR nim ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var mahasiswaList []models.Mahasiswa
	var total int64

	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Find(&mahasiswaList).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch mahasiswa")
		return
	}

	var responses []gin.H
	for _, mhs := range mahasiswaList {
		responses = append(responses, gin.H{
			"id":              mhs.ID,
			"nim":             mhs.NIM,
			"nama":            mhs.Nama,
			"jurusan":         mhs.Jurusan,
			"semester":        mhs.Semester,
			"status_akademik": mhs.StatusAkademik,
			"ipk":             mhs.IPK,
			"dosen_wali_id":   mhs.DosenWaliID,
			"created_at":      mhs.CreatedAt,
		})
	}

	utils.SuccessResponse(c, gin.H{
		"data": responses,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// Lihat semua dosen di jurusan
func (kc *KajurController) GetDosenDiJurusan(c *gin.Context) {
	kajurID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")
	search := c.Query("search")

	// Ambil data kajur untuk mendapatkan jurusan
	var kajur models.Kajur
	if err := config.DB.Where("id = ?", kajurID).First(&kajur).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Kajur not found")
		return
	}

	offset := (page - 1) * limit
	query := config.DB.Preload("Courses").Where("jurusan = ?", kajur.Jurusan)

	// Filter berdasarkan parameter
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		query = query.Where("nama ILIKE ? OR nidn ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	var dosenList []models.Dosen
	var total int64

	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Find(&dosenList).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch dosen")
		return
	}

	var responses []gin.H
	for _, dosen := range dosenList {
		// Hitung jumlah mahasiswa yang dibimbing
		var mahasiswaCount int64
		config.DB.Model(&models.Mahasiswa{}).Where("dosen_wali_id = ?", dosen.ID).Count(&mahasiswaCount)

		responses = append(responses, gin.H{
			"id":                 dosen.ID,
			"nidn":               dosen.NIDN,
			"nama":               dosen.Nama,
			"email":              dosen.Email,
			"jurusan":            dosen.Jurusan,
			"status":             dosen.Status,
			"jumlah_mata_kuliah": len(dosen.Courses),
			"jumlah_mahasiswa":   mahasiswaCount,
			"created_at":         dosen.CreatedAt,
		})
	}

	utils.SuccessResponse(c, gin.H{
		"data": responses,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// Approve/Validasi KRS mahasiswa di jurusan
func (kc *KajurController) GetPendingKRSValidation(c *gin.Context) {
	kajurID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	semester := c.Query("semester")
	tahunAjaran := c.Query("tahun_ajaran")

	// Ambil data kajur untuk mendapatkan jurusan
	var kajur models.Kajur
	if err := config.DB.Where("id = ?", kajurID).First(&kajur).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Kajur not found")
		return
	}

	offset := (page - 1) * limit
	query := config.DB.Preload("Mahasiswa").Preload("Course").
		Joins("JOIN mahasiswas ON krs.mahasiswa_id = mahasiswas.id").
		Where("mahasiswas.jurusan = ? AND krs.approval_status = 'pending'", kajur.Jurusan)

	if semester != "" {
		query = query.Where("krs.semester = ?", semester)
	}
	if tahunAjaran != "" {
		query = query.Where("krs.tahun_ajaran = ?", tahunAjaran)
	}

	var pendingKRS []models.KRS
	var total int64

	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Find(&pendingKRS).Error; err != nil {
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
				"ipk":      krs.Mahasiswa.IPK,
			},
			"course": gin.H{
				"id":      krs.Course.ID,
				"code":    krs.Course.Code,
				"name":    krs.Course.Name,
				"credits": krs.Course.Credits,
			},
			"semester":        krs.Semester,
			"tahun_ajaran":    krs.TahunAjaran,
			"approval_status": krs.ApprovalStatus,
			"created_at":      krs.CreatedAt,
		})
	}

	utils.SuccessResponse(c, gin.H{
		"data": responses,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// Validasi KRS (approve/reject) oleh Kajur
func (kc *KajurController) ProcessKRSValidation(c *gin.Context) {
	kajurID, _ := c.Get("user_id")
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

	// Ambil data kajur
	var kajur models.Kajur
	if err := config.DB.Where("id = ?", kajurID).First(&kajur).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Kajur not found")
		return
	}

	// Ambil KRS dan verifikasi mahasiswa di jurusan yang sama
	var krs models.KRS
	if err := config.DB.Preload("Mahasiswa").Where("id = ?", krsID).First(&krs).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "KRS not found")
		return
	}

	// Verifikasi mahasiswa di jurusan kajur
	if krs.Mahasiswa.Jurusan != kajur.Jurusan {
		utils.ErrorResponse(c, http.StatusForbidden, "You can only validate KRS for students in your department")
		return
	}

	// Update status approval
	now := time.Now()
	kajurIDUint := kajurID.(uint)
	if req.Action == "approve" {
		krs.ApprovalStatus = "approved"
		krs.Status = "diambil"
		krs.ApprovedBy = &kajurIDUint
		krs.ApprovedAt = &now
		krs.RejectionReason = ""
	} else {
		krs.ApprovalStatus = "rejected"
		krs.Status = "ditolak"
		krs.ApprovedBy = &kajurIDUint
		krs.ApprovedAt = &now
		krs.RejectionReason = req.RejectionReason
	}

	if err := config.DB.Save(&krs).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to process KRS validation")
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message": "KRS " + req.Action + "d successfully by department head",
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

// Monitoring nilai & absensi dosen di jurusan
func (kc *KajurController) GetMonitoringDosenPerformance(c *gin.Context) {
	kajurID, _ := c.Get("user_id")
	semester := c.DefaultQuery("semester", "")
	tahunAjaran := c.DefaultQuery("tahun_ajaran", getCurrentAcademicYear())

	// Ambil data kajur untuk mendapatkan jurusan
	var kajur models.Kajur
	if err := config.DB.Where("id = ?", kajurID).First(&kajur).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Kajur not found")
		return
	}

	// Ambil semua dosen di jurusan
	var dosenList []models.Dosen
	if err := config.DB.Preload("Courses").Where("jurusan = ? AND status = 'aktif'", kajur.Jurusan).Find(&dosenList).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch dosen")
		return
	}

	var monitoringData []gin.H
	for _, dosen := range dosenList {
		dosenData := gin.H{
			"dosen_id": dosen.ID,
			"nama":     dosen.Nama,
			"nidn":     dosen.NIDN,
			"email":    dosen.Email,
			"courses":  []gin.H{},
		}

		var totalNilaiInput, totalAbsensiInput int64
		var avgNilaiMahasiswa float64
		var coursesData []gin.H

		for _, course := range dosen.Courses {
			// Hitung berapa banyak nilai yang sudah diinput
			var nilaiCount int64
			query := config.DB.Model(&models.Nilai{}).Where("course_id = ?", course.ID)
			if semester != "" {
				query = query.Where("semester = ?", semester)
			}
			if tahunAjaran != "" {
				query = query.Where("tahun_ajaran = ?", tahunAjaran)
			}
			query.Count(&nilaiCount)
			totalNilaiInput += nilaiCount

			// Hitung berapa banyak absensi yang sudah diinput
			var absensiCount int64
			absensiQuery := config.DB.Model(&models.Absensi{}).Where("course_id = ?", course.ID)
			absensiQuery.Count(&absensiCount)
			totalAbsensiInput += absensiCount

			// Hitung rata-rata nilai mahasiswa di mata kuliah ini
			var avgNilai struct {
				Average float64
			}
			config.DB.Model(&models.Nilai{}).
				Select("AVG(nilai_akhir) as average").
				Where("course_id = ?", course.ID).
				Scan(&avgNilai)

			// Hitung jumlah mahasiswa terdaftar
			var mahasiswaCount int64
			config.DB.Model(&models.KRS{}).
				Where("course_id = ? AND approval_status = 'approved'", course.ID).
				Count(&mahasiswaCount)

			// Hitung tingkat kehadiran rata-rata
			var avgKehadiran struct {
				Percentage float64
			}
			config.DB.Raw(`
				SELECT 
				CASE 
					WHEN COUNT(*) > 0 THEN 
						CAST(SUM(CASE WHEN status = 'hadir' THEN 1 ELSE 0 END) AS FLOAT) / COUNT(*) * 100
					ELSE 0 
				END as percentage
				FROM absensis 
				WHERE course_id = ?
			`, course.ID).Scan(&avgKehadiran)

			courseData := gin.H{
				"course_id":        course.ID,
				"course_code":      course.Code,
				"course_name":      course.Name,
				"credits":          course.Credits,
				"jumlah_mahasiswa": mahasiswaCount,
				"nilai_terinput":   nilaiCount,
				"absensi_terinput": absensiCount,
				"rata_nilai":       avgNilai.Average,
				"rata_kehadiran":   avgKehadiran.Percentage,
			}
			coursesData = append(coursesData, courseData)
			avgNilaiMahasiswa += avgNilai.Average
		}

		// Hitung rata-rata keseluruhan dosen
		if len(dosen.Courses) > 0 {
			avgNilaiMahasiswa = avgNilaiMahasiswa / float64(len(dosen.Courses))
		}

		dosenData["courses"] = coursesData
		dosenData["total_mata_kuliah"] = len(dosen.Courses)
		dosenData["total_nilai_input"] = totalNilaiInput
		dosenData["total_absensi_input"] = totalAbsensiInput
		dosenData["avg_nilai_mahasiswa"] = avgNilaiMahasiswa

		monitoringData = append(monitoringData, dosenData)
	}

	utils.SuccessResponse(c, gin.H{
		"jurusan":      kajur.Jurusan,
		"semester":     semester,
		"tahun_ajaran": tahunAjaran,
		"total_dosen":  len(dosenList),
		"monitoring":   monitoringData,
	})
}

// Generate laporan jurusan
func (kc *KajurController) GenerateLaporanJurusan(c *gin.Context) {
	kajurID, _ := c.Get("user_id")
	tahunAjaran := c.DefaultQuery("tahun_ajaran", getCurrentAcademicYear())

	// Ambil data kajur untuk mendapatkan jurusan
	var kajur models.Kajur
	if err := config.DB.Where("id = ?", kajurID).First(&kajur).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Kajur not found")
		return
	}

	// Statistik Mahasiswa
	var totalMhs, aktifMhs, cutiMhs, dropOutMhs, lulusMhs int64
	config.DB.Model(&models.Mahasiswa{}).Where("jurusan = ?", kajur.Jurusan).Count(&totalMhs)
	config.DB.Model(&models.Mahasiswa{}).Where("jurusan = ? AND status_akademik = 'aktif'", kajur.Jurusan).Count(&aktifMhs)
	config.DB.Model(&models.Mahasiswa{}).Where("jurusan = ? AND status_akademik = 'cuti'", kajur.Jurusan).Count(&cutiMhs)
	config.DB.Model(&models.Mahasiswa{}).Where("jurusan = ? AND status_akademik = 'drop_out'", kajur.Jurusan).Count(&dropOutMhs)
	config.DB.Model(&models.Mahasiswa{}).Where("jurusan = ? AND status_akademik = 'lulus'", kajur.Jurusan).Count(&lulusMhs)

	// IPK Statistics
	var ipkStats struct {
		Average float64
		Max     float64
		Min     float64
	}
	config.DB.Model(&models.Mahasiswa{}).
		Select("AVG(ipk) as average, MAX(ipk) as max, MIN(ipk) as min").
		Where("jurusan = ? AND status_akademik = 'aktif'", kajur.Jurusan).
		Scan(&ipkStats)

	statistikMahasiswa := models.StatistikMahasiswaKajur{
		Total:        int(totalMhs),
		Aktif:        int(aktifMhs),
		Cuti:         int(cutiMhs),
		DropOut:      int(dropOutMhs),
		Lulus:        int(lulusMhs),
		IPKRataRata:  ipkStats.Average,
		IPKTertinggi: ipkStats.Max,
		IPKTerendah:  ipkStats.Min,
	}

	// Statistik Dosen
	var totalDosen, aktifDosen, nonAktifDosen int64
	config.DB.Model(&models.Dosen{}).Where("jurusan = ?", kajur.Jurusan).Count(&totalDosen)
	config.DB.Model(&models.Dosen{}).Where("jurusan = ? AND status = 'aktif'", kajur.Jurusan).Count(&aktifDosen)
	nonAktifDosen = totalDosen - aktifDosen

	// Rata-rata mata kuliah per dosen
	var avgMatkulPerDosen struct {
		Average float64
	}
	config.DB.Table("courses").
		Select("AVG(course_count) as average").
		Joins("JOIN (SELECT dosen_id, COUNT(*) as course_count FROM courses JOIN dosens ON courses.dosen_id = dosens.id WHERE dosens.jurusan = ? GROUP BY dosen_id) as course_counts ON true", kajur.Jurusan).
		Scan(&avgMatkulPerDosen)

	statistikDosen := models.StatistikDosenKajur{
		Total:              int(totalDosen),
		Aktif:              int(aktifDosen),
		NonAktif:           int(nonAktifDosen),
		RataMatkulPerDosen: avgMatkulPerDosen.Average,
	}

	// Statistik Akademik
	var totalMatkul, matkulAktif int64
	config.DB.Table("courses").
		Joins("JOIN dosens ON courses.dosen_id = dosens.id").
		Where("dosens.jurusan = ?", kajur.Jurusan).
		Count(&totalMatkul)

	config.DB.Table("courses").
		Joins("JOIN dosens ON courses.dosen_id = dosens.id").
		Joins("JOIN jadwals ON courses.id = jadwals.course_id").
		Where("dosens.jurusan = ? AND jadwals.tahun_ajaran = ?", kajur.Jurusan, tahunAjaran).
		Distinct("courses.id").
		Count(&matkulAktif)

	// Hitung total kelas (berdasarkan jadwal)
	var totalKelas int64
	config.DB.Table("jadwals").
		Joins("JOIN courses ON jadwals.course_id = courses.id").
		Joins("JOIN dosens ON courses.dosen_id = dosens.id").
		Where("dosens.jurusan = ? AND jadwals.tahun_ajaran = ?", kajur.Jurusan, tahunAjaran).
		Count(&totalKelas)

	// Rata-rata kehadiran mahasiswa
	var avgKehadiranMhs struct {
		Percentage float64
	}
	config.DB.Raw(`
		SELECT 
		CASE 
			WHEN COUNT(*) > 0 THEN 
				CAST(SUM(CASE WHEN absensis.status = 'hadir' THEN 1 ELSE 0 END) AS FLOAT) / COUNT(*) * 100
			ELSE 0 
		END as percentage
		FROM absensis 
		JOIN courses ON absensis.course_id = courses.id
		JOIN dosens ON courses.dosen_id = dosens.id
		WHERE dosens.jurusan = ?
	`, kajur.Jurusan).Scan(&avgKehadiranMhs)

	// Tingkat kelulusan
	tingkatKelulusan := 0.0
	totalAlumni := lulusMhs + dropOutMhs
	if totalAlumni > 0 {
		tingkatKelulusan = float64(lulusMhs) / float64(totalAlumni) * 100
	}

	statistikAkademik := models.StatistikAkademikKajur{
		TotalMataKuliah:    int(totalMatkul),
		MataKuliahAktif:    int(matkulAktif),
		TotalKelas:         int(totalKelas),
		RataKehadiranDosen: 85.0, // Placeholder - bisa dihitung dari jadwal vs kehadiran aktual
		RataKehadiranMhs:   avgKehadiranMhs.Percentage,
		TingkatKelulusan:   tingkatKelulusan,
	}

	// Performance Kelas
	var performanceKelas []models.PerformanceKelasKajur
	rows, err := config.DB.Raw(`
		SELECT 
			c.id as course_id,
			c.code as course_code,
			c.name as course_name,
			d.nama as dosen_pengampu,
			COUNT(DISTINCT k.mahasiswa_id) as jumlah_mhs,
			AVG(n.nilai_akhir) as rata_nilai,
			CASE 
				WHEN COUNT(n.id) > 0 THEN 
					CAST(SUM(CASE WHEN n.grade_huruf != 'E' THEN 1 ELSE 0 END) AS FLOAT) / COUNT(n.id) * 100
				ELSE 0 
			END as tingkat_lulus,
			CASE 
				WHEN COUNT(a.id) > 0 THEN 
					CAST(SUM(CASE WHEN a.status = 'hadir' THEN 1 ELSE 0 END) AS FLOAT) / COUNT(a.id) * 100
				ELSE 0 
			END as rata_kehadiran
		FROM courses c
		JOIN dosens d ON c.dosen_id = d.id
		LEFT JOIN krs k ON c.id = k.course_id AND k.approval_status = 'approved'
		LEFT JOIN nilais n ON c.id = n.course_id
		LEFT JOIN absensis a ON c.id = a.course_id
		WHERE d.jurusan = ?
		GROUP BY c.id, c.code, c.name, d.nama
		ORDER BY c.code
	`, kajur.Jurusan).Rows()

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var perf models.PerformanceKelasKajur
			rows.Scan(&perf.CourseID, &perf.CourseCode, &perf.CourseName, &perf.DosenPengampu,
				&perf.JumlahMhs, &perf.RataNilai, &perf.TingkatLulus, &perf.RataKehadiran)
			performanceKelas = append(performanceKelas, perf)
		}
	}

	// Trend Semester (data 3 semester terakhir)
	var trendSemester []models.TrendSemesterKajur
	// Ini bisa dikembangkan lebih lanjut dengan query yang lebih kompleks

	laporan := models.LaporanJurusanResponse{
		Jurusan:            kajur.Jurusan,
		PeriodeLaporan:     tahunAjaran,
		StatistikMahasiswa: statistikMahasiswa,
		StatistikDosen:     statistikDosen,
		StatistikAkademik:  statistikAkademik,
		PerformanceKelas:   performanceKelas,
		TrendSemester:      trendSemester,
		CreatedAt:          time.Now(),
	}

	utils.SuccessResponse(c, laporan)
}

// Manage mata kuliah - Lihat semua mata kuliah di jurusan
func (kc *KajurController) GetMataKuliahDiJurusan(c *gin.Context) {
	kajurID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	semester := c.Query("semester")
	search := c.Query("search")

	// Ambil data kajur untuk mendapatkan jurusan
	var kajur models.Kajur
	if err := config.DB.Where("id = ?", kajurID).First(&kajur).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Kajur not found")
		return
	}

	offset := (page - 1) * limit
	query := config.DB.Preload("Dosen").
		Joins("JOIN dosens ON courses.dosen_id = dosens.id").
		Where("dosens.jurusan = ?", kajur.Jurusan)

	if semester != "" {
		query = query.Where("courses.semester = ?", semester)
	}
	if search != "" {
		query = query.Where("courses.name ILIKE ? OR courses.code ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var courseList []models.Course
	var total int64

	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Find(&courseList).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch courses")
		return
	}

	var responses []gin.H
	for _, course := range courseList {
		// Hitung jumlah mahasiswa terdaftar
		var mahasiswaCount int64
		config.DB.Model(&models.KRS{}).
			Where("course_id = ? AND approval_status = 'approved'", course.ID).
			Count(&mahasiswaCount)

		// Cek apakah ada jadwal aktif
		var jadwalCount int64
		config.DB.Model(&models.Jadwal{}).
			Where("course_id = ? AND tahun_ajaran = ?", course.ID, getCurrentAcademicYear()).
			Count(&jadwalCount)

		dosenPengampu := ""
		if course.Dosen != nil {
			dosenPengampu = course.Dosen.Nama
		}

		statusKelas := "tidak_aktif"
		if jadwalCount > 0 {
			statusKelas = "aktif"
		}

		responses = append(responses, gin.H{
			"id":               course.ID,
			"code":             course.Code,
			"name":             course.Name,
			"credits":          course.Credits,
			"semester":         course.Semester,
			"prasyarat":        course.Prasyarat,
			"deskripsi":        course.Deskripsi,
			"dosen_pengampu":   dosenPengampu,
			"jumlah_mahasiswa": mahasiswaCount,
			"status_kelas":     statusKelas,
			"created_at":       course.CreatedAt,
		})
	}

	utils.SuccessResponse(c, gin.H{
		"data": responses,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// Buka/Tutup Kelas baru - Update status mata kuliah
func (kc *KajurController) UpdateStatusMataKuliah(c *gin.Context) {
	kajurID, _ := c.Get("user_id")
	courseID := c.Param("courseId")

	var req struct {
		Action string `json:"action" validate:"required,oneof=buka tutup"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// Ambil data kajur
	var kajur models.Kajur
	if err := config.DB.Where("id = ?", kajurID).First(&kajur).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Kajur not found")
		return
	}

	// Ambil mata kuliah dan verifikasi di jurusan yang sama
	var course models.Course
	if err := config.DB.Preload("Dosen").Where("id = ?", courseID).First(&course).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Course not found")
		return
	}

	// Verifikasi mata kuliah di jurusan kajur
	if course.Dosen == nil || course.Dosen.Jurusan != kajur.Jurusan {
		utils.ErrorResponse(c, http.StatusForbidden, "You can only manage courses in your department")
		return
	}

	if req.Action == "buka" {
		// Cek apakah sudah ada jadwal untuk tahun ajaran ini
		var jadwal models.Jadwal
		if err := config.DB.Where("course_id = ? AND tahun_ajaran = ?", courseID, getCurrentAcademicYear()).First(&jadwal).Error; err != nil {
			// Belum ada jadwal, buat jadwal default
			newJadwal := models.Jadwal{
				CourseID:    course.ID,
				Hari:        "senin", // Default, bisa diubah nanti
				JamMulai:    "08:00",
				JamSelesai:  "10:00",
				Ruangan:     "TBD",
				Dosen:       course.Dosen.Nama,
				TipeKelas:   "kuliah",
				Semester:    course.Semester,
				TahunAjaran: getCurrentAcademicYear(),
			}

			if err := config.DB.Create(&newJadwal).Error; err != nil {
				utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create schedule")
				return
			}
		}

		utils.SuccessResponse(c, gin.H{
			"message": "Kelas berhasil dibuka untuk mata kuliah " + course.Name,
			"status":  "aktif",
		})
	} else {
		// Tutup kelas - hapus jadwal untuk tahun ajaran ini
		if err := config.DB.Where("course_id = ? AND tahun_ajaran = ?", courseID, getCurrentAcademicYear()).Delete(&models.Jadwal{}).Error; err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to close class")
			return
		}

		utils.SuccessResponse(c, gin.H{
			"message": "Kelas berhasil ditutup untuk mata kuliah " + course.Name,
			"status":  "tidak_aktif",
		})
	}
}

// Helper function untuk mendapatkan tahun ajaran saat ini
func getCurrentAcademicYear() string {
	now := time.Now()
	year := now.Year()

	// Jika bulan Juli-Desember, maka tahun ajaran dimulai tahun ini
	// Jika bulan Januari-Juni, maka tahun ajaran dimulai tahun lalu
	if now.Month() >= 7 {
		return strconv.Itoa(year) + "/" + strconv.Itoa(year+1)
	}
	return strconv.Itoa(year-1) + "/" + strconv.Itoa(year)
}
