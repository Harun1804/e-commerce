package access

import (
	"harun1804/e-commerce/modules/access/controllers"
	"harun1804/e-commerce/modules/access/repositories"

	"harun1804/e-commerce/modules/access/usecases"

	"gorm.io/gorm"
)

func BuildAccessModule(db *gorm.DB) (
	controllers.RoleControllerInterface,
	usecases.RoleUsecaseInterface,
) {
	roleRepo := repositories.NewRoleRepository(db)
	roleUsecase := usecases.NewRoleUsecase(roleRepo)
	roleController := controllers.NewRoleController(roleUsecase)

	return roleController, roleUsecase
}
