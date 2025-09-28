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

type RektorController struct{}

func NewRektorController() *RektorController {
	return &RektorController{}
}

// GetUniversityDashboard - Dashboard summary seluruh universitas
func (rc *RektorController) GetUniversityDashboard(c *gin.Context) {
	userRole := c.GetString("user_role")
	if userRole != "rektor" {
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied: Rektor role required")
		return
	}

	semester := c.DefaultQuery("semester", "ganjil_2024_2025")
	academicYear := c.DefaultQuery("academic_year", "2024-2025")

	// Get university-wide statistics
	var totalStudents int64
	var totalLecturers int64
	var totalDepartments int64

	config.DB.Model(&models.Mahasiswa{}).Count(&totalStudents)
	config.DB.Model(&models.Dosen{}).Count(&totalLecturers)

	// Count unique departments from both dosen and mahasiswa
	var uniqueDepartments []string
	config.DB.Model(&models.Mahasiswa{}).Distinct("jurusan").Pluck("jurusan", &uniqueDepartments)
	totalDepartments = int64(len(uniqueDepartments))

	// Calculate overall GPA
	var avgGPA float64
	config.DB.Model(&models.Nilai{}).Select("AVG(CASE WHEN nilai >= 85 THEN 4.0 WHEN nilai >= 80 THEN 3.7 WHEN nilai >= 75 THEN 3.3 WHEN nilai >= 70 THEN 3.0 WHEN nilai >= 65 THEN 2.7 WHEN nilai >= 60 THEN 2.3 WHEN nilai >= 55 THEN 2.0 WHEN nilai >= 50 THEN 1.7 WHEN nilai >= 45 THEN 1.3 WHEN nilai >= 40 THEN 1.0 ELSE 0.0 END) as avg_gpa").Scan(&avgGPA)

	// Calculate graduation rate (simplified)
	var totalGraduated int64
	var totalEligible int64
	config.DB.Model(&models.Mahasiswa{}).Where("created_at <= ?", time.Now().AddDate(-4, 0, 0)).Count(&totalEligible)
	if totalEligible > 0 {
		totalGraduated = totalEligible * 85 / 100 // Simulated 85% graduation rate
	}
	graduationRate := float64(totalGraduated) / float64(totalEligible) * 100

	// Build faculty reports
	var facultyReports []models.FacultyReport
	for i, dept := range uniqueDepartments {
		var deptStudents int64
		var deptLecturers int64
		var deptGPA float64

		config.DB.Model(&models.Mahasiswa{}).Where("jurusan = ?", dept).Count(&deptStudents)
		config.DB.Model(&models.Dosen{}).Where("jurusan = ?", dept).Count(&deptLecturers)

		// Get department GPA
		config.DB.Table("nilais").
			Joins("JOIN mahasiswas ON mahasiswas.id = nilais.mahasiswa_id").
			Where("mahasiswas.jurusan = ?", dept).
			Select("AVG(CASE WHEN nilai >= 85 THEN 4.0 WHEN nilai >= 80 THEN 3.7 WHEN nilai >= 75 THEN 3.3 WHEN nilai >= 70 THEN 3.0 WHEN nilai >= 65 THEN 2.7 WHEN nilai >= 60 THEN 2.3 WHEN nilai >= 55 THEN 2.0 WHEN nilai >= 50 THEN 1.7 WHEN nilai >= 45 THEN 1.3 WHEN nilai >= 40 THEN 1.0 ELSE 0.0 END) as avg_gpa").
			Scan(&deptGPA)

		deptSummary := models.DepartmentSummary{
			DepartmentID:   uint(i + 1),
			DepartmentName: dept,
			TotalStudents:  int(deptStudents),
			AverageGPA:     deptGPA,
			GraduationRate: 85.0 + float64(i*2), // Simulated varying rates
			Accreditation:  "A",
			HeadOfDept:     "Prof. Head " + dept,
		}

		facultyReport := models.FacultyReport{
			FacultyID:        uint(i + 1),
			FacultyName:      "Faculty of " + dept,
			TotalStudents:    int(deptStudents),
			TotalLecturers:   int(deptLecturers),
			TotalDepartments: 1,
			AverageGPA:       deptGPA,
			GraduationRate:   85.0 + float64(i*2),
			Departments:      []models.DepartmentSummary{deptSummary},
			Performance: models.FacultyPerformance{
				Rank:                i + 1,
				ResearchOutput:      50 + i*10,
				PublicationCount:    20 + i*5,
				StudentSatisfaction: 4.2 + float64(i)*0.1,
				EmployabilityRate:   88.0 + float64(i)*2,
			},
		}
		facultyReports = append(facultyReports, facultyReport)
	}

	// Top performing departments
	var topPerformers []models.TopPerformingDepartment
	for i, report := range facultyReports {
		if i < 3 { // Top 3
			topPerformer := models.TopPerformingDepartment{
				DepartmentName: report.Departments[0].DepartmentName,
				FacultyName:    report.FacultyName,
				AverageGPA:     report.AverageGPA,
				GraduationRate: report.GraduationRate,
				Accreditation:  "A",
			}
			topPerformers = append(topPerformers, topPerformer)
		}
	}

	// University alerts
	alerts := []models.UniversityAlert{
		{
			AlertType: "academic",
			Title:     "Low GPA Alert",
			Message:   "Some departments showing declining GPA trends",
			Severity:  "medium",
			CreatedAt: time.Now().AddDate(0, 0, -1),
		},
		{
			AlertType: "administrative",
			Title:     "Accreditation Review",
			Message:   "Upcoming accreditation review for 3 departments",
			Severity:  "high",
			CreatedAt: time.Now().AddDate(0, 0, -3),
		},
	}

	dashboard := models.RektorDashboardResponse{
		UniversityName: "Universitas Example",
		AcademicYear:   academicYear,
		Semester:       semester,
		Summary: models.UniversitySummary{
			TotalStudents:    int(totalStudents),
			TotalLecturers:   int(totalLecturers),
			TotalDepartments: int(totalDepartments),
			TotalFaculties:   len(facultyReports),
			OverallGPA:       avgGPA,
			GraduationRate:   graduationRate,
			StudentRetention: 92.5,
			AccreditationA:   len(facultyReports),
			AccreditationB:   0,
			AccreditationC:   0,
		},
		Faculties:     facultyReports,
		TopPerformers: topPerformers,
		Alerts:        alerts,
		GeneratedAt:   time.Now(),
	}

	utils.SuccessResponse(c, dashboard)
}

