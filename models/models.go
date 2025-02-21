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

type Toko struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	IDUser    uint      `json:"id_user"`
	NamaToko  string    `json:"nama_toko"`
	URLFoto   string    `json:"url_foto"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Produk    []Produk  `json:"produk,omitempty" gorm:"foreignKey:IDToko"`
}

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

type FotoProduk struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	IDProduk  uint      `json:"id_produk"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Category struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	NamaCategory string    `json:"nama_category"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

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

type DetailTransaction struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	IDTrx       uint      `json:"id_trx"`
	IDLogProduk uint      `json:"id_log_produk"`
	IDToko      uint      `json:"id_toko"`
	Kuantitas   int       `json:"kuantitas"`
	HargaTotal  int       `json:"harga_total"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

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
