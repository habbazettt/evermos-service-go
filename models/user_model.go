package models

import (
	"time"
)

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Nama         string    `json:"nama"`
	KataSandi    string    `json:"kata_sandi"`
	NoTelp       string    `json:"no_telp" gorm:"unique"`
	TanggalLahir string    `json:"tanggal_lahir"`
	JenisKelamin string    `json:"jenis_kelamin"`
	Tentang      string    `json:"tentang"`
	Pekerjaan    string    `json:"pekerjaan"`
	Email        string    `json:"email"`
	IDProvinsi   string    `json:"id_provinsi"`
	IDKota       string    `json:"id_kota"`
	IsAdmin      bool      `json:"is_admin" gorm:"default:false"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Toko         *Toko     `json:"toko,omitempty" gorm:"foreignKey:IDUser"`
	Alamat       *[]Alamat `json:"alamat,omitempty" gorm:"foreignKey:IDUser"`
}