// GetFacultyReport - Lihat laporan per fakultas/jurusan
func (rc *RektorController) GetFacultyReport(c *gin.Context) {
	userRole := c.GetString("user_role")
	if userRole != "rektor" {
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied: Rektor role required")
		return
	}

	facultyName := c.Param("faculty")
	if facultyName == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Faculty name is required")
		return
	}

	// For simplicity, treating jurusan as faculty
	var totalStudents int64
	var totalLecturers int64
	var avgGPA float64

	config.DB.Model(&models.Mahasiswa{}).Where("jurusan = ?", facultyName).Count(&totalStudents)
	config.DB.Model(&models.Dosen{}).Where("jurusan = ?", facultyName).Count(&totalLecturers)

	config.DB.Table("nilais").
		Joins("JOIN mahasiswas ON mahasiswas.id = nilais.mahasiswa_id").
		Where("mahasiswas.jurusan = ?", facultyName).
		Select("AVG(CASE WHEN nilai >= 85 THEN 4.0 WHEN nilai >= 80 THEN 3.7 WHEN nilai >= 75 THEN 3.3 WHEN nilai >= 70 THEN 3.0 WHEN nilai >= 65 THEN 2.7 WHEN nilai >= 60 THEN 2.3 WHEN nilai >= 55 THEN 2.0 WHEN nilai >= 50 THEN 1.7 WHEN nilai >= 45 THEN 1.3 WHEN nilai >= 40 THEN 1.0 ELSE 0.0 END) as avg_gpa").
		Scan(&avgGPA)

	deptSummary := models.DepartmentSummary{
		DepartmentID:   1,
		DepartmentName: facultyName,
		TotalStudents:  int(totalStudents),
		AverageGPA:     avgGPA,
		GraduationRate: 87.5,
		Accreditation:  "A",
		HeadOfDept:     "Prof. Head " + facultyName,
	}

	facultyReport := models.FacultyReport{
		FacultyID:        1,
		FacultyName:      "Faculty of " + facultyName,
		TotalStudents:    int(totalStudents),
		TotalLecturers:   int(totalLecturers),
		TotalDepartments: 1,
		AverageGPA:       avgGPA,
		GraduationRate:   87.5,
		Departments:      []models.DepartmentSummary{deptSummary},
		Performance: models.FacultyPerformance{
			Rank:                1,
			ResearchOutput:      75,
			PublicationCount:    30,
			StudentSatisfaction: 4.3,
			EmployabilityRate:   90.0,
		},
	}

	utils.SuccessResponse(c, facultyReport)
}

