package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/habbazettt/evermos-service-go/config"
	_ "github.com/habbazettt/evermos-service-go/docs"
	"github.com/habbazettt/evermos-service-go/routes"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()

	err = config.SetupCloudinary()
	if err != nil {
		log.Fatal("Cloudinary setup failed:", err)
	}

	app := fiber.New()

	// @title Evermos Store and Product API
	// @version 1.0
	// @description API documentation for Evermos service backend.
	// @termsOfService http://swagger.io/terms/
	// @contact.name API Support
	// @contact.url http://www.swagger.io/support
	// @host localhost:8080
	// @BasePath /api/v1
	// @securityDefinitions.apikey BearerAuth
	// @in header
	// @name token
	// @description Enter your token in the format: Bearer <token>
	// @security BearerAuth
	app.Get("/swagger/*", swagger.HandlerDefault)

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":8080"))
	fmt.Println("Server started on port 8080")
}
