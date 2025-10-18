package models

import "time"

// Users - Tabel utama untuk semua akun yang bisa login
type Users struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username" validate:"required,min=3,max=50"`
	Email     string    `gorm:"unique;not null" json:"email" validate:"required,email"`
	Password  string    `gorm:"not null" json:"-" validate:"required,min=6"`
	Role      string    `gorm:"type:varchar(20);not null" json:"role" validate:"required,oneof=mahasiswa dosen kajur rektor"`
	Status    string    `gorm:"type:varchar(20);default:'aktif'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Request & Response structures for new architecture
type UserRegistrationRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Nama     string `json:"nama" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=mahasiswa dosen kajur rektor"`

	// Detail fields based on role
	NIM            string `json:"nim,omitempty" validate:"omitempty,min=8,max=20"`      // For mahasiswa
	NIDN           string `json:"nidn,omitempty" validate:"omitempty,min=8,max=20"`     // For dosen, kajur, rektor
	PhoneNumber    string `json:"phone_number,omitempty" validate:"omitempty,min=10,max=15"` // For all roles
	Jurusan        string `json:"jurusan,omitempty" validate:"omitempty,min=2,max=100"` // For mahasiswa, dosen, kajur
	Semester       int    `json:"semester,omitempty"`                                   // For mahasiswa
	StatusAkademik string `json:"status_akademik,omitempty"`                            // For mahasiswa
}

type UserLoginRequest struct {
	Identifier string `json:"identifier" validate:"required"` // Could be username, email, nim, or nidn
	Password   string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Role-specific data (simplified)
	RoleData interface{} `json:"role_data,omitempty"`
}
