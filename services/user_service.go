package services

import (
	"errors"

	"github.com/habbazettt/evermos-service-go/config"
	"github.com/habbazettt/evermos-service-go/models"
)

type UpdateUserRequest struct {
	Nama         string `json:"nama"`
	KataSandi    string `json:"kata_sandi"`
	NoTelp       string `json:"no_telp"`
	TanggalLahir string `json:"tanggal_Lahir"`
	Pekerjaan    string `json:"pekerjaan"`
	JenisKelamin string `json:"jenis_kelamin"`
	Tentang      string `json:"tentang"`
	Email        string `json:"email"`
	IDProvinsi   string `json:"id_provinsi"`
	IDKota       string `json:"id_kota"`
	IsAdmin      bool   `json:"is_admin"`
}

// GetUserByID mengambil user berdasarkan ID
func GetUserByID(userID string) (*models.User, error) {
	var user models.User

	// Query database untuk mendapatkan user berdasarkan ID
	result := config.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	return &user, nil
}

// UpdateUserByID memperbarui data user berdasarkan ID
func UpdateUserByID(userID string, updateData UpdateUserRequest) (*models.User, error) {
	var user models.User

	// Cek apakah user ada di database
	result := config.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	// Update data user
	user.Nama = updateData.Nama
	user.KataSandi = updateData.KataSandi
	user.NoTelp = updateData.NoTelp
	user.TanggalLahir = updateData.TanggalLahir
	user.Pekerjaan = updateData.Pekerjaan
	user.JenisKelamin = updateData.JenisKelamin
	user.Tentang = updateData.Tentang
	user.Email = updateData.Email
	user.IDProvinsi = updateData.IDProvinsi
	user.IDKota = updateData.IDKota
	user.IsAdmin = updateData.IsAdmin

	// Simpan perubahan
	if err := config.DB.Save(&user).Error; err != nil {
		return nil, errors.New("gagal memperbarui user")
	}

	// Perbarui nama toko jika ada
	if err := config.DB.Model(&models.Toko{}).Where("id_user = ?", user.ID).Update("nama_toko", user.Nama+" Store").Error; err != nil {
		return nil, errors.New("gagal memperbarui nama toko")
	}

	return &user, nil
}
