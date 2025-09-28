package models

import "time"

type Absensi struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CourseID    uint      `gorm:"not null" json:"course_id"`
	MahasiswaID uint      `gorm:"not null" json:"mahasiswa_id"`
	Pertemuan   int       `gorm:"not null" json:"pertemuan" validate:"required,min=1,max=16"`
	Tanggal     time.Time `gorm:"not null" json:"tanggal"`
	Status      string    `gorm:"type:varchar(10);not null" json:"status" validate:"required,oneof=hadir izin sakit alfa"`
	Keterangan  string    `gorm:"type:text" json:"keterangan"`
	Course      Course    `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Mahasiswa   Mahasiswa `gorm:"foreignKey:MahasiswaID" json:"mahasiswa,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AbsensiRequest struct {
	CourseID    uint   `json:"course_id" validate:"required"`
	MahasiswaID uint   `json:"mahasiswa_id" validate:"required"`
	Pertemuan   int    `json:"pertemuan" validate:"required,min=1,max=16"`
	Tanggal     string `json:"tanggal" validate:"required"`
	Status      string `json:"status" validate:"required,oneof=hadir izin sakit alfa"`
	Keterangan  string `json:"keterangan"`
}

type AbsensiResponse struct {
	ID          uint      `json:"id"`
	CourseID    uint      `json:"course_id"`
	CourseName  string    `json:"course_name"`
	CourseCode  string    `json:"course_code"`
	MahasiswaID uint      `json:"mahasiswa_id"`
	NIM         string    `json:"nim"`
	NamaMhs     string    `json:"nama_mahasiswa"`
	Pertemuan   int       `json:"pertemuan"`
	Tanggal     time.Time `json:"tanggal"`
	Status      string    `json:"status"`
	Keterangan  string    `json:"keterangan"`
	CreatedAt   time.Time `json:"created_at"`
}

type AbsensiPertemuanRequest struct {
	CourseID  uint                    `json:"course_id" validate:"required"`
	Pertemuan int                     `json:"pertemuan" validate:"required,min=1,max=16"`
	Tanggal   string                  `json:"tanggal" validate:"required"`
	Absensi   []AbsensiMahasiswaInput `json:"absensi" validate:"required,dive"`
}

type AbsensiMahasiswaInput struct {
	MahasiswaID uint   `json:"mahasiswa_id" validate:"required"`
	Status      string `json:"status" validate:"required,oneof=hadir izin sakit alfa"`
	Keterangan  string `json:"keterangan"`
}

type RekapAbsensiResponse struct {
	MahasiswaID         uint    `json:"mahasiswa_id"`
	NIM                 string  `json:"nim"`
	NamaMahasiswa       string  `json:"nama_mahasiswa"`
	TotalHadir          int     `json:"total_hadir"`
	TotalIzin           int     `json:"total_izin"`
	TotalSakit          int     `json:"total_sakit"`
	TotalAlfa           int     `json:"total_alfa"`
	TotalPertemuan      int     `json:"total_pertemuan"`
	PersentaseKehadiran float64 `json:"persentase_kehadiran"`
}
