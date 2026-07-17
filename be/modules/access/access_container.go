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
	controllers.UserControllerInterface,
	usecases.UserUsecaseInterface,
) {
	roleRepo := repositories.NewRoleRepository(db)
	permissionRepo := repositories.NewPermissionRepository(db)
	rolePermissionRepo := repositories.NewRolePermissionRepository(db)
	userRepo := repositories.NewUserRepository(db)
	roleUserRepo := repositories.NewRoleUserRepository(db)

	permissionUsecase := usecases.NewPermissionUsecase(permissionRepo)
	roleUsecase := usecases.NewRoleUsecase(roleRepo, permissionUsecase, rolePermissionRepo)
	userUsecase := usecases.NewUserUsecase(userRepo, roleUsecase, roleUserRepo)

	roleController := controllers.NewRoleController(roleUsecase)
	permissionController := controllers.NewPermissionController(permissionUsecase)
	userController := controllers.NewUserController(userUsecase)

	return roleController, roleUsecase, permissionController, permissionUsecase, userController, userUsecase
}
