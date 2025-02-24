package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/evermos-service-go/controllers"
	"github.com/habbazettt/evermos-service-go/middleware"
)

func AlamatRoutes(app *fiber.App) {
	alamat := app.Group("/api/v1/user/alamat", middleware.JWTMiddleware())

	alamat.Get("/", controllers.GetListAddress)
	alamat.Get("/:id", controllers.GetAlamatByID)
	alamat.Post("/", controllers.CreateAlamat)
	alamat.Put("/:id", controllers.UpdateAlamatByID)
	alamat.Delete("/:id", controllers.DeleteAlamatByID)
}
