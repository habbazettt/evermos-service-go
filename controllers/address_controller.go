package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/evermos-service-go/middleware"
	"github.com/habbazettt/evermos-service-go/models"
	"github.com/habbazettt/evermos-service-go/services"
)

// Get All Addresses
// @Summary Get All Addresses
// @Description Get all addresses for the authenticated user.
// @Tags Address
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param judul_alamat query string false "Filter by address title"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Router /user/alamat [get]
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

	judulAlamat := c.Query("judul_alamat")

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

// Get Address by ID
// @Summary Get Address by ID
// @Description Get a specific address for the authenticated user.
// @Tags Address
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Address ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /user/alamat/{id} [get]
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

	alamatID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid address ID",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	alamat, err := services.GetAlamatByID(userID, uint(alamatID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Address not found",
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

// Create Address
// @Summary Create a new address
// @Description Create a new address for the authenticated user.
// @Tags Address
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.Alamat true "Address Data"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Router /user/alamat [post]
func CreateAlamat(c *fiber.Ctx) error {
	userID, err := middleware.ExtractUserID(c)
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

	alamat.IDUser = userID

	err = services.CreateAlamat(alamat)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to create address",
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  true,
		"message": "Address created successfully",
		"data":    alamat,
	})
}

// Update Address
// @Summary Update Address by ID
// @Description Update an existing address for the authenticated user.
// @Tags Address
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Address ID"
// @Param request body models.Alamat true "Updated Address Data"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /user/alamat/{id} [put]
func UpdateAlamatByID(c *fiber.Ctx) error {
	alamatIDStr := c.Params("id")
	alamatID, err := strconv.ParseUint(alamatIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid address ID",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	userID, err := middleware.ExtractUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Unauthorized",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	var alamatRequest models.Alamat
	if err := c.BodyParser(&alamatRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid request body",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	err = services.UpdateAlamatByID(uint(alamatID), userID, map[string]interface{}{
		"judul_alamat": alamatRequest.JudulAlamat,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to update address",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Address updated successfully",
	})
}

// Delete Address
// @Summary Delete Address by ID
// @Description Delete a specific address for the authenticated user.
// @Tags Address
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Address ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /user/alamat/{id} [delete]
func DeleteAlamatByID(c *fiber.Ctx) error {
	alamatIDStr := c.Params("id") // Ambil ID dari parameter URL
	alamatID, err := strconv.ParseUint(alamatIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid address ID",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	userID, err := middleware.ExtractUserID(c) // Ambil user_id dari middleware
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Unauthorized",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	err = services.DeleteAlamatByID(uint(alamatID), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to delete address",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Address deleted successfully",
		"data":    nil,
	})
}
