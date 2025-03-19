package services

import (
	"api-service/models"
	"api-service/repositories"
// 	"api-service/requests"
// 	"api-service/utils"
 	"errors"
)

type AdminService struct {
	adminRepo repositories.AdminRepository
}

func NewAdminService(adminRepo repositories.AdminRepository) *AdminService {
	return &AdminService{adminRepo: adminRepo}
}

func (s *AdminService) GetPendingArticles() ([]models.Resource, error) {
	// Fetch structured response from repository
	data, err := s.adminRepo.GetPendingArticles()
	if err != nil {
		return nil, errors.New("failed to fetch pending articles")
	}

	// Extract recent resources from the map
	resources, ok := data["resources"].(map[string]interface{})["recent_resources"].([]models.Resource)
	if !ok {
		return nil, errors.New("failed to parse recent resources")
	}

	return resources, nil
}
