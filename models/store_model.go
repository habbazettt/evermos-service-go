package models

import (
	"time"
)

type Toko struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	IDUser    uint      `json:"id_user"`
	NamaToko  string    `json:"nama_toko"`
	URLFoto   string    `json:"url_foto"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Produk    []Produk  `json:"produk,omitempty" gorm:"foreignKey:IDToko"`
}