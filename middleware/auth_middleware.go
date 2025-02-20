package middleware

import (
	"errors"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// Middleware JWT untuk validasi token
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("token")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  false,
				"message": "Token tidak ditemukan",
				"errors":  nil,
				"data":    nil,
			})
		}

		// Parsing token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Metode signing tidak valid")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  false,
				"message": "Token tidak valid atau kadaluarsa",
				"errors":  err.Error(),
				"data":    nil,
			})
		}

		// Ambil claims dari token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  false,
				"message": "Gagal membaca token claims",
				"errors":  nil,
				"data":    nil,
			})
		}

		// Pastikan ada user_id dalam token
		userIDFloat, ok := claims["user_id"].(float64) // JWT menyimpan angka sebagai float64
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  false,
				"message": "User ID tidak ditemukan dalam token",
				"errors":  nil,
				"data":    nil,
			})
		}

		// Konversi ke uint dan simpan di Locals
		userID := uint(userIDFloat)
		c.Locals("user_id", userID)

		return c.Next()
	}
}

// ExtractUserID mengambil user_id dari context
func ExtractUserID(c *fiber.Ctx) (uint, error) {
	userIDInterface := c.Locals("user_id")
	if userIDInterface == nil {
		return 0, errors.New("user ID tidak ditemukan dalam token")
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		return 0, errors.New("gagal mengonversi user ID dari token")
	}

	return userID, nil
}
