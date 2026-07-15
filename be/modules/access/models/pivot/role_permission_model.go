package pivot

import "time"

type RolePermission struct {
	Id           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleId       uint      `gorm:"not null;uniqueIndex:idx_role_permission" json:"role_id"`
	PermissionId uint      `gorm:"not null;uniqueIndex:idx_role_permission" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
