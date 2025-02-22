package models

import "time"

type FotoProduk struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	IDProduk  uint      `json:"id_produk"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
