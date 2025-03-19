package models

import (
	"time"
	"gorm.io/gorm"
)

type ResourceStatus string

const (
	Pending    ResourceStatus = "PENDING"
	InReview   ResourceStatus = "IN-REVIEW"
	Rejected   ResourceStatus = "REJECTED"
	Approved   ResourceStatus = "APPROVED"
)

type Resource struct {
	ID                           int            `gorm:"primaryKey;autoIncrement"`
	StartDate                    time.Time      `gorm:"type:date"`
	EndDate                      time.Time      `gorm:"type:date"`
	Status                       ResourceStatus `gorm:"type:enum('PENDING','IN-REVIEW','REJECTED','APPROVED')"`
	IPAddress                    string         `gorm:"type:text"`
	Views                        int            `gorm:"default:0"`
	Pages                        int            `gorm:"default:0"`
	Likes                        int            `gorm:"default:0"`
	Progress                     string         `gorm:"type:text"`
	DurationInSeconds            string         `gorm:"type:text"`
	Finished                     string         `gorm:"type:text"`
	RecordedDate                 string         `gorm:"type:text"`
	ResponseID                   string         `gorm:"type:text"`
	ExternalReference            string         `gorm:"type:text"`
	LocationLatitude             string         `gorm:"type:text"`
	LocationLongitude            string         `gorm:"type:text"`
	DistributionChannel          string         `gorm:"type:text"`
	UserLanguage                 string         `gorm:"type:text"`
	ContentTitle                 string         `gorm:"type:text"`
	ContentOrganization          string         `gorm:"type:text"`
	UserID                       int            `gorm:"not null;index"`
	ResourceTitle                string         `gorm:"type:text"`
	ResourceDescription          string         `gorm:"type:text"`
	ResourceTypeOptionId         string         `gorm:"type:text"`
	OtherTitle                   string         `gorm:"type:text"`
	ResourceCategoryOptionId     string         `gorm:"type:text"`
	ResourceIdentityGroupId      string         `gorm:"type:text;not null"`
	ResourceTypeOtherText        string         `gorm:"type:text"`
	TargetAudience               string         `gorm:"type:text"`
	ResourceSuppose              string         `gorm:"type:text"`
	ResourceSupposeOtherTitle    string         `gorm:"type:text"`
	SustainableDevelopmentGoals  string         `gorm:"type:text"`
	RacialEqualityConsciousness  string         `gorm:"type:text"`
	YearInitiated                int            `gorm:"default:null"`
	ResourceStatus               string         `gorm:"type:text"`
	ResourceStartDate            *time.Time     `gorm:"type:date;default:null"`
	ResourceEndDate              *time.Time     `gorm:"type:date;default:null"`
	ResourceLink                 string         `gorm:"type:text"`
	File                         string         `gorm:"type:text"`
	FileName                     string         `gorm:"type:text"`
	FileSize                     int            `gorm:"default:null"`
	FileType                     string         `gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
