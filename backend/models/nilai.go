package models

import "time"

type Nilai struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	MahasiswaID uint      `gorm:"not null" json:"mahasiswa_id"`
	CourseID    uint      `gorm:"not null" json:"course_id"`
	Semester    int       `gorm:"not null" json:"semester"`
	TahunAjaran string    `gorm:"type:varchar(20);not null" json:"tahun_ajaran"`
	NilaiTugas  float64   `gorm:"type:decimal(5,2);default:0" json:"nilai_tugas"`
	NilaiUTS    float64   `gorm:"type:decimal(5,2);default:0" json:"nilai_uts"`
	NilaiUAS    float64   `gorm:"type:decimal(5,2);default:0" json:"nilai_uas"`
	NilaiAkhir  float64   `gorm:"type:decimal(5,2);default:0" json:"nilai_akhir"`
	GradeHuruf  string    `gorm:"type:varchar(2)" json:"grade_huruf"`
	GradePoint  float64   `gorm:"type:decimal(3,2);default:0" json:"grade_point"`
	Status      string    `gorm:"type:varchar(20);default:'belum_dinilai'" json:"status"`
	Mahasiswa   Mahasiswa `gorm:"foreignKey:MahasiswaID" json:"mahasiswa,omitempty"`
	Course      Course    `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type NilaiResponse struct {
	ID          uint      `json:"id"`
	CourseCode  string    `json:"course_code"`
	CourseName  string    `json:"course_name"`
	Credits     int       `json:"credits"`
	Semester    int       `json:"semester"`
	TahunAjaran string    `json:"tahun_ajaran"`
	NilaiTugas  float64   `json:"nilai_tugas"`
	NilaiUTS    float64   `json:"nilai_uts"`
	NilaiUAS    float64   `json:"nilai_uas"`
	NilaiAkhir  float64   `json:"nilai_akhir"`
	GradeHuruf  string    `json:"grade_huruf"`
	GradePoint  float64   `json:"grade_point"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type TranskripResponse struct {
	Mahasiswa       MahasiswaResponse `json:"mahasiswa"`
	TotalSKS        int               `json:"total_sks"`
	TotalSKSLulus   int               `json:"total_sks_lulus"`
	IPKKumulatif    float64           `json:"ipk_kumulatif"`
	RiwayatNilai    []NilaiResponse   `json:"riwayat_nilai"`
	StatusKelulusan string            `json:"status_kelulusan"`
}
