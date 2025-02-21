package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/evermos-service-go/controllers"
	"github.com/habbazettt/evermos-service-go/middleware"
)

func ProductRoutes(app *fiber.App) {
	product := app.Group("/api/v1/product", middleware.JWTMiddleware())

	product.Get("/", controllers.GetAllProducts)
	product.Get("/:id", controllers.GetProductByID)
	product.Post("/", controllers.CreateProduct)
	product.Put("/:id", controllers.UpdateProduct)
	product.Delete("/:id", controllers.DeleteProduct)
}
