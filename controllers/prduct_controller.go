package controllers

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/evermos-service-go/config"
	"github.com/habbazettt/evermos-service-go/middleware"
	"github.com/habbazettt/evermos-service-go/models"
	"github.com/habbazettt/evermos-service-go/services"
	"github.com/habbazettt/evermos-service-go/utils"
	"gorm.io/gorm"
)

func GetAllProducts(c *fiber.Ctx) error {
	db := config.DB
	var products []models.Produk

	// Ambil query params
	namaProduk := c.Query("nama_produk")
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	categoryID, _ := strconv.Atoi(c.Query("category_id"))
	tokoID, _ := strconv.Atoi(c.Query("toko_id"))
	maxHarga, _ := strconv.Atoi(c.Query("max_harga"))
	minHarga, _ := strconv.Atoi(c.Query("min_harga"))

	offset := (page - 1) * limit

	query := db.Model(&models.Produk{}).Preload("FotoProduk")

	if namaProduk != "" {
		query = query.Where("nama_produk LIKE ?", "%"+namaProduk+"%")
	}
	if categoryID > 0 {
		query = query.Where("id_category = ?", categoryID)
	}
	if tokoID > 0 {
		query = query.Where("id_toko = ?", tokoID)
	}
	if maxHarga > 0 {
		query = query.Where("harga_konsumen <= ?", maxHarga)
	}
	if minHarga > 0 {
		query = query.Where("harga_konsumen >= ?", minHarga)
	}

	query = query.Limit(limit).Offset(offset)

	if err := query.Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil produk",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Berhasil mengambil produk",
		"products": products,
	})
}

func GetProductByID(c *fiber.Ctx) error {
	// Extract parameter ID dari request
	id := c.Params("id")

	// Validasi ID harus angka
	produkID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID produk tidak valid",
		})
	}

	// Ambil produk berdasarkan ID
	var produk models.Produk
	err = config.DB.Preload("FotoProduk").First(&produk, produkID).Error

	// Jika produk tidak ditemukan
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Produk tidak ditemukan",
		})
	}

	// Jika ada error lain saat query
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data produk",
		})
	}

	// Return data produk
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil data produk",
		"data":    produk,
	})
}

// CreateProduct - Tambah produk baru
func CreateProduct(c *fiber.Ctx) error {
	// Ambil user_id dari middleware
	userID, err := middleware.ExtractUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Unauthorized",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Cek apakah user memiliki toko
	var toko models.Toko
	if err := config.DB.Where("id_user = ?", userID).First(&toko).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Hanya pemilik toko yang bisa menambah produk",
		})
	}

	// Ambil data dari request form
	namaProduk := c.FormValue("nama_produk")
	deskripsi := c.FormValue("deskripsi")
	idCategory, err := strconv.Atoi(c.FormValue("id_category"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Kategori tidak valid"})
	}
	hargaReseller, err := strconv.Atoi(c.FormValue("harga_reseller"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Harga reseller tidak valid"})
	}
	hargaKonsumen, err := strconv.Atoi(c.FormValue("harga_konsumen"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Harga konsumen tidak valid"})
	}
	stok, err := strconv.Atoi(c.FormValue("stok"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Stok tidak valid"})
	}

	// Buat slug dari nama produk
	slug := utils.GenerateSlug(namaProduk)

	// Gunakan transaksi database untuk memastikan konsistensi data
	tx := config.DB.Begin()

	// Simpan produk ke database
	product := models.Produk{
		NamaProduk:    namaProduk,
		Slug:          slug,
		HargaReseller: hargaReseller,
		HargaKonsumen: hargaKonsumen,
		Stok:          stok,
		Deskripsi:     deskripsi,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		IDToko:        toko.ID,
		IDCategory:    uint(idCategory),
	}

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menyimpan produk"})
	}

	// Upload foto produk ke Cloudinary
	form, err := c.MultipartForm()
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Gagal mengambil file"})
	}

	files := form.File["photos"]
	if len(files) == 0 {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Minimal satu foto produk harus diunggah"})
	}

	var fotoProdukList []models.FotoProduk

	for _, file := range files {
		// Buka file untuk diupload
		src, err := file.Open()
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Gagal membuka file",
				"error":   err.Error(),
			})
		}
		defer src.Close()

		// Upload ke Cloudinary
		url, err := services.UploadToCloudinary(src)
		if err != nil {
			tx.Rollback()
			fmt.Println("Error uploading to Cloudinary:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengunggah foto produk"})
		}

		// Simpan foto produk ke database
		fotoProduk := models.FotoProduk{
			IDProduk:  product.ID,
			URL:       url,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := tx.Create(&fotoProduk).Error; err != nil {
			tx.Rollback()
			fmt.Println("Error saving photo to database:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menyimpan foto produk"})
		}

		fotoProdukList = append(fotoProdukList, fotoProduk)
	}

	// Commit transaksi jika semua berhasil
	tx.Commit()

	// Tambahkan foto produk ke response
	product.FotoProduk = fotoProdukList

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Produk berhasil ditambahkan",
		"produk":  product,
	})
}

