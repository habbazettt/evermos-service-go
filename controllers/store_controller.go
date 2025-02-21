package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/evermos-service-go/config"
	"github.com/habbazettt/evermos-service-go/middleware"
	"github.com/habbazettt/evermos-service-go/models"
	"github.com/habbazettt/evermos-service-go/services"
)

// GetMyStore - Mendapatkan toko milik user yang sedang login
func GetMyStore(c *fiber.Ctx) error {
	// Ambil User ID dari Middleware JWT
	userID, err := middleware.ExtractUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "User ID tidak ditemukan dalam token",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Cek apakah toko dengan IDUser ini ada
	var toko models.Toko
	if err := config.DB.Where("id_user = ?", userID).Preload("Produk").First(&toko).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Toko tidak ditemukan",
			"errors":  nil,
			"data":    nil,
		})
	}

	// Return data toko jika ditemukan
	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Berhasil mendapatkan toko",
		"errors":  nil,
		"data":    toko,
	})
}

// Get All Stores with Pagination & Search
func GetAllStores(c *fiber.Ctx) error {
	// Parse query params
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("nama", "")

	// Fetch stores
	stores, totalPages, err := services.GetAllStores(page, limit, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Gagal mengambil daftar toko",
			"errors":  err.Error(),
		})
	}

	// Response
	return c.JSON(fiber.Map{
		"status":      true,
		"message":     "Daftar toko berhasil diambil",
		"total_pages": totalPages,
		"data":        stores,
	})
}

// Get Store by ID
func GetStoreByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "ID tidak valid",
		})
	}

	store, err := services.GetStoreByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Toko ditemukan",
		"data":    store,
	})
}

// Create Store
func CreateStore(c *fiber.Ctx) error {
	// Extract user ID from JWT
	userID, err := middleware.ExtractUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	var storeData models.Toko
	if err := c.BodyParser(&storeData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid request body",
		})
	}

	store, err := services.CreateStore(userID, storeData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Gagal membuat toko",
			"errors":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Toko berhasil dibuat",
		"data":    store,
	})
}

func UpdateStore(c *fiber.Ctx) error {
	// Extract user ID dari JWT
	userID, err := middleware.ExtractUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Unauthorized",
			"errors":  err.Error(),
		})
	}

	// Ambil store ID dari parameter URL
	storeID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "ID tidak valid",
			"errors":  err.Error(),
		})
	}

	// Ambil data toko dari form-data
	namaToko := c.FormValue("nama_toko")

	// Ambil file dari form-data
	file, _ := c.FormFile("photo")
	var photoURL string

	if file != nil {
		// Buka file langsung dari memori (seperti multer di Node.js)
		src, err := file.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  false,
				"message": "Gagal membuka file",
				"errors":  err.Error(),
			})
		}
		defer src.Close()

		// Upload ke Cloudinary langsung dari memory
		uploadResult, err := services.UploadToCloudinary(src)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  false,
				"message": "Gagal mengunggah ke Cloudinary",
				"errors":  err.Error(),
			})
		}
		photoURL = uploadResult
	}

	// Siapkan data update
	updateData := models.Toko{
		NamaToko: namaToko,
	}
	if photoURL != "" {
		updateData.URLFoto = photoURL
	}

	// Panggil service untuk update store
	store, err := services.UpdateStore(userID, uint(storeID), updateData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Gagal memperbarui toko",
			"errors":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Toko berhasil diperbarui",
		"data":    store,
	})
}

// Delete Store
func DeleteStore(c *fiber.Ctx) error {
	// Extract user ID from JWT
	userID, err := middleware.ExtractUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	storeID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "ID tidak valid",
		})
	}

	err = services.DeleteStore(userID, uint(storeID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Toko berhasil dihapus",
	})
}
