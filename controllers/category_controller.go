package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/evermos-service-go/services"
)

// Get All Categories
// @summary Get All Categories
// @description Get a list of all categories.
// @tags Category
// @accept json
// @produce json
// @success 200 {object} Response
// @failure 500 {object} Response
// @router /category [get]
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
// @summary Get Category By ID
// @description Get detailed information of a specific category by ID.
// @tags Category
// @accept json
// @produce json
// @param id path int true "Category ID"
// @success 200 {object} Response
// @failure 400 {object} Response
// @failure 404 {object} Response
// @router /category/{id} [get]
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
// @summary Create a new category
// @description Create a new category (Admin only).
// @tags Category
// @accept json
// @produce json
// @security BearerAuth
// @param request body object{nama_category=string} true "Category Data"
// @success 201 {object} Response
// @failure 400 {object} Response
// @failure 401 {object} Response
// @failure 500 {object} Response
// @router /category [post]
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
// @summary Update a category
// @description Update a category's name (Admin only).
// @tags Category
// @accept json
// @produce json
// @security BearerAuth
// @param id path int true "Category ID"
// @param request body object{nama_category=string} true "Category Data"
// @success 200 {object} Response
// @failure 400 {object} Response
// @failure 401 {object} Response
// @failure 404 {object} Response
// @failure 500 {object} Response
// @router /category/{id} [put]
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
// @summary Delete a category
// @description Delete a category by ID (Admin only).
// @tags Category
// @accept json
// @produce json
// @security BearerAuth
// @param id path int true "Category ID"
// @success 200 {object} Response
// @failure 400 {object} Response
// @failure 401 {object} Response
// @failure 404 {object} Response
// @failure 500 {object} Response
// @router /category/{id} [delete]
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
