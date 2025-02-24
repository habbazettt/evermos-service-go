package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/evermos-service-go/services"
)

// Get List of Provinces
// @Summary Get List of Provinces
// @Description Get a list of provinces with optional search, limit, and pagination.
// @Tags Location
// @Accept json
// @Produce json
// @Param search query string false "Search province by name"
// @Param limit query int false "Limit results per page" default(10)
// @Param page query int false "Page number" default(1)
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router /provcity/listprovincies [get]
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

// Get Province Detail
// @Summary Get Province Detail
// @Description Get detailed information of a specific province by ID.
// @Tags Location
// @Accept json
// @Produce json
// @Param prov_id path string true "Province ID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Router /provcity/detailprovince/{prov_id} [get]
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

// Get List of Cities
// @Summary Get List of Cities
// @Description Get a list of cities in a specific province.
// @Tags Location
// @Accept json
// @Produce json
// @Param prov_id path string true "Province ID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Router /provcity/listcities/{prov_id} [get]
func GetListCities(c *fiber.Ctx) error {
	provinceID := c.Params("prov_id")

	// Ambil data kota berdasarkan provinsi
	cities, err := services.GetCitiesByProvince(provinceID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to fetch cities",
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
// @Summary Get City Detail
// @Description Get detailed information of a specific city by ID.
// @Tags Location
// @Accept json
// @Produce json
// @Param city_id path string true "City ID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Router /provcity/detailcity/{city_id} [get]
func GetCityDetail(c *fiber.Ctx) error {
	cityID := c.Params("city_id")

	// Ambil data kota berdasarkan city ID
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
