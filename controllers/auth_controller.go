package controllers

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/habbazettt/evermos-service-go/config"
	"github.com/habbazettt/evermos-service-go/models"
	"github.com/habbazettt/evermos-service-go/services"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
	Data    interface{} `json:"data"`
}

type LoginRequest struct {
	NoTelp    string `json:"no_telp"`
	KataSandi string `json:"kata_sandi"`
}

type RegisterRequest struct {
	Nama         string `json:"nama"`
	KataSandi    string `json:"kata_sandi"`
	NoTelp       string `json:"no_telp"`
	Email        string `json:"email"`
	TanggalLahir string `json:"tanggal_lahir"`
	JenisKelamin string `json:"jenis_kelamin"`
	Pekerjaan    string `json:"pekerjaan"`
	Tentang      string `json:"tentang"`
	IsAdmin      bool   `json:"is_admin"`
	IDProvinsi   string `json:"id_provinsi"`
	IDKota       string `json:"id_kota"`
}

type LoginResponse struct {
	Nama         string            `json:"nama"`
	NoTelp       string            `json:"no_telp"`
	TanggalLahir string            `json:"tanggal_Lahir"`
	Tentang      string            `json:"tentang"`
	Pekerjaan    string            `json:"pekerjaan"`
	JenisKelamin string            `json:"jenis_kelamin"`
	Email        string            `json:"email"`
	IsAdmin      bool              `json:"is_admin"`
	Provinsi     services.Province `json:"id_provinsi"`
	Kota         services.City     `json:"id_kota"`
	Token        string            `json:"token"`
}

// Register - Register a new user
// @Summary Register a new user to the system
// @Description Register a new user with the provided details (name, phone number, email, password, etc.)
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param request body RegisterRequest true "Register Request Body"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /auth/register [post]
func Register(c *fiber.Ctx) error {
	var data RegisterRequest
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid request",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	var existingUser models.User
	if err := config.DB.Where("no_telp = ?", data.NoTelp).Or("email = ?", data.Email).First(&existingUser).Error; err == nil {
		var errorMessage string
		if existingUser.NoTelp == data.NoTelp {
			errorMessage = "Phone number already exists"
		} else if existingUser.Email == data.Email {
			errorMessage = "Email already exists"
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{errorMessage},
			"data":    nil,
		})
	}

	// Hash kata sandi sebelum disimpan
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(data.KataSandi), bcrypt.DefaultCost)

	// Simpan user ke database
	user := models.User{
		Nama:         data.Nama,
		KataSandi:    string(hashedPassword),
		NoTelp:       data.NoTelp,
		Email:        data.Email,
		TanggalLahir: data.TanggalLahir,
		JenisKelamin: data.JenisKelamin,
		Pekerjaan:    data.Pekerjaan,
		IsAdmin:      data.IsAdmin,
		Tentang:      data.Tentang,
		IDProvinsi:   data.IDProvinsi,
		IDKota:       data.IDKota,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to create user",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Buat toko otomatis setelah register
	toko := models.Toko{IDUser: user.ID, NamaToko: user.Nama + " Store"}
	config.DB.Create(&toko)

	// Fetch provinsi & kota dari services
	province, err := services.GetAllProvinces()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to fetch province data",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	city, err := services.GetCitiesByProvince(user.IDProvinsi)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to fetch city data",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Cari provinsi yang sesuai dengan IDProvinsi user
	var selectedProvince services.Province
	for _, p := range province {
		if p.ID == user.IDProvinsi {
			selectedProvince = p
			break
		}
	}

	// Cari kota yang sesuai dengan IDKota user
	var selectedCity services.City
	for _, c := range city {
		if c.ID == user.IDKota {
			selectedCity = c
			break
		}
	}

	// Response sukses tanpa token
	response := fiber.Map{
		"nama":          user.Nama,
		"no_telp":       user.NoTelp,
		"tanggal_Lahir": user.TanggalLahir,
		"tentang":       user.Tentang,
		"pekerjaan":     user.Pekerjaan,
		"email":         user.Email,
		"is_admin":      user.IsAdmin,
		"id_provinsi":   selectedProvince,
		"id_kota":       selectedCity,
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "User registered successfully",
		"errors":  nil,
		"data":    response,
	})
}

// Login - User login
// @Summary Login a user with phone number and password
// @Description Login a user and return a JWT token along with user details.
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param request body LoginRequest true "Login Request Body"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /auth/login [post]
func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid request body",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Cek apakah user ada di database
	var user models.User
	result := config.DB.Where("no_telp = ?", req.NoTelp).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  false,
				"message": "User not found",
				"errors":  nil,
				"data":    nil,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to fetch user",
			"errors":  result.Error.Error(),
			"data":    nil,
		})
	}

	// Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.KataSandi), []byte(req.KataSandi)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid password",
			"errors":  nil,
			"data":    nil,
		})
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"user_id": float64(user.ID), // Pastikan sebagai float64
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to generate token",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Fetch provinsi & kota dari services
	province, err := services.GetAllProvinces()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to fetch province data",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	city, err := services.GetCitiesByProvince(user.IDProvinsi)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to fetch city data",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Cari provinsi yang sesuai dengan IDProvinsi user
	var selectedProvince services.Province
	for _, p := range province {
		if p.ID == user.IDProvinsi {
			selectedProvince = p
			break
		}
	}

	// Cari kota yang sesuai dengan IDKota user
	var selectedCity services.City
	for _, c := range city {
		if c.ID == user.IDKota {
			selectedCity = c
			break
		}
	}

	// Response sukses
	response := LoginResponse{
		Nama:         user.Nama,
		NoTelp:       user.NoTelp,
		TanggalLahir: user.TanggalLahir,
		Tentang:      user.Tentang,
		Pekerjaan:    user.Pekerjaan,
		Email:        user.Email,
		IsAdmin:      user.IsAdmin,
		Provinsi:     selectedProvince,
		Kota:         selectedCity,
		Token:        signedToken,
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Login successful",
		"errors":  nil,
		"data":    response,
	})
}
