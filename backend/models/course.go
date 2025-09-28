package models

import "time"

type Course struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	Code       string      `gorm:"unique;not null" json:"code" validate:"required,min=3,max=10"`
	Name       string      `gorm:"type:varchar(100);not null" json:"name" validate:"required,min=3,max=100"`
	Credits    int         `gorm:"not null" json:"credits" validate:"required,min=1,max=6"`
	Semester   int         `gorm:"not null;default:1" json:"semester" validate:"required,min=1,max=14"`
	Prasyarat  string      `gorm:"type:text" json:"prasyarat"`
	Deskripsi  string      `gorm:"type:text" json:"deskripsi"`
	DosenID    *uint       `gorm:"default:null" json:"dosen_id,omitempty"`
	Dosen      *Dosen      `gorm:"foreignKey:DosenID" json:"dosen,omitempty"`
	Mahasiswas []Mahasiswa `gorm:"many2many:mahasiswa_courses;" json:"mahasiswas,omitempty"`
	Jadwal     []Jadwal    `gorm:"foreignKey:CourseID" json:"jadwal,omitempty"`
	KRS        []KRS       `gorm:"foreignKey:CourseID" json:"krs,omitempty"`
	Nilai      []Nilai     `gorm:"foreignKey:CourseID" json:"nilai,omitempty"`
	Absensi    []Absensi   `gorm:"foreignKey:CourseID" json:"absensi,omitempty"`
	Materi     []Materi    `gorm:"foreignKey:CourseID" json:"materi,omitempty"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

// CourseRequest untuk input validation
type CourseRequest struct {
	Code      string `json:"code" validate:"required,min=3,max=10"`
	Name      string `json:"name" validate:"required,min=3,max=100"`
	Credits   int    `json:"credits" validate:"required,min=1,max=6"`
	Semester  int    `json:"semester" validate:"required,min=1,max=14"`
	Prasyarat string `json:"prasyarat"`
	Deskripsi string `json:"deskripsi"`
}
