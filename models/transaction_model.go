package models

import "time"

type Transaction struct {
	ID               uint                `json:"id" gorm:"primaryKey;autoIncrement"`
	IDUser           uint                `json:"id_user"`
	AlamatPengiriman uint                `json:"alamat_pengiriman"`
	HargaTotal       int                 `json:"harga_total"`
	KodeInvoice      string              `json:"kode_invoice"`
	MethodBayar      string              `json:"method_bayar"`
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
	DetailTransaksi  []DetailTransaction `json:"detail_transaksi,omitempty" gorm:"foreignKey:IDTrx"`
}
