package models

import "time"

type Kajur struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	NIDN      string    `gorm:"unique;not null" json:"nidn" validate:"required,min=8,max=20"`
	Nama      string    `gorm:"type:varchar(100);not null" json:"nama" validate:"required,min=2,max=100"`
	Email     string    `gorm:"type:varchar(100);unique;not null" json:"email" validate:"required,email"`
	Jurusan   string    `gorm:"type:varchar(100);not null" json:"jurusan" validate:"required,min=2,max=100"`
	Status    string    `gorm:"type:varchar(20);default:'aktif'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type KajurRequest struct {
	NIDN    string `json:"nidn" validate:"required,min=8,max=20"`
	Nama    string `json:"nama" validate:"required,min=2,max=100"`
	Email   string `json:"email" validate:"required,email"`
	Jurusan string `json:"jurusan" validate:"required,min=2,max=100"`
}

type KajurResponse struct {
	ID        uint      `json:"id"`
	NIDN      string    `json:"nidn"`
	Nama      string    `json:"nama"`
	Email     string    `json:"email"`
	Jurusan   string    `json:"jurusan"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Response untuk dashboard kajur
type KajurDashboardResponse struct {
	TotalMahasiswa     int     `json:"total_mahasiswa"`
	MahasiswaAktif     int     `json:"mahasiswa_aktif"`
	MahasiswaCuti      int     `json:"mahasiswa_cuti"`
	MahasiswaDropOut   int     `json:"mahasiswa_drop_out"`
	TotalDosen         int     `json:"total_dosen"`
	DosenAktif         int     `json:"dosen_aktif"`
	TotalMataKuliah    int     `json:"total_mata_kuliah"`
	MataKuliahAktif    int     `json:"mata_kuliah_aktif"`
	IPKRataRata        float64 `json:"ipk_rata_rata"`
	TingkatKelulusan   float64 `json:"tingkat_kelulusan"`
	PendingKRSApproval int     `json:"pending_krs_approval"`
}

// Response untuk laporan jurusan
type LaporanJurusanResponse struct {
	Jurusan            string                  `json:"jurusan"`
	PeriodeLaporan     string                  `json:"periode_laporan"`
	StatistikMahasiswa StatistikMahasiswaKajur `json:"statistik_mahasiswa"`
	StatistikDosen     StatistikDosenKajur     `json:"statistik_dosen"`
	StatistikAkademik  StatistikAkademikKajur  `json:"statistik_akademik"`
	PerformanceKelas   []PerformanceKelasKajur `json:"performance_kelas"`
	TrendSemester      []TrendSemesterKajur    `json:"trend_semester"`
	CreatedAt          time.Time               `json:"created_at"`
}

type StatistikMahasiswaKajur struct {
	Total        int     `json:"total"`
	Aktif        int     `json:"aktif"`
	Cuti         int     `json:"cuti"`
	DropOut      int     `json:"drop_out"`
	Lulus        int     `json:"lulus"`
	IPKRataRata  float64 `json:"ipk_rata_rata"`
	IPKTertinggi float64 `json:"ipk_tertinggi"`
	IPKTerendah  float64 `json:"ipk_terendah"`
}

type StatistikDosenKajur struct {
	Total              int     `json:"total"`
	Aktif              int     `json:"aktif"`
	NonAktif           int     `json:"non_aktif"`
	RataMatkulPerDosen float64 `json:"rata_matkul_per_dosen"`
}

type StatistikAkademikKajur struct {
	TotalMataKuliah    int     `json:"total_mata_kuliah"`
	MataKuliahAktif    int     `json:"mata_kuliah_aktif"`
	TotalKelas         int     `json:"total_kelas"`
	RataKehadiranDosen float64 `json:"rata_kehadiran_dosen"`
	RataKehadiranMhs   float64 `json:"rata_kehadiran_mahasiswa"`
	TingkatKelulusan   float64 `json:"tingkat_kelulusan"`
}

type PerformanceKelasKajur struct {
	CourseID      uint    `json:"course_id"`
	CourseCode    string  `json:"course_code"`
	CourseName    string  `json:"course_name"`
	DosenPengampu string  `json:"dosen_pengampu"`
	JumlahMhs     int     `json:"jumlah_mahasiswa"`
	RataNilai     float64 `json:"rata_nilai"`
	TingkatLulus  float64 `json:"tingkat_lulus"`
	RataKehadiran float64 `json:"rata_kehadiran"`
}

type TrendSemesterKajur struct {
	Semester     string  `json:"semester"`
	TahunAjaran  string  `json:"tahun_ajaran"`
	JumlahMhs    int     `json:"jumlah_mahasiswa"`
	IPKRataRata  float64 `json:"ipk_rata_rata"`
	TingkatLulus float64 `json:"tingkat_lulus"`
}
