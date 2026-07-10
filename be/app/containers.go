package app

import (
	"harun1804/e-commerce/configs"
	accessMod "harun1804/e-commerce/modules/access"
	accessCtrl "harun1804/e-commerce/modules/access/controllers"

	"gorm.io/gorm"
)

type Containers struct {
	RoleController accessCtrl.RoleControllerInterface
}

func BuildContainers() *Containers {
	return NewContainers(configs.DB)
}

func NewContainers(db *gorm.DB) *Containers {
	roleController, _ := accessMod.BuildAccessModule(db)
	return &Containers{
		RoleController: roleController,
	}
}