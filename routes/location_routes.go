package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/evermos-service-go/controllers"
)

func LocationRoutes(app *fiber.App) {
	location := app.Group("/api/v1/provcity")

	location.Get("/listprovincies", controllers.GetListProvinces)
	location.Get("/detailprovince/:prov_id", controllers.GetProvinceDetail)
	location.Get("/listcities/:prov_id", controllers.GetListCities)
	location.Get("/detailcity/:city_id", controllers.GetCityDetail)

}
