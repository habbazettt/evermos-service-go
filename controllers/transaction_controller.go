package controllers

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/evermos-service-go/config"
	"github.com/habbazettt/evermos-service-go/middleware"
	"github.com/habbazettt/evermos-service-go/models"
)

// CreateTransaction - Membuat transaksi baru
func CreateTransaction(c *fiber.Ctx) error {
	// Ambil user_id dari middleware
	userID, err := middleware.ExtractUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	// Parsing request body
	var req struct {
		MethodBayar      string `json:"method_bayar"`
		AlamatPengiriman uint   `json:"alamat_kirim"`
		DetailTransaksi  []struct {
			ProductID uint `json:"product_id"`
			Kuantitas int  `json:"kuantitas"`
		} `json:"detail_transaksi"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	// Periksa apakah alamat kirim valid
	var alamat models.Alamat
	if err := config.DB.Where("id = ? AND id_user = ?", req.AlamatPengiriman, userID).First(&alamat).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Alamat tidak ditemukan"})
	}

	// Generate kode invoice unik
	kodeInvoice := fmt.Sprintf("INV-%d", time.Now().Unix())

	// Inisialisasi transaksi baru
	transaction := models.Transaction{
		IDUser:           userID,
		AlamatPengiriman: req.AlamatPengiriman,
		KodeInvoice:      kodeInvoice,
		MethodBayar:      req.MethodBayar,
		HargaTotal:       0,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Simpan transaksi ke database
	tx := config.DB.Begin()
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menyimpan transaksi"})
	}

	// Proses detail transaksi
	var totalHarga int
	var detailTrxList []models.DetailTransaction

	for _, detail := range req.DetailTransaksi {
		// Cek apakah produk ada di tabel Produk
		var produk models.Produk
		if err := tx.Preload("FotoProduk").First(&produk, detail.ProductID).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Produk tidak ditemukan"})
		}

		// Pastikan stok mencukupi
		if produk.Stok < detail.Kuantitas {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Stok produk tidak mencukupi"})
		}

		// Cek apakah sudah ada di LogProduk
		var logProduk models.LogProduk
		if err := tx.Where("id_produk = ?", produk.ID).First(&logProduk).Error; err != nil {
			// Jika belum ada, buat entri baru di LogProduk
			logProduk = models.LogProduk{
				IDProduk:      produk.ID,
				NamaProduk:    produk.NamaProduk,
				Slug:          produk.Slug,
				HargaReseller: produk.HargaReseller,
				HargaKonsumen: produk.HargaKonsumen,
				Deskripsi:     produk.Deskripsi,
				IDToko:        produk.IDToko,
				IDCategory:    produk.IDCategory,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}

			if err := tx.Create(&logProduk).Error; err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menyimpan ke LogProduk"})
			}
		}

		// Hitung harga total
		hargaTotal := produk.HargaKonsumen * detail.Kuantitas
		totalHarga += hargaTotal

		// Simpan detail transaksi
		detailTrx := models.DetailTransaction{
			IDTrx:       transaction.ID,
			IDLogProduk: logProduk.ID, // Pakai ID dari LogProduk
			IDToko:      produk.IDToko,
			Kuantitas:   detail.Kuantitas,
			HargaTotal:  hargaTotal,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := tx.Create(&detailTrx).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menyimpan detail transaksi"})
		}

		detailTrxList = append(detailTrxList, detailTrx)

		// Update stok produk
		produk.Stok -= detail.Kuantitas
		if err := tx.Save(&produk).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal memperbarui stok produk"})
		}
	}

	// Update total harga transaksi
	transaction.HargaTotal = totalHarga
	if err := tx.Save(&transaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal memperbarui total harga transaksi"})
	}

	tx.Commit()

	// Ambil transaksi yang baru dibuat
	var trx models.Transaction
	if err := config.DB.Preload("DetailTransaksi.LogProduct").Preload("Alamat").First(&trx, transaction.ID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengambil transaksi yang baru dibuat"})
	}

	// Format detail transaksi dengan informasi lengkap
	formattedDetailTrx := []map[string]interface{}{}
	for _, detail := range trx.DetailTransaksi {
		formattedDetailTrx = append(formattedDetailTrx, map[string]interface{}{
			"id":            detail.ID,
			"id_trx":        detail.IDTrx,
			"id_log_produk": detail.IDLogProduk,
			"id_toko":       detail.IDToko,
			"kuantitas":     detail.Kuantitas,
			"harga_total":   detail.HargaTotal,
			"created_at":    detail.CreatedAt,
			"updated_at":    detail.UpdatedAt,
			"log_product": map[string]interface{}{
				"id":             detail.LogProduct.ID,
				"id_produk":      detail.LogProduct.IDProduk,
				"nama_produk":    detail.LogProduct.NamaProduk,
				"slug":           detail.LogProduct.Slug,
				"harga_reseller": detail.LogProduct.HargaReseller,
				"harga_konsumen": detail.LogProduct.HargaKonsumen,
				"deskripsi":      detail.LogProduct.Deskripsi,
				"created_at":     detail.LogProduct.CreatedAt,
				"updated_at":     detail.LogProduct.UpdatedAt,
				"id_toko":        detail.LogProduct.IDToko,
				"id_category":    detail.LogProduct.IDCategory,
			},
		})
	}

	// **Return response sukses dengan format lengkap**
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Transaksi berhasil dibuat",
		"data": map[string]interface{}{
			"id":           trx.ID,
			"harga_total":  trx.HargaTotal,
			"kode_invoice": trx.KodeInvoice,
			"method_bayar": trx.MethodBayar,
			"alamat_kirim": map[string]interface{}{
				"id":            trx.AlamatPengiriman,
				"judul_alamat":  alamat.JudulAlamat,
				"nama_penerima": alamat.NamaPenerima,
				"no_telp":       alamat.NoTelp,
				"detail_alamat": alamat.DetailAlamat,
			},
			"detail_trx": formattedDetailTrx,
		},
	})
}

// GetAllTransactions - Mendapatkan semua transaksi user
func GetAllTransactions(c *fiber.Ctx) error {
	// Ambil query params
	search := c.Query("search")                      // Filter berdasarkan kode invoice
	limit, _ := strconv.Atoi(c.Query("limit", "10")) // Default 10
	page, _ := strconv.Atoi(c.Query("page", "1"))    // Default 1
	offset := (page - 1) * limit

	// Query transaksi
	var transactions []models.Transaction
	query := config.DB.Preload("DetailTransaksi.LogProduct").Preload("Alamat")

	// Jika ada parameter search, filter berdasarkan kode_invoice
	if search != "" {
		query = query.Where("kode_invoice LIKE ?", "%"+search+"%")
	}

	// Eksekusi query dengan pagination
	if err := query.Limit(limit).Offset(offset).Find(&transactions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengambil transaksi"})
	}

	// Hitung total transaksi untuk pagination
	var total int64
	config.DB.Model(&models.Transaction{}).Where("kode_invoice LIKE ?", "%"+search+"%").Count(&total)

	// Response
	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil transaksi",
		"data":    transactions,
		"pagination": fiber.Map{
			"page":       page,
			"limit":      limit,
			"total_data": total,
			"total_page": int(math.Ceil(float64(total) / float64(limit))),
		},
	})
}

// GetTransactionByID - Mendapatkan transaksi berdasarkan ID
func GetTransactionByID(c *fiber.Ctx) error {
	// Ambil query params
	transactionID := c.Params("id")
	limit, _ := strconv.Atoi(c.Query("limit", "10")) // Default 10
	page, _ := strconv.Atoi(c.Query("page", "1"))    // Default 1
	offset := (page - 1) * limit

	// Query transaksi berdasarkan ID
	var transaction models.Transaction
	query := config.DB.Preload("DetailTransaksi.LogProduct").Preload("Alamat")

	// Cari transaksi dengan parameter ID
	query = query.Where("id = ?", transactionID)

	// Eksekusi query dengan pagination
	if err := query.Limit(limit).Offset(offset).First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Transaksi tidak ditemukan"})
	}

	// Response dengan format yang sesuai
	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil transaksi",
		"data": map[string]interface{}{
			"id":           transaction.ID,
			"harga_total":  transaction.HargaTotal,
			"kode_invoice": transaction.KodeInvoice,
			"method_bayar": transaction.MethodBayar,
			"alamat_kirim": map[string]interface{}{
				"id":            transaction.Alamat.ID,
				"judul_alamat":  transaction.Alamat.JudulAlamat,
				"nama_penerima": transaction.Alamat.NamaPenerima,
				"no_telp":       transaction.Alamat.NoTelp,
				"detail_alamat": transaction.Alamat.DetailAlamat,
			},
			"detail_trx": transaction.DetailTransaksi,
		},
	})
}
