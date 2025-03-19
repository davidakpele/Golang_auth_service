package models

import (
	"time"
)

type Bookmark struct {
	ID         uint           `gorm:"primaryKey;autoIncrement"`
	UserId 			uint      `gorm:"not null;"` 
	ResourceId      string    `gorm:"size:255;not null"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
}
