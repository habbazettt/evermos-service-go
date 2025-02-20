package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/evermos-service-go/middleware"
	"github.com/habbazettt/evermos-service-go/models"
	"github.com/habbazettt/evermos-service-go/services"
)

// GetListAddress handler
func GetListAddress(c *fiber.Ctx) error {
	userID, err := middleware.ExtractUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Unauthorized",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Ambil query param 'judul_alamat' jika ada
	judulAlamat := c.Query("judul_alamat")

	// Panggil service dengan query params
	addresses, err := services.GetAddressesByUserID(userID, judulAlamat)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to fetch addresses",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Success",
		"errors":  nil,
		"data":    addresses,
	})
}

// GetAlamatByID handler
func GetAlamatByID(c *fiber.Ctx) error {
	userID, err := middleware.ExtractUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Unauthorized",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Ambil `id` dari path parameter
	alamatID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "ID alamat tidak valid",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Panggil service
	alamat, err := services.GetAlamatByID(userID, uint(alamatID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Gagal mengambil alamat",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Success",
		"errors":  nil,
		"data":    alamat,
	})
}

// CreateAlamat handler
func CreateAlamat(c *fiber.Ctx) error {
	userID, err := middleware.ExtractUserID(c) // Pastikan ini mengambil userID dari token JWT
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Unauthorized",
			"errors":  err.Error(),
		})
	}

	alamat := new(models.Alamat)
	if err := c.BodyParser(alamat); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
	}

	// Set userID ke alamat
	alamat.IDUser = userID

	// Panggil service untuk menyimpan alamat
	err = services.CreateAlamat(alamat)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Gagal menambahkan alamat",
			"errors":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Berhasil menambahkan alamat",
		"data":    alamat,
	})
}

// UpdateAlamatByID handler
// UpdateAlamatByID handler
func UpdateAlamatByID(c *fiber.Ctx) error {
	// Ambil alamat ID dari parameter URL
	alamatIDStr := c.Params("id")
	alamatID, err := strconv.ParseUint(alamatIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "ID harus berupa angka",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Gunakan middleware ExtractUserID untuk mendapatkan user_id
	userID, err := middleware.ExtractUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Unauthorized",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Parsing request body hanya untuk field yang diizinkan
	var alamatRequest struct {
		NamaPenerima string `json:"nama_penerima"`
		NoTelp       string `json:"no_telp"`
		DetailAlamat string `json:"detail_alamat"`
	}

	if err := c.BodyParser(&alamatRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid request body",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Mapping data yang akan diupdate
	updateData := map[string]interface{}{}
	if alamatRequest.NamaPenerima != "" {
		updateData["nama_penerima"] = alamatRequest.NamaPenerima
	}
	if alamatRequest.NoTelp != "" {
		updateData["no_telp"] = alamatRequest.NoTelp
	}
	if alamatRequest.DetailAlamat != "" {
		updateData["detail_alamat"] = alamatRequest.DetailAlamat
	}

	// Cek apakah ada data yang akan diupdate
	if len(updateData) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Data yang diupdate tidak boleh kosong",
			"errors":  "Tidak ada perubahan yang dikirim",
			"data":    nil,
		})
	}

	// Panggil service untuk update alamat
	err = services.UpdateAlamatByID(uint(alamatID), userID, updateData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Gagal memperbarui alamat",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Respon jika berhasil
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Alamat berhasil diperbarui",
		"data":    updateData,
	})
}

// DeleteAlamatByID handler
func DeleteAlamatByID(c *fiber.Ctx) error {
	alamatIDStr := c.Params("id") // Ambil ID dari parameter URL
	alamatID, err := strconv.ParseUint(alamatIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "ID harus berupa angka",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	userID := c.Locals("user_id").(uint) // Ambil user_id dari middleware

	err = services.DeleteAlamatByID(uint(alamatID), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Gagal menghapus alamat",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Alamat berhasil dihapus",
		"data":    nil,
	})
}
