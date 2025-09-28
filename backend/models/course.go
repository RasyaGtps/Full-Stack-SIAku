package models

import "time"

type Course struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	Code       string      `gorm:"unique;not null" json:"code" validate:"required,min=3,max=10"`
	Name       string      `gorm:"type:varchar(100);not null" json:"name" validate:"required,min=3,max=100"`
	Credits    int         `gorm:"not null" json:"credits" validate:"required,min=1,max=6"`
	Mahasiswas []Mahasiswa `gorm:"many2many:mahasiswa_courses;" json:"mahasiswas,omitempty"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

// CourseRequest untuk input validation
type CourseRequest struct {
	Code    string `json:"code" validate:"required,min=3,max=10"`
	Name    string `json:"name" validate:"required,min=3,max=100"`
	Credits int    `json:"credits" validate:"required,min=1,max=6"`
}
