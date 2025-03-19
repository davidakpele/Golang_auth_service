package repositories

import (
	"api-service/models"
	// "log"
	"gorm.io/gorm"

)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) GetPendingArticles() (map[string]interface{}, error) {
	var (
		totalPending   int64
		totalReview    int64
		totalRejected  int64
		totalApproved  int64
		totalUsers     int64
		totalResources int64
		recentResources []models.Resource
		recentUsers []models.User
	)

	// Fetch total pending resources
	r.db.Model(&models.Resource{}).Where("status = ?", "PENDING").Count(&totalPending)

	// Fetch total in-review resources
	r.db.Model(&models.Resource{}).Where("status = ?", "IN-REVIEW").Count(&totalReview)

	// Fetch total rejected resources
	r.db.Model(&models.Resource{}).Where("status = ?", "REJECTED").Count(&totalRejected)

	// Fetch total approved resources
	r.db.Model(&models.Resource{}).Where("status = ?", "APPROVED").Count(&totalApproved)

	// Fetch total users
	r.db.Model(&models.User{}).Count(&totalUsers)

	// Fetch total resources
	r.db.Model(&models.Resource{}).Count(&totalResources)

	// Fetch last 10 recent posted resources
	err := r.db.Preload("User").Order("startdate DESC").Limit(10).Find(&recentResources).Error
	if err != nil {
		return nil, err
	}

	// Fetch last 10 recent users
	err = r.db.Order("created_at DESC").Limit(10).Find(&recentUsers).Error
	if err != nil {
		return nil, err
	}

	// Return structured response
	response := map[string]interface{}{
		"info": map[string]int64{
			"total_pending_resources":  totalPending,
			"total_approved_resources": totalApproved,
			"total_rejected_resources": totalRejected,
			"total_in_review_resources": totalReview,
			"total_no_of_users": totalUsers,
			"total_no_of_resources": totalResources,
		},
		"resources": map[string]interface{}{
			"recent_resources": recentResources,
			"recent_users": recentUsers,
		},
	}

	return response, nil
}
