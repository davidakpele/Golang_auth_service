package models

import (
	"time"
)

type Reports struct {
	ID         	uint           `gorm:"primaryKey;autoIncrement"`
	ResourceId 	uint      		`gorm:"not null;"` 
	Fullname   string      		`gorm:"type:text"`
	Email      string         `gorm:"type:text"`
	Content     string         `gorm:"type:text"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
}
