package models

import (
	"time"
)

type Messages struct {
	ID         	uint           `gorm:"primaryKey;autoIncrement"`
	SenderId 	uint      `gorm:"not null;"` 
	ReceiverId   uint      `gorm:"not null"`
	Content     string         `gorm:"type:text"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
}
