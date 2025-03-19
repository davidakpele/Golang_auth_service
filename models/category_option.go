package models

import (
	"time"
)

type CategoryOption struct {
	ID         uint           `gorm:"primaryKey;autoIncrement"`
	CategoryID uint           `gorm:"not null;index"` 
	Title      string         `gorm:"size:255;not null"`
	Icon       string         `gorm:"type:text"` 
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
}
