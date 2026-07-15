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
	controllers.PermissionControllerInterface,
	usecases.PermissionUsecaseInterface,
) {
	roleRepo := repositories.NewRoleRepository(db)
	permissionRepo := repositories.NewPermissionRepository(db)
	rolePermissionRepo := repositories.NewRolePermissionRepository(db)

	roleUsecase := usecases.NewRoleUsecase(roleRepo, permissionRepo, rolePermissionRepo)
	permissionUsecase := usecases.NewPermissionUsecase(permissionRepo)
	
	roleController := controllers.NewRoleController(roleUsecase)
	permissionController := controllers.NewPermissionController(permissionUsecase)

	return roleController, roleUsecase, permissionController, permissionUsecase
}
