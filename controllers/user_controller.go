package controllers

import (
	"strconv"

	"github.com/habbazettt/evermos-service-go/middleware"
	"github.com/habbazettt/evermos-service-go/services"

	"github.com/gofiber/fiber/v2"
)

// GetMyProfile handler
func GetMyProfile(c *fiber.Ctx) error {
	// Ambil user_id dari token JWT
	userID, err := middleware.ExtractUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Unauthorized",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Panggil service untuk mendapatkan data user
	user, err := services.GetUserByID(strconv.Itoa(int(userID)))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "User tidak ditemukan",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Berikan respons dengan data user
	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Success",
		"errors":  nil,
		"data":    user,
	})
}

// UpdateProfile handler
func UpdateProfile(c *fiber.Ctx) error {
	// Ambil user_id dari token JWT
	userID, err := middleware.ExtractUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Unauthorized",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Parse request body
	var updateData services.UpdateUserRequest
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid request data",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Ambil data user saat ini dari database
	currentUser, err := services.GetUserByID(strconv.Itoa(int(userID)))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "User not found",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Jika user bukan admin, abaikan perubahan is_admin
	if !currentUser.IsAdmin {
		updateData.IsAdmin = currentUser.IsAdmin
	}

	// Panggil service untuk update user
	updatedUser, err := services.UpdateUserByID(strconv.Itoa(int(userID)), updateData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to update profile",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Berikan respons dengan data user yang telah diperbarui
	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Profile updated successfully",
		"errors":  nil,
		"data":    updatedUser,
	})
}