// AssignRole - Mengatur role (otorisasi global)
func (rc *RektorController) AssignRole(c *gin.Context) {
	userRole := c.GetString("user_role")
	if userRole != "rektor" {
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied: Rektor role required")
		return
	}

	var req models.RoleAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// Find the user to be assigned
	var dosen models.Dosen
	if err := config.DB.Where("id = ?", req.UserID).First(&dosen).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	oldRole := "dosen" // Default role from context

	// Update role based on new assignment
	if req.NewRole == "kajur" {
		// Create Kajur record
		kajur := models.Kajur{
			NIDN:    dosen.NIDN,
			Nama:    dosen.Nama,
			Email:   dosen.Email,
			Jurusan: req.Department,
			Status:  "active",
		}

		if err := config.DB.Create(&kajur).Error; err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create kajur role")
			return
		}

		// Update dosen status
		config.DB.Model(&dosen).Update("status", "promoted_to_kajur")
	} else {
		// Update existing dosen status
		config.DB.Model(&dosen).Update("status", "active")
	}

	effectiveDate, _ := time.Parse("2006-01-02", req.EffectiveDate)
	assignedBy := c.GetString("user_name")

	response := models.RoleAssignmentResponse{
		AssignmentID:  uint(time.Now().Unix()),
		UserID:        req.UserID,
		UserName:      dosen.Nama,
		UserNIDN:      dosen.NIDN,
		OldRole:       oldRole,
		NewRole:       req.NewRole,
		Department:    req.Department,
		Faculty:       req.Faculty,
		EffectiveDate: effectiveDate,
		AssignedBy:    assignedBy,
		Reason:        req.Reason,
		Status:        "active",
		CreatedAt:     time.Now(),
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Role assigned successfully",
		"data":    response,
	})
}

// ApprovePolicy - Approve/reject kebijakan besar
func (rc *RektorController) ApprovePolicy(c *gin.Context) {
	userRole := c.GetString("user_role")
	if userRole != "rektor" {
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied: Rektor role required")
		return
	}

	policyIDStr := c.Param("id")
	policyID, err := strconv.ParseUint(policyIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	var req models.PolicyApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	req.PolicyID = uint(policyID)

	if err := utils.ValidateStruct(req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	approvalDate, _ := time.Parse("2006-01-02", req.ApprovalDate)
	approvedBy := c.GetString("user_name")

	// Simulate policy approval (in real app, would update policy table)
	response := models.PolicyApprovalResponse{
		PolicyID:     req.PolicyID,
		PolicyTitle:  "Academic Policy Update " + policyIDStr,
		SubmittedBy:  "Kajur Teknik Informatika",
		Department:   "Teknik Informatika",
		Faculty:      "Faculty of Engineering",
		Action:       req.Action,
		ApprovedBy:   approvedBy,
		Comments:     req.Comments,
		ApprovalDate: approvalDate,
		Status:       req.Action + "d",
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Policy " + req.Action + "d successfully",
		"data":    response,
	})
}

// GetPendingPolicies - View policies awaiting approval
func (rc *RektorController) GetPendingPolicies(c *gin.Context) {
	userRole := c.GetString("user_role")
	if userRole != "rektor" {
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied: Rektor role required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	// Simulate pending policies (in real app, would query policy table)
	policies := []models.PolicyApprovalResponse{
		{
			PolicyID:    1,
			PolicyTitle: "New Academic Curriculum 2025",
			SubmittedBy: "Kajur Teknik Informatika",
			Department:  "Teknik Informatika",
			Faculty:     "Faculty of Engineering",
			Status:      "pending",
		},
		{
			PolicyID:    2,
			PolicyTitle: "Student Exchange Program Guidelines",
			SubmittedBy: "Kajur Manajemen",
			Department:  "Manajemen",
			Faculty:     "Faculty of Economics",
			Status:      "pending",
		},
		{
			PolicyID:    3,
			PolicyTitle: "Research Grant Allocation Policy",
			SubmittedBy: "Kajur Teknik Sipil",
			Department:  "Teknik Sipil",
			Faculty:     "Faculty of Engineering",
			Status:      "pending",
		},
	}

	// Apply pagination
	start := offset
	end := offset + limit
	if start > len(policies) {
		start = len(policies)
	}
	if end > len(policies) {
		end = len(policies)
	}

	paginatedPolicies := policies[start:end]

	response := gin.H{
		"policies": paginatedPolicies,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       len(policies),
			"total_pages": (len(policies) + limit - 1) / limit,
		},
	}

	utils.SuccessResponse(c, response)
}

// GetRoleAssignments - View current role assignments
func (rc *RektorController) GetRoleAssignments(c *gin.Context) {
	userRole := c.GetString("user_role")
	if userRole != "rektor" {
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied: Rektor role required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var dosens []models.Dosen
	var total int64

	query := config.DB.Model(&models.Dosen{})
	query.Count(&total)
	query.Offset(offset).Limit(limit).Find(&dosens)

	var assignments []models.RoleAssignmentResponse
	for _, dosen := range dosens {
		assignment := models.RoleAssignmentResponse{
			UserID:     dosen.ID,
			UserName:   dosen.Nama,
			UserNIDN:   dosen.NIDN,
			NewRole:    "dosen",
			Department: dosen.Jurusan,
			Status:     dosen.Status,
			CreatedAt:  dosen.CreatedAt,
		}
		assignments = append(assignments, assignment)
	}

	response := gin.H{
		"assignments": assignments,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       int(total),
			"total_pages": (int(total) + limit - 1) / limit,
		},
	}

	utils.SuccessResponse(c, response)
}
