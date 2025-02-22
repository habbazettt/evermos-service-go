package models

import (
	"time"
)

type DetailTransaction struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	IDTrx       uint      `json:"id_trx"`
	IDLogProduk uint      `json:"id_log_produk"`
	IDToko      uint      `json:"id_toko"`
	Kuantitas   int       `json:"kuantitas"`
	HargaTotal  int       `json:"harga_total"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	LogProduct  LogProduk `json:"log_product,omitempty" gorm:"foreignKey:IDLogProduk"`
}
