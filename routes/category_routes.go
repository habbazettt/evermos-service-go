package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/evermos-service-go/controllers"
	"github.com/habbazettt/evermos-service-go/middleware"
)

func CategoryRoutes(app *fiber.App) {
	category := app.Group("/api/v1/category")

	category.Get("/", controllers.GetAllCategories)
	category.Get("/:id", controllers.GetCategoryByID)

	category.Post("/", middleware.JWTMiddleware(), middleware.AdminMiddleware(), controllers.CreateCategory)
	category.Put("/:id", middleware.JWTMiddleware(), middleware.AdminMiddleware(), controllers.UpdateCategory)
	category.Delete("/:id", middleware.JWTMiddleware(), middleware.AdminMiddleware(), controllers.DeleteCategory)
}
