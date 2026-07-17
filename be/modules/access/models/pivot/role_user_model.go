package pivot

import "time"

type RoleUser struct {
	Id           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleId       uint      `gorm:"not null;uniqueIndex:idx_role_user" json:"role_id"`
	UserId       uint      `gorm:"not null;uniqueIndex:idx_role_user" json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
