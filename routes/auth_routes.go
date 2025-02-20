package routes

import (
	"github.com/habbazettt/evermos-service-go/controllers"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	route := app.Group("/api/v1/auth")

	route.Post("/register", controllers.Register)
	route.Post("/login", controllers.Login)
}
