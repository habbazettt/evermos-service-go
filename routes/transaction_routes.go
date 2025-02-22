package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/evermos-service-go/controllers"
	"github.com/habbazettt/evermos-service-go/middleware"
)

func TransactionRoutes(app *fiber.App) {
	transaction := app.Group("/api/v1/trx", middleware.JWTMiddleware())

	transaction.Get("/", controllers.GetAllTransactions)
	transaction.Get("/:id", controllers.GetTransactionByID)
	transaction.Post("/", controllers.CreateTransaction)
}