// UpdateProduct memperbarui data produk berdasarkan ID
func UpdateProduct(c *fiber.Ctx) error {
	// Ambil user_id dari middleware
	userID, err := middleware.ExtractUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"errors":  err.Error(),
		})
	}

	// Ambil ID produk dari parameter URL
	id := c.Params("id")
	produkID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID produk tidak valid",
		})
	}

	// Cari produk berdasarkan ID
	var produk models.Produk
	if err := config.DB.Preload("FotoProduk").First(&produk, produkID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Produk tidak ditemukan",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data produk",
		})
	}

	// Cek apakah user adalah pemilik toko dari produk tersebut
	var toko models.Toko
	if err := config.DB.Where("id = ?", produk.IDToko).First(&toko).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data toko",
		})
	}
	if toko.IDUser != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Anda tidak berhak mengubah produk ini",
		})
	}

	// Gunakan transaksi database
	tx := config.DB.Begin()

	// Parse form-data request
	form, err := c.MultipartForm()
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Gagal membaca data",
		})
	}

	// Update hanya field yang dikirim dalam request
	if values, ok := form.Value["nama_produk"]; ok && len(values) > 0 {
		produk.NamaProduk = values[0]
		produk.Slug = utils.GenerateSlug(values[0]) // Update slug
	}
	if values, ok := form.Value["deskripsi"]; ok && len(values) > 0 {
		produk.Deskripsi = values[0]
	}
	if values, ok := form.Value["harga_reseller"]; ok && len(values) > 0 {
		hargaReseller, err := strconv.Atoi(values[0])
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Harga reseller harus berupa angka",
			})
		}
		produk.HargaReseller = hargaReseller
	}
	if values, ok := form.Value["harga_konsumen"]; ok && len(values) > 0 {
		hargaKonsumen, err := strconv.Atoi(values[0])
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Harga konsumen harus berupa angka",
			})
		}
		produk.HargaKonsumen = hargaKonsumen
	}
	if values, ok := form.Value["stok"]; ok && len(values) > 0 {
		stok, err := strconv.Atoi(values[0])
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Stok harus berupa angka",
			})
		}
		produk.Stok = stok
	}

	produk.UpdatedAt = time.Now()

	// Simpan perubahan produk
	if err := tx.Save(&produk).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal memperbarui produk",
		})
	}

	// **Handle Update Foto Produk (Opsional)**
	files := form.File["photos"]
	if len(files) > 0 {
		// Hapus foto lama dari database
		if err := tx.Where("id_produk = ?", produk.ID).Delete(&models.FotoProduk{}).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Gagal menghapus foto lama",
			})
		}

		var fotoProdukList []models.FotoProduk

		for _, file := range files {
			// Buka file untuk diupload
			src, err := file.Open()
			if err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Gagal membuka file",
					"error":   err.Error(),
				})
			}
			defer src.Close()

			// Upload ke Cloudinary
			url, err := services.UploadToCloudinary(src)
			if err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Gagal mengunggah foto produk",
				})
			}

			// Simpan foto produk baru ke database
			fotoProduk := models.FotoProduk{
				IDProduk:  produk.ID,
				URL:       url,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err := tx.Create(&fotoProduk).Error; err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Gagal menyimpan foto produk",
				})
			}

			fotoProdukList = append(fotoProdukList, fotoProduk)
		}

		// Update daftar foto di produk
		produk.FotoProduk = fotoProdukList
	}

	// Commit transaksi jika semua berhasil
	tx.Commit()

	// Return response dengan produk yang diperbarui
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Produk berhasil diperbarui",
		"produk":  produk,
	})
}

func DeleteProduct(c *fiber.Ctx) error {
	// Ambil user_id dari middleware
	userID, err := middleware.ExtractUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"errors":  err.Error(),
		})
	}

	// Ambil ID produk dari parameter URL
	id := c.Params("id")
	produkID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID produk tidak valid",
		})
	}

	// Cari produk berdasarkan ID
	var produk models.Produk
	if err := config.DB.Preload("FotoProduk").First(&produk, produkID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Produk tidak ditemukan",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data produk",
		})
	}

	// Cek apakah user adalah pemilik toko dari produk tersebut
	var toko models.Toko
	if err := config.DB.Where("id = ?", produk.IDToko).First(&toko).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data toko",
		})
	}
	if toko.IDUser != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Anda tidak berhak menghapus produk ini",
		})
	}

	// Gunakan transaksi database
	tx := config.DB.Begin()

	// Hapus semua foto produk dari Cloudinary
	for _, foto := range produk.FotoProduk {
		err := services.DeleteFromCloudinary(foto.URL)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Gagal menghapus foto produk dari Cloudinary",
			})
		}
	}

	// Hapus foto produk dari database
	if err := tx.Where("id_produk = ?", produk.ID).Delete(&models.FotoProduk{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghapus foto produk dari database",
		})
	}

	// Hapus produk dari database
	if err := tx.Delete(&produk).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghapus produk",
		})
	}

	// Commit transaksi jika semua berhasil
	tx.Commit()

	// Return response sukses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Produk berhasil dihapus",
	})
}

