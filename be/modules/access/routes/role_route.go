package routes

import (
	"harun1804/e-commerce/modules/access/controllers"

	"github.com/gofiber/fiber/v3"
)

func RoleRoutes(router fiber.Router, roleController controllers.RoleControllerInterface) {
	roleRoutes := router.Group("/roles")

	roleRoutes.Get("/", roleController.GetAllRoles)
	roleRoutes.Get("/:id", roleController.GetRoleByID)
	roleRoutes.Post("/", roleController.CreateRole)
	roleRoutes.Put("/:id", roleController.UpdateRole)
	roleRoutes.Delete("/:id", roleController.DeleteRole)
}