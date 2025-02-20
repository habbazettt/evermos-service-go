package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/evermos-service-go/services"
)

// Get All Categories
func GetAllCategories(c *fiber.Ctx) error {
	categories, err := services.GetAllCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Gagal mengambil kategori",
			"errors":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Daftar kategori berhasil diambil",
		"data":    categories,
	})
}

// Get Category By ID
func GetCategoryByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "ID tidak valid",
		})
	}

	category, err := services.GetCategoryByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Kategori ditemukan",
		"data":    category,
	})
}

// Create Category (Admin Only)
func CreateCategory(c *fiber.Ctx) error {
	type Request struct {
		NamaCategory string `json:"nama_category"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid request body",
		})
	}

	category, err := services.CreateCategory(req.NamaCategory)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Gagal membuat kategori",
			"errors":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Kategori berhasil dibuat",
		"data":    category,
	})
}

// Update Category (Admin Only)
func UpdateCategory(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "ID tidak valid",
		})
	}

	type Request struct {
		NamaCategory string `json:"nama_category"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid request body",
		})
	}

	category, err := services.UpdateCategory(uint(id), req.NamaCategory)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Kategori berhasil diperbarui",
		"data":    category,
	})
}

// Delete Category (Admin Only)
func DeleteCategory(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "ID tidak valid",
		})
	}

	err = services.DeleteCategory(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Kategori berhasil dihapus",
	})
}
