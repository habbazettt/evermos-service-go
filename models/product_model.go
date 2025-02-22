package models

import "time"

type Produk struct {
	ID            uint         `json:"id" gorm:"primaryKey;autoIncrement"`
	NamaProduk    string       `json:"nama_produk"`
	Slug          string       `json:"slug"`
	HargaReseller int          `json:"harga_reseller"`
	HargaKonsumen int          `json:"harga_konsumen"`
	Stok          int          `json:"stok"`
	Deskripsi     string       `json:"deskripsi"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	IDToko        uint         `json:"id_toko"`
	IDCategory    uint         `json:"id_category"`
	FotoProduk    []FotoProduk `json:"foto_produk,omitempty" gorm:"foreignKey:IDProduk"`
}
