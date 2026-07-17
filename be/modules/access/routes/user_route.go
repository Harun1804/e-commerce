package routes

import (
	"harun1804/e-commerce/modules/access/controllers"

	"github.com/gofiber/fiber/v3"
)

func UserRoutes(router fiber.Router, userController controllers.UserControllerInterface) {
	userRoutes := router.Group("/users")

	userRoutes.Get("/", userController.GetAllUsers)
	userRoutes.Get("/:id", userController.GetUserByID)
	userRoutes.Post("/", userController.CreateUser)
	userRoutes.Put("/:id", userController.UpdateUser)
	userRoutes.Delete("/:id", userController.DeleteUser)
}
