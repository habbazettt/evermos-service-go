package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/habbazettt/evermos-service-go/config"
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

	routes.SetupRoutes(app)

	app.Get("/swagger/*", swagger.HandlerDefault)

	log.Fatal(app.Listen(":8080"))
	fmt.Println("Server started on port 8080")
}
