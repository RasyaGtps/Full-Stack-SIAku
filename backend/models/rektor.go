package models

import "time"

type Rektor struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	NIDN      string    `gorm:"unique;not null" json:"nidn" validate:"required,min=8,max=20"`
	Nama      string    `gorm:"type:varchar(100);not null" json:"nama" validate:"required,min=2,max=100"`
	Email     string    `gorm:"type:varchar(100);unique;not null" json:"email" validate:"required,email"`
	Status    string    `gorm:"type:varchar(20);default:'aktif'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RektorRequest struct {
	NIDN  string `json:"nidn" validate:"required,min=8,max=20"`
	Nama  string `json:"nama" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
}

type RektorResponse struct {
	ID        uint      `json:"id"`
	NIDN      string    `json:"nidn"`
	Nama      string    `json:"nama"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Dashboard Response for Rektor - University-wide summary
type RektorDashboardResponse struct {
	UniversityName string                    `json:"university_name"`
	AcademicYear   string                    `json:"academic_year"`
	Semester       string                    `json:"semester"`
	Summary        UniversitySummary         `json:"summary"`
	Faculties      []FacultyReport           `json:"faculties"`
	TopPerformers  []TopPerformingDepartment `json:"top_performers"`
	Alerts         []UniversityAlert         `json:"alerts"`
	GeneratedAt    time.Time                 `json:"generated_at"`
}

type UniversitySummary struct {
	TotalStudents    int     `json:"total_students"`
	TotalLecturers   int     `json:"total_lecturers"`
	TotalDepartments int     `json:"total_departments"`
	TotalFaculties   int     `json:"total_faculties"`
	OverallGPA       float64 `json:"overall_gpa"`
	GraduationRate   float64 `json:"graduation_rate"`
	StudentRetention float64 `json:"student_retention"`
	AccreditationA   int     `json:"accreditation_a"`
	AccreditationB   int     `json:"accreditation_b"`
	AccreditationC   int     `json:"accreditation_c"`
}

type FacultyReport struct {
	FacultyID        uint                `json:"faculty_id"`
	FacultyName      string              `json:"faculty_name"`
	TotalStudents    int                 `json:"total_students"`
	TotalLecturers   int                 `json:"total_lecturers"`
	TotalDepartments int                 `json:"total_departments"`
	AverageGPA       float64             `json:"average_gpa"`
	GraduationRate   float64             `json:"graduation_rate"`
	Departments      []DepartmentSummary `json:"departments"`
	Performance      FacultyPerformance  `json:"performance"`
}

type DepartmentSummary struct {
	DepartmentID   uint    `json:"department_id"`
	DepartmentName string  `json:"department_name"`
	TotalStudents  int     `json:"total_students"`
	AverageGPA     float64 `json:"average_gpa"`
	GraduationRate float64 `json:"graduation_rate"`
	Accreditation  string  `json:"accreditation"`
	HeadOfDept     string  `json:"head_of_department"`
}

type FacultyPerformance struct {
	Rank                int     `json:"rank"`
	ResearchOutput      int     `json:"research_output"`
	PublicationCount    int     `json:"publication_count"`
	StudentSatisfaction float64 `json:"student_satisfaction"`
	EmployabilityRate   float64 `json:"employability_rate"`
}

type TopPerformingDepartment struct {
	DepartmentName string  `json:"department_name"`
	FacultyName    string  `json:"faculty_name"`
	AverageGPA     float64 `json:"average_gpa"`
	GraduationRate float64 `json:"graduation_rate"`
	Accreditation  string  `json:"accreditation"`
}

type UniversityAlert struct {
	AlertType  string    `json:"alert_type"`
	Title      string    `json:"title"`
	Message    string    `json:"message"`
	Severity   string    `json:"severity"`
	Department string    `json:"department,omitempty"`
	Faculty    string    `json:"faculty,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// Role Management Structures
type RoleAssignmentRequest struct {
	UserID        uint   `json:"user_id" validate:"required"`
	UserType      string `json:"user_type" validate:"required,oneof=dosen"`
	NewRole       string `json:"new_role" validate:"required,oneof=dosen kajur dekan"`
	Department    string `json:"department,omitempty"`
	Faculty       string `json:"faculty,omitempty"`
	EffectiveDate string `json:"effective_date" validate:"required"`
	Reason        string `json:"reason" validate:"required,min=10"`
}

type RoleAssignmentResponse struct {
	AssignmentID  uint      `json:"assignment_id"`
	UserID        uint      `json:"user_id"`
	UserName      string    `json:"user_name"`
	UserNIDN      string    `json:"user_nidn"`
	OldRole       string    `json:"old_role"`
	NewRole       string    `json:"new_role"`
	Department    string    `json:"department,omitempty"`
	Faculty       string    `json:"faculty,omitempty"`
	EffectiveDate time.Time `json:"effective_date"`
	AssignedBy    string    `json:"assigned_by"`
	Reason        string    `json:"reason"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type PolicyApprovalRequest struct {
	PolicyID     uint   `json:"policy_id" validate:"required"`
	Action       string `json:"action" validate:"required,oneof=approve reject"`
	Comments     string `json:"comments,omitempty"`
	ApprovalDate string `json:"approval_date" validate:"required"`
}

type PolicyApprovalResponse struct {
	PolicyID     uint      `json:"policy_id"`
	PolicyTitle  string    `json:"policy_title"`
	SubmittedBy  string    `json:"submitted_by"`
	Department   string    `json:"department"`
	Faculty      string    `json:"faculty"`
	Action       string    `json:"action"`
	ApprovedBy   string    `json:"approved_by"`
	Comments     string    `json:"comments"`
	ApprovalDate time.Time `json:"approval_date"`
	Status       string    `json:"status"`
}
