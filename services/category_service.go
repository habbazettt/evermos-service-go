package services

import (
	"errors"

	"github.com/habbazettt/evermos-service-go/config"
	"github.com/habbazettt/evermos-service-go/models"
)

// Get All Categories
func GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	err := config.DB.Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// Get Category By ID
func GetCategoryByID(id uint) (*models.Category, error) {
	var category models.Category
	err := config.DB.First(&category, id).Error
	if err != nil {
		return nil, errors.New("kategori tidak ditemukan")
	}
	return &category, nil
}

// Create Category (Admin Only)
func CreateCategory(namaCategory string) (*models.Category, error) {
	category := models.Category{NamaCategory: namaCategory}
	err := config.DB.Create(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// Update Category (Admin Only)
func UpdateCategory(id uint, namaCategory string) (*models.Category, error) {
	var category models.Category
	err := config.DB.First(&category, id).Error
	if err != nil {
		return nil, errors.New("kategori tidak ditemukan")
	}

	category.NamaCategory = namaCategory
	err = config.DB.Save(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// Delete Category (Admin Only)
func DeleteCategory(id uint) error {
	var category models.Category
	err := config.DB.First(&category, id).Error
	if err != nil {
		return errors.New("kategori tidak ditemukan")
	}

	err = config.DB.Delete(&category).Error
	if err != nil {
		return err
	}
	return nil
}
