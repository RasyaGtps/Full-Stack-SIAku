package models

import "time"

type Mahasiswa struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	NIM            string    `gorm:"unique;not null" json:"nim" validate:"required,min=8,max=20"`
	Nama           string    `gorm:"type:varchar(100);not null" json:"nama" validate:"required,min=2,max=100"`
	Jurusan        string    `gorm:"type:varchar(100)" json:"jurusan" validate:"required,min=2,max=100"`
	PhoneNumber    string    `gorm:"type:varchar(20)" json:"phone_number,omitempty"`
	StatusAkademik string    `gorm:"type:varchar(20);default:'aktif'" json:"status_akademik"`
	Semester       int       `gorm:"default:1" json:"semester"`
	IPK            float64   `gorm:"type:decimal(3,2);default:0.00" json:"ipk"`
	DosenWaliID    *uint     `gorm:"default:null" json:"dosen_wali_id,omitempty"`
	Courses        []Course  `gorm:"many2many:mahasiswa_courses;" json:"courses,omitempty"`
	KRS            []KRS     `gorm:"foreignKey:MahasiswaID" json:"krs,omitempty"`
	Nilai          []Nilai   `gorm:"foreignKey:MahasiswaID" json:"nilai,omitempty"`
	DosenWali      *Dosen    `gorm:"foreignKey:DosenWaliID" json:"dosen_wali,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// MahasiswaRequest untuk input validation
type MahasiswaRequest struct {
	NIM     string `json:"nim" validate:"required,min=8,max=20"`
	Nama    string `json:"nama" validate:"required,min=2,max=100"`
	Jurusan string `json:"jurusan" validate:"required,min=2,max=100"`
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

// LoginRequest sudah tidak diperlukan, gunakan UserLoginRequest di users.go
