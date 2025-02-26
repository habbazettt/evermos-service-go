package controllers

import (
	"errors"
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

// Get All Products
// @Summary Get All Products
// @Description Get all products with optional filters.
// @Tags Product
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param nama_produk query string false "Filter by product name"
// @Param limit query int false "Limit per page" default(10)
// @Param page query int false "Page number" default(1)
// @Param category_id query int false "Filter by category ID"
// @Param toko_id query int false "Filter by store ID"
// @Param max_harga query int false "Filter by maximum price"
// @Param min_harga query int false "Filter by minimum price"
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router /product [get]
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

// Get Product by ID
// @Summary Get Product by ID
// @Description Get a product by its ID.
// @Tags Product
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /product/{id} [get]
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

// Create Product
// @Summary Create Product
// @Description Create a new product.
// @Tags Product
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param nama_produk formData string true "Product name"
// @Param deskripsi formData string true "Product description"
// @Param id_category formData int true "Category ID"
// @Param harga_reseller formData int true "Reseller price"
// @Param harga_konsumen formData int true "Consumer price"
// @Param stok formData int true "Product stock"
// @Param photos formData file true "Product photos (multiple files allowed)"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Router /product [post]
func CreateProduct(c *fiber.Ctx) error {
	userID, err := middleware.ExtractUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	// Cek apakah user memiliki toko
	var toko models.Toko
	if err := config.DB.Where("id_user = ?", userID).First(&toko).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Hanya pemilik toko yang bisa menambah produk"})
	}

	// Ambil data dari request
	namaProduk := c.FormValue("nama_produk")
	deskripsi := c.FormValue("deskripsi")
	idCategory, _ := strconv.Atoi(c.FormValue("id_category"))
	hargaReseller, _ := strconv.Atoi(c.FormValue("harga_reseller"))
	hargaKonsumen, _ := strconv.Atoi(c.FormValue("harga_konsumen"))
	stok, _ := strconv.Atoi(c.FormValue("stok"))

	slug := utils.GenerateSlug(namaProduk)

	// Mulai transaksi database
	tx := config.DB.Begin()

	product := models.Produk{
		NamaProduk:    namaProduk,
		Slug:          slug,
		HargaReseller: hargaReseller,
		HargaKonsumen: hargaKonsumen,
		Stok:          stok,
		Deskripsi:     deskripsi,
		IDToko:        toko.ID,
		IDCategory:    uint(idCategory),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menyimpan produk"})
	}

	// Tambahkan ke LogProduk
	logProduct := models.LogProduk{
		IDProduk:      product.ID,
		NamaProduk:    product.NamaProduk,
		Slug:          product.Slug,
		HargaReseller: product.HargaReseller,
		HargaKonsumen: product.HargaKonsumen,
		Deskripsi:     product.Deskripsi,
		IDToko:        product.IDToko,
		IDCategory:    product.IDCategory,
	}
	if err := tx.Create(&logProduct).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menyimpan ke log produk"})
	}

	// Upload foto produk ke Cloudinary
	form, err := c.MultipartForm()
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Gagal mengambil file"})
	}

	files := form.File["photos"]
	var fotoProdukList []models.FotoProduk

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal membuka file"})
		}
		defer src.Close()

		url, err := services.UploadToCloudinary(src)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengunggah foto produk"})
		}

		fotoProduk := models.FotoProduk{
			IDProduk:  product.ID,
			URL:       url,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := tx.Create(&fotoProduk).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menyimpan foto produk"})
		}

		fotoProdukList = append(fotoProdukList, fotoProduk)
	}

	tx.Commit()

	product.FotoProduk = fotoProdukList

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Produk berhasil ditambahkan",
		"produk":  product,
	})
}

// Update Product
// @Summary Update Product
// @Description Update a product's information.
// @Tags Product
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param nama_produk formData string false "Product name"
// @Param deskripsi formData string false "Product description"
// @Param photos formData file false "Product photos (multiple files allowed)"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /product/{id} [put]
func UpdateProduct(c *fiber.Ctx) error {
	userID, err := middleware.ExtractUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	id := c.Params("id")
	produkID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID produk tidak valid"})
	}

	var produk models.Produk
	if err := config.DB.Preload("FotoProduk").First(&produk, produkID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Produk tidak ditemukan"})
	}

	var toko models.Toko
	if err := config.DB.Where("id = ?", produk.IDToko).First(&toko).Error; err != nil || toko.IDUser != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Anda tidak memiliki izin untuk mengubah produk ini"})
	}

	tx := config.DB.Begin()

	form, _ := c.MultipartForm()

	if values, ok := form.Value["nama_produk"]; ok && len(values) > 0 {
		produk.NamaProduk = values[0]
		produk.Slug = utils.GenerateSlug(values[0])
	}
	if values, ok := form.Value["deskripsi"]; ok && len(values) > 0 {
		produk.Deskripsi = values[0]
	}
	produk.UpdatedAt = time.Now()

	if err := tx.Save(&produk).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal memperbarui produk"})
	}

	files := form.File["photos"]
	var fotoProdukList []models.FotoProduk

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal membuka file"})
		}
		defer src.Close()

		url, err := services.UploadToCloudinary(src)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengunggah foto produk"})
		}

		fotoProduk := models.FotoProduk{
			IDProduk: produk.ID,
			URL:      url,
		}

		if err := tx.Create(&fotoProduk).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menyimpan foto produk"})
		}

		fotoProdukList = append(fotoProdukList, fotoProduk)
	}

	tx.Commit()

	produk.FotoProduk = fotoProdukList

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Produk berhasil diperbarui",
		"produk":  produk,
	})
}

// Delete Product
// @Summary Delete Product
// @Description Delete a product by its ID.
// @Tags Product
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /product/{id} [delete]
func DeleteProduct(c *fiber.Ctx) error {
	userID, err := middleware.ExtractUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	id := c.Params("id")
	produkID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID produk tidak valid"})
	}

	var produk models.Produk
	if err := config.DB.Preload("FotoProduk").First(&produk, produkID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Produk tidak ditemukan"})
	}

	var toko models.Toko
	if err := config.DB.Where("id = ?", produk.IDToko).First(&toko).Error; err != nil || toko.IDUser != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Anda tidak memiliki izin untuk menghapus produk ini"})
	}

	tx := config.DB.Begin()

	for _, foto := range produk.FotoProduk {
		services.DeleteFromCloudinary(foto.URL)
	}

	if err := tx.Where("id_produk = ?", produk.ID).Delete(&models.FotoProduk{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menghapus foto produk dari database"})
	}

	if err := tx.Delete(&produk).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menghapus produk"})
	}

	tx.Commit()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Produk berhasil dihapus",
		"produk":  produk,
	})
}
