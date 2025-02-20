package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/evermos-service-go/services"
)

// Get List Provinces Handler
func GetListProvinces(c *fiber.Ctx) error {
	search := c.Query("search", "")
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	provinces, err := services.GetListProvinces(search, limit, page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to fetch province data",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Success",
		"errors":  nil,
		"data":    provinces,
	})
}

// Get Detail Province Handler
func GetProvinceDetail(c *fiber.Ctx) error {
	provinceID := c.Params("prov_id")

	province, err := services.GetProvinceDetail(provinceID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Province not found",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Success",
		"errors":  nil,
		"data":    province,
	})
}

// Get List Cities Handler
func GetListCities(c *fiber.Ctx) error {
	provinceID := c.Params("prov_id")
	search := c.Query("search", "")
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	cities, err := services.GetListCities(provinceID, search, limit, page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to fetch city data",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Success",
		"errors":  nil,
		"data":    cities,
	})
}

// Get Detail City Handler
func GetCityDetail(c *fiber.Ctx) error {
	cityID := c.Params("city_id")

	// Validasi jika city_id kosong
	if cityID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "City ID is required",
			"errors":  nil,
			"data":    nil,
		})
	}

	// Fetch data city berdasarkan city_id
	city, err := services.GetCityDetailByID(cityID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "City not found",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Success",
		"errors":  nil,
		"data":    city,
	})
}
