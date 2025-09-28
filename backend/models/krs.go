package models

import "time"

type KRS struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	MahasiswaID     uint       `gorm:"not null" json:"mahasiswa_id"`
	CourseID        uint       `gorm:"not null" json:"course_id"`
	Semester        int        `gorm:"not null" json:"semester"`
	TahunAjaran     string     `gorm:"type:varchar(20);not null" json:"tahun_ajaran"`
	Status          string     `gorm:"type:varchar(20);default:'pending'" json:"status"`
	ApprovalStatus  string     `gorm:"type:varchar(20);default:'pending'" json:"approval_status"`
	ApprovedBy      *uint      `gorm:"default:null" json:"approved_by,omitempty"`
	ApprovedAt      *time.Time `gorm:"default:null" json:"approved_at,omitempty"`
	RejectionReason string     `gorm:"type:text" json:"rejection_reason,omitempty"`
	Mahasiswa       Mahasiswa  `gorm:"foreignKey:MahasiswaID" json:"mahasiswa,omitempty"`
	Course          Course     `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Approver        *Dosen     `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type KRSRequest struct {
	CourseID    uint   `json:"course_id" validate:"required"`
	Semester    int    `json:"semester" validate:"required,min=1,max=14"`
	TahunAjaran string `json:"tahun_ajaran" validate:"required"`
}

type KRSResponse struct {
	ID              uint       `json:"id"`
	CourseID        uint       `json:"course_id"`
	CourseName      string     `json:"course_name"`
	CourseCode      string     `json:"course_code"`
	Credits         int        `json:"credits"`
	Semester        int        `json:"semester"`
	TahunAjaran     string     `json:"tahun_ajaran"`
	Status          string     `json:"status"`
	ApprovalStatus  string     `json:"approval_status"`
	ApprovedBy      *uint      `json:"approved_by,omitempty"`
	ApprovedAt      *time.Time `json:"approved_at,omitempty"`
	RejectionReason string     `json:"rejection_reason,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type KRSApprovalRequest struct {
	Action          string `json:"action" validate:"required,oneof=approve reject"`
	RejectionReason string `json:"rejection_reason"`
}
