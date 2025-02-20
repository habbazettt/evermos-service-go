package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/evermos-service-go/controllers"
	"github.com/habbazettt/evermos-service-go/middleware"
)

func AlamatRoutes(app *fiber.App) {
	alamat := app.Group("/api/v1/user/alamat", middleware.JWTMiddleware())

	alamat.Get("/", controllers.GetListAddress)         // Get all addresses
	alamat.Get("/:id", controllers.GetAlamatByID)       // Get address by ID
	alamat.Post("/", controllers.CreateAlamat)          // Create new address
	alamat.Put("/:id", controllers.UpdateAlamatByID)    // Update address by ID
	alamat.Delete("/:id", controllers.DeleteAlamatByID) // Delete address by ID
}
