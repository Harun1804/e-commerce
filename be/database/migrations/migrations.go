package migrations

import (
	"harun1804/e-commerce/configs"
	accessModel "harun1804/e-commerce/modules/access/models"

	"go.uber.org/zap"
)

func RunMigrations() {
	err := configs.DB.AutoMigrate(
		&accessModel.Role{},
		&accessModel.Permission{},
	)

	if err != nil {
		zap.L().Fatal("Failed to run migrations", zap.Error(err))
	}
}
