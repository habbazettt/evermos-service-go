package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/evermos-service-go/config"
	"github.com/habbazettt/evermos-service-go/models"
)

// Middleware untuk mengecek apakah user adalah admin
func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := ExtractUserID(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  false,
				"message": "Unauthorized",
				"errors":  err.Error(),
			})
		}

		// Cek apakah user adalah admin
		var user models.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  false,
				"message": "User tidak ditemukan",
				"errors":  err.Error(),
			})
		}

		if !user.IsAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  false,
				"message": "Forbidden: Anda bukan admin",
			})
		}

		return c.Next()
	}
}
