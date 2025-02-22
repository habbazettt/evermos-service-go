package models

import "time"

type Alamat struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	IDUser       uint      `json:"id_user"`
	JudulAlamat  string    `json:"judul_alamat"`
	NamaPenerima string    `json:"nama_penerima"`
	NoTelp       string    `json:"no_telp"`
	DetailAlamat string    `json:"detail_alamat"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
