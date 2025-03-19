package models

import (
	"time"
)

type Comments struct {
	ID         	uint           `gorm:"primaryKey;autoIncrement"`
	UserId 		uint      `gorm:"not null;"` 
	ParentId    string      `gorm:"size:255;not null"`
	Content     string         `gorm:"type:text"`
	PressId		string         `gorm:"type:text"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
}
