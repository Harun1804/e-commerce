package app

import (
	"github.com/gofiber/fiber/v3"
	accessRoutes "harun1804/e-commerce/modules/access/routes"
)

func SetupRoutes(app *fiber.App, container *Containers) {
	api := app.Group("/api")

	// Access module routes
	access := api.Group("/access")
	accessRoutes.RoleRoutes(access, container.RoleController)
}