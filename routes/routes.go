package routes

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) {
	UserRoutes(app)
	AuthRoutes(app)
	LocationRoutes(app)
	AlamatRoutes(app)
	TokoRoutes(app)
	CategoryRoutes(app)
	ProductRoutes(app)
	TransactionRoutes(app)
}
