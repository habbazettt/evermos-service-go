package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Secret Key untuk JWT
var jwtSecret = []byte("mysecretkey")

// GenerateJWT membuat token JWT
func GenerateJWT(email string, userID int) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"id":    userID,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
