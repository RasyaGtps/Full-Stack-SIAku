package models

import "time"

type Mahasiswa struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	NIM       string    `gorm:"unique;not null" json:"nim"`
	Nama      string    `gorm:"type:varchar(100);not null" json:"nama"`
	Jurusan   string    `gorm:"type:varchar(100)" json:"jurusan"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
