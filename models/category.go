package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID            uint           `gorm:"primaryKey;autoIncrement"`
	ResourceTitle string         `gorm:"size:255;not null"`
	ParentLevel   *uint          `gorm:"default:null"` 
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}
