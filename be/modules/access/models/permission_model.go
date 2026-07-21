package models

import "time"

type Permission struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(100);;not null" json:"name"`
	Slug        string    `gorm:"type:text;uniqueIndex;not null" json:"slug"`
	Module      string    `gorm:"type:varchar(100);not null" json:"module"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
