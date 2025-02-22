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

	// Map untuk menyimpan data yang akan diperbarui
	updateFields := map[string]interface{}{}

	// Hanya update jika field diisi
	if updateData.Nama != "" {
		updateFields["nama"] = updateData.Nama
	}
	if updateData.KataSandi != "" {
		updateFields["kata_sandi"] = updateData.KataSandi
	}
	if updateData.NoTelp != "" {
		updateFields["no_telp"] = updateData.NoTelp
	}
	if updateData.TanggalLahir != "" {
		updateFields["tanggal_lahir"] = updateData.TanggalLahir
	}
	if updateData.Pekerjaan != "" {
		updateFields["pekerjaan"] = updateData.Pekerjaan
	}
	if updateData.JenisKelamin != "" {
		updateFields["jenis_kelamin"] = updateData.JenisKelamin
	}
	if updateData.Tentang != "" {
		updateFields["tentang"] = updateData.Tentang
	}
	if updateData.Email != "" {
		updateFields["email"] = updateData.Email
	}
	if updateData.IDProvinsi != "" {
		updateFields["id_provinsi"] = updateData.IDProvinsi
	}
	if updateData.IDKota != "" {
		updateFields["id_kota"] = updateData.IDKota
	}

	// Pastikan hanya admin yang bisa mengubah `is_admin`
	if updateData.IsAdmin != user.IsAdmin {
		updateFields["is_admin"] = updateData.IsAdmin
	}

	// Jika tidak ada perubahan, kembalikan user tanpa update ke database
	if len(updateFields) == 0 {
		return &user, nil
	}

	// Update user di database
	if err := config.DB.Model(&user).Updates(updateFields).Error; err != nil {
		return nil, errors.New("gagal memperbarui user")
	}

	// Jika nama diperbarui, update juga nama toko
	if namaBaru, ok := updateFields["nama"]; ok {
		if err := config.DB.Model(&models.Toko{}).Where("id_user = ?", user.ID).Update("nama_toko", namaBaru.(string)+" Store").Error; err != nil {
			return nil, errors.New("gagal memperbarui nama toko")
		}
	}

	return &user, nil
}
