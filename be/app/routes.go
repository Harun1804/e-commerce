package app

import (
	accessRoutes "harun1804/e-commerce/modules/access/routes"

	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App, container *Containers) {
	api := app.Group("/api")

	// Access module routes
	access := api.Group("/access")
	accessRoutes.RoleRoutes(access, container.RoleController)
	accessRoutes.PermissionRoutes(access, container.PermissionController)
	accessRoutes.UserRoutes(access, container.UserController)
}
