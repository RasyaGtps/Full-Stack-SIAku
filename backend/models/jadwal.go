package models

import "time"

type Jadwal struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CourseID    uint      `gorm:"not null" json:"course_id"`
	Hari        string    `gorm:"type:varchar(20);not null" json:"hari"`
	JamMulai    string    `gorm:"type:varchar(10);not null" json:"jam_mulai"`
	JamSelesai  string    `gorm:"type:varchar(10);not null" json:"jam_selesai"`
	Ruangan     string    `gorm:"type:varchar(50)" json:"ruangan"`
	Dosen       string    `gorm:"type:varchar(100)" json:"dosen"`
	TipeKelas   string    `gorm:"type:varchar(20);default:'kuliah'" json:"tipe_kelas"`
	Semester    int       `gorm:"not null" json:"semester"`
	TahunAjaran string    `gorm:"type:varchar(20);not null" json:"tahun_ajaran"`
	Course      Course    `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type JadwalRequest struct {
	CourseID    uint   `json:"course_id" validate:"required"`
	Hari        string `json:"hari" validate:"required,oneof=senin selasa rabu kamis jumat sabtu"`
	JamMulai    string `json:"jam_mulai" validate:"required"`
	JamSelesai  string `json:"jam_selesai" validate:"required"`
	Ruangan     string `json:"ruangan" validate:"required"`
	Dosen       string `json:"dosen" validate:"required"`
	TipeKelas   string `json:"tipe_kelas" validate:"required,oneof=kuliah ujian praktikum"`
	Semester    int    `json:"semester" validate:"required,min=1,max=14"`
	TahunAjaran string `json:"tahun_ajaran" validate:"required"`
}

type JadwalResponse struct {
	ID          uint      `json:"id"`
	CourseCode  string    `json:"course_code"`
	CourseName  string    `json:"course_name"`
	Credits     int       `json:"credits"`
	Hari        string    `json:"hari"`
	JamMulai    string    `json:"jam_mulai"`
	JamSelesai  string    `json:"jam_selesai"`
	Ruangan     string    `json:"ruangan"`
	Dosen       string    `json:"dosen"`
	TipeKelas   string    `json:"tipe_kelas"`
	Semester    int       `json:"semester"`
	TahunAjaran string    `json:"tahun_ajaran"`
	CreatedAt   time.Time `json:"created_at"`
}
