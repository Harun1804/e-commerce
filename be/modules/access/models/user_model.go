package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"type:varchar(255)" json:"password"`
	Roles     []Role    `gorm:"many2many:role_users;joinForeignKey:UserId;References:ID" json:"roles,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
