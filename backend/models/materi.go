package models

import "time"

type Materi struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CourseID   uint      `gorm:"not null" json:"course_id"`
	Judul      string    `gorm:"type:varchar(200);not null" json:"judul" validate:"required,min=3,max=200"`
	Deskripsi  string    `gorm:"type:text" json:"deskripsi"`
	Pertemuan  int       `gorm:"not null" json:"pertemuan" validate:"required,min=1,max=16"`
	TipeMateri string    `gorm:"type:varchar(50);not null" json:"tipe_materi" validate:"required,oneof=slide video document link"`
	FilePath   string    `gorm:"type:varchar(500)" json:"file_path"`
	FileSize   int64     `gorm:"default:0" json:"file_size"`
	URL        string    `gorm:"type:varchar(500)" json:"url"`
	Status     string    `gorm:"type:varchar(20);default:'aktif'" json:"status"`
	Course     Course    `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type MateriRequest struct {
	CourseID   uint   `json:"course_id" validate:"required"`
	Judul      string `json:"judul" validate:"required,min=3,max=200"`
	Deskripsi  string `json:"deskripsi"`
	Pertemuan  int    `json:"pertemuan" validate:"required,min=1,max=16"`
	TipeMateri string `json:"tipe_materi" validate:"required,oneof=slide video document link"`
	URL        string `json:"url"`
}

type MateriResponse struct {
	ID         uint      `json:"id"`
	CourseID   uint      `json:"course_id"`
	CourseName string    `json:"course_name"`
	CourseCode string    `json:"course_code"`
	Judul      string    `json:"judul"`
	Deskripsi  string    `json:"deskripsi"`
	Pertemuan  int       `json:"pertemuan"`
	TipeMateri string    `json:"tipe_materi"`
	FilePath   string    `json:"file_path"`
	FileSize   int64     `json:"file_size"`
	URL        string    `json:"url"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
