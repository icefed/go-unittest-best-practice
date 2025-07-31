package store

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string    `gorm:"primaryKey"`
	Name      string    `gorm:"size:128;not null"`
	Email     string    `gorm:"size:128;uniqueIndex;not null"`
	Password  string    `gorm:"size:256;not null"`
	Age       int       `gorm:"default:0"`
	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;autoUpdateTime"`
	DeletedAt gorm.DeletedAt
}
