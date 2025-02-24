package routes

import (
	"github.com/habbazettt/evermos-service-go/controllers"
	"github.com/habbazettt/evermos-service-go/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {
	user := app.Group("/api/v1/user", middleware.JWTMiddleware())

	user.Get("/", controllers.GetMyProfile)
	user.Put("/", controllers.UpdateProfile)
}
