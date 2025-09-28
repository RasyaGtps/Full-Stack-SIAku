package models

import "time"

type Dosen struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	NIDN      string      `gorm:"unique;not null" json:"nidn" validate:"required,min=8,max=20"`
	Nama      string      `gorm:"type:varchar(100);not null" json:"nama" validate:"required,min=2,max=100"`
	Email     string      `gorm:"type:varchar(100);unique;not null" json:"email" validate:"required,email"`
	Jurusan   string      `gorm:"type:varchar(100)" json:"jurusan" validate:"required,min=2,max=100"`
	Status    string      `gorm:"type:varchar(20);default:'aktif'" json:"status"`
	Courses   []Course    `gorm:"foreignKey:DosenID" json:"courses,omitempty"`
	Mahasiswa []Mahasiswa `gorm:"foreignKey:DosenWaliID" json:"mahasiswa_wali,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type DosenRequest struct {
	NIDN    string `json:"nidn" validate:"required,min=8,max=20"`
	Nama    string `json:"nama" validate:"required,min=2,max=100"`
	Email   string `json:"email" validate:"required,email"`
	Jurusan string `json:"jurusan" validate:"required,min=2,max=100"`
}

type DosenResponse struct {
	ID        uint      `json:"id"`
	NIDN      string    `json:"nidn"`
	Nama      string    `json:"nama"`
	Email     string    `json:"email"`
	Jurusan   string    `json:"jurusan"`
	Status    string    `json:"status"`
	Courses   []Course  `json:"courses,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
