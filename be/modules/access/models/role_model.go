package models

import "time"

type Role struct {
	ID          uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string       `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Description string       `gorm:"type:text" json:"description,omitempty"`
	Permissions []Permission `gorm:"many2many:role_permissions;joinForeignKey:RoleId;joinReferences:PermissionId" json:"permissions,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}
