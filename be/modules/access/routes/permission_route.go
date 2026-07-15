package routes

import (
	"harun1804/e-commerce/modules/access/controllers"

	"github.com/gofiber/fiber/v3"
)

func PermissionRoutes(router fiber.Router, permissionController controllers.PermissionControllerInterface) {
	permissionRoutes := router.Group("/permissions")

	permissionRoutes.Get("/", permissionController.GetAllPermissions)
	permissionRoutes.Get("/:id", permissionController.GetPermissionByID)
	permissionRoutes.Post("/", permissionController.CreatePermission)
	permissionRoutes.Put("/:id", permissionController.UpdatePermission)
	permissionRoutes.Delete("/:id", permissionController.DeletePermission)
}
