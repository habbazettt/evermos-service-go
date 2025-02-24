package services

import (
	"errors"

	"github.com/habbazettt/evermos-service-go/config"
	"github.com/habbazettt/evermos-service-go/models"
)

// GetAddressesByUserID mengambil daftar alamat berdasarkan user ID
func GetAddressesByUserID(userID uint, judulAlamat string) ([]models.Alamat, error) {
	var addresses []models.Alamat
	query := config.DB.Where("id_user = ?", userID)

	// Jika ada filter 'judul_alamat', tambahkan ke query
	if judulAlamat != "" {
		query = query.Where("judul_alamat LIKE ?", "%"+judulAlamat+"%")
	}

	// Jalankan query
	result := query.Find(&addresses)
	if result.Error != nil {
		return nil, errors.New("gagal mengambil daftar alamat")
	}

	return addresses, nil
}

// GetUserAlamat mengambil alamat berdasarkan user ID dan filter judul_alamat
func GetUserAlamat(userID uint, judulAlamat string) ([]models.Alamat, error) {
	var alamat []models.Alamat
	query := config.DB.Where("id_user = ?", userID)
	if judulAlamat != "" {
		query = query.Where("judul_alamat LIKE ?", "%"+judulAlamat+"%")
	}
	if err := query.Find(&alamat).Error; err != nil {
		return nil, err
	}
	return alamat, nil
}

// GetAlamatByID mengambil alamat berdasarkan ID
func GetAlamatByID(userID, alamatID uint) (*models.Alamat, error) {
	var alamat models.Alamat

	if err := config.DB.Where("id_user = ? AND id = ?", userID, alamatID).First(&alamat).Error; err != nil {

		return nil, errors.New("alamat tidak ditemukan")
	}

	return &alamat, nil
}

// CreateAlamat menambahkan alamat baru untuk user
func CreateAlamat(alamat *models.Alamat) error {
	var user models.User
	// Cek apakah user dengan ID tersebut ada
	if err := config.DB.First(&user, alamat.IDUser).Error; err != nil {
		return errors.New("user tidak ditemukan")
	}

	// Simpan alamat ke database
	result := config.DB.Create(&alamat)
	if result.Error != nil {
		return errors.New("gagal menambahkan alamat")
	}
	return nil
}

// UpdateAlamatByID memperbarui alamat berdasarkan ID
func UpdateAlamatByID(alamatID uint, userID uint, updateData map[string]interface{}) error {
	// Debugging: Cek apakah alamat ada sebelum update
	var alamat models.Alamat
	if err := config.DB.Where("id = ? AND id_user = ?", alamatID, userID).First(&alamat).Error; err != nil {
		return errors.New("alamat tidak ditemukan")
	}

	result := config.DB.Model(&alamat).Updates(updateData)
	if result.Error != nil {
		return errors.New("gagal memperbarui alamat")
	}

	if result.RowsAffected == 0 {
		return errors.New("alamat tidak ditemukan atau tidak ada perubahan")
	}

	return nil
}

// DeleteAlamatByID menghapus alamat berdasarkan ID
func DeleteAlamatByID(userID, alamatID uint) error {
	result := config.DB.Where("id_user = ? AND id = ?", userID, alamatID).Delete(&models.Alamat{})
	if result.Error != nil {
		return errors.New("gagal menghapus alamat")
	}
	if result.RowsAffected == 0 {
		return errors.New("alamat tidak ditemukan")
	}
	return nil
}
