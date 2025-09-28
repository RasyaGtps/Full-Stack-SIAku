package models

import "time"

type Mahasiswa struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	NIM       string    `gorm:"unique;not null" json:"nim" validate:"required,min=8,max=20"`
	Nama      string    `gorm:"type:varchar(100);not null" json:"nama" validate:"required,min=2,max=100"`
	Jurusan   string    `gorm:"type:varchar(100)" json:"jurusan" validate:"required,min=2,max=100"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-" validate:"required,min=6"`
	Courses   []Course  `gorm:"many2many:mahasiswa_courses;" json:"courses,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MahasiswaRequest untuk input validation
type MahasiswaRequest struct {
	NIM      string `json:"nim" validate:"required,min=8,max=20"`
	Nama     string `json:"nama" validate:"required,min=2,max=100"`
	Jurusan  string `json:"jurusan" validate:"required,min=2,max=100"`
	Password string `json:"password" validate:"required,min=6"`
}

// MahasiswaResponse untuk output (tanpa password)
type MahasiswaResponse struct {
	ID        uint      `json:"id"`
	NIM       string    `json:"nim"`
	Nama      string    `json:"nama"`
	Jurusan   string    `json:"jurusan"`
	Courses   []Course  `json:"courses,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginRequest untuk login
type LoginRequest struct {
	NIM      string `json:"nim" validate:"required"`
	Password string `json:"password" validate:"required"`
}
