package services

import (
	"context"
	"errors"
	"math"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/habbazettt/evermos-service-go/config"
	"github.com/habbazettt/evermos-service-go/models"
	"gorm.io/gorm"
)

// Get All Stores with Pagination and Search
func GetAllStores(page, limit int, search string) ([]models.Toko, int, error) {
	var stores []models.Toko
	var total int64

	query := config.DB.Model(&models.Toko{})

	// Search by store name
	if search != "" {
		query = query.Where("nama_toko LIKE ?", "%"+search+"%")
	}

	// Count total stores
	query.Count(&total)

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch data with limit and offset
	err := query.Limit(limit).Offset(offset).Find(&stores).Error
	if err != nil {
		return nil, 0, err
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return stores, totalPages, nil
}

// Get Store By ID
func GetStoreByID(id uint) (*models.Toko, error) {
	var store models.Toko
	err := config.DB.Where("id = ?", id).First(&store).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("toko tidak ditemukan")
	}
	return &store, err
}

// Create Store
func CreateStore(userID uint, storeData models.Toko) (*models.Toko, error) {
	storeData.IDUser = userID
	err := config.DB.Create(&storeData).Error
	if err != nil {
		return nil, err
	}
	return &storeData, nil
}

// UploadToCloudinary mengunggah gambar langsung ke Cloudinary dari memory
func UploadToCloudinary(file multipart.File) (string, error) {
	ctx := context.Background()

	uploadResult, err := config.CLD.Upload.Upload(ctx, file, uploader.UploadParams{})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}

// UpdateStore memperbarui informasi toko
func UpdateStore(userID, storeID uint, updateData models.Toko) (*models.Toko, error) {
	var store models.Toko

	// Cek apakah toko ada dan milik user
	err := config.DB.Where("id = ? AND id_user = ?", storeID, userID).First(&store).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("toko tidak ditemukan atau tidak memiliki akses")
	}

	// Perbarui data toko
	err = config.DB.Model(&store).Updates(updateData).Error
	if err != nil {
		return nil, err
	}

	return &store, nil
}

// Delete Store
func DeleteStore(userID, storeID uint) error {
	var store models.Toko

	// Check if store exists and user has access
	err := config.DB.Where("id = ? AND id_user = ?", storeID, userID).First(&store).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("toko tidak ditemukan atau tidak memiliki akses")
	}

	// Delete the store
	err = config.DB.Delete(&store).Error
	if err != nil {
		return err
	}

	return nil
}
