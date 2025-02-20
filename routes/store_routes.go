package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/evermos-service-go/controllers"
	"github.com/habbazettt/evermos-service-go/middleware"
)

func TokoRoutes(app *fiber.App) {
	toko := app.Group("/api/v1/toko", middleware.JWTMiddleware())

	toko.Get("/", controllers.GetAllStores)
	toko.Get("/my", controllers.GetMyStore)
	toko.Get("/:id", controllers.GetStoreByID)
	toko.Post("/", controllers.CreateStore)
	toko.Put("/:id", controllers.UpdateStore)
	toko.Delete("/:id", controllers.DeleteStore)
}	
