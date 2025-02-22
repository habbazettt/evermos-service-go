package models

import "time"

type LogProduk struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	IDProduk      uint      `json:"id_produk"`
	NamaProduk    string    `json:"nama_produk"`
	Slug          string    `json:"slug"`
	HargaReseller int       `json:"harga_reseller"`
	HargaKonsumen int       `json:"harga_konsumen"`
	Deskripsi     string    `json:"deskripsi"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	IDToko        uint      `json:"id_toko"`
	IDCategory    uint      `json:"id_category"`
}
