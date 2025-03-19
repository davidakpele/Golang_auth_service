package services

import (
	"api-service/repositories"
	"api-service/requests"
	"api-service/models"
	"errors"
)

type ResourceService struct {
	resourceRepo repositories.ResourceRepository
}

func NewResourceService(resourceRepo repositories.ResourceRepository) *ResourceService {
	return &ResourceService{resourceRepo: resourceRepo}
}

// CreateResource processes the resource creation request
func (s *ResourceService) CreateResource(req requests.CreateResourceRequest, files []requests.UploadedFile) (string, error) {
	// Ensure files are uploaded
	if len(files) == 0 {
		return "", errors.New("at least one file must be uploaded")
	}

	// Create resource data structure
	resource := map[string]interface{}{
		"resourceType":         req.ResourceType,
		"startDate":            req.StartDate,
		"endDate":              req.EndDate,
		"resourceContactTitle": req.ResourceContactTitle,
		"resourceOrganization": req.ResourceOrganization,
		"resourceTitle":        req.ResourceTitle,
		"resourceDescription":  req.ResourceDescription,
		"resourceCategory":     req.ResourceCategory,
		"identityGroup":        req.IdentityGroup,
		"racialSphere":         req.RacialSphere,
		"sustainable":          req.Sustainable,
		"targetAudience":       req.TargetAudience,
		"yearResource":         req.YearResource,
		"status":               req.Status,
		"weblink":              req.Weblink,
		"files":                files,
	}

	// Save resource to database via repository
	resourceID, err := s.resourceRepo.SaveResource(resource)
	if err != nil {
		return "", err
	}

	return resourceID, nil
}

// GetAllResources fetches all resources with pagination
func (rs *ResourceService) GetAllResources(offset, itemsPerPage int) (map[string]interface{}, error) {
	return rs.resourceRepo.GetAllResources(offset, itemsPerPage)
}

// GetResourceByID fetches a single resource by ID
func (rs *ResourceService) GetResourceByID(id int) (map[string]interface{}, error) {
	return rs.resourceRepo.GetResourceByID(id)
}

func (rs *ResourceService) DeleteResource(id int) (bool, error) {
	err := rs.resourceRepo.DeleteResource(id)
	if err != nil {
		return false, err 
	}
	return true, nil
}

func (rs *ResourceService) GetLimitedResource(offset, itemsPerPage int) (map[string]interface{}, error) {
	return rs.resourceRepo.GetLimitedResource(offset, itemsPerPage)
}

func (rs *ResourceService) CreateViewOnResource(resourceID int) (bool, error) {
	// Delegate the task to the repository to create a view for the resource
	return rs.resourceRepo.CreateViewOnResource(resourceID)
}

func (s *ResourceService) CreateLikeOnResource(resourceID int) (bool, error) {
	return s.resourceRepo.CreateLikeOnResource(resourceID)
}

// UpdateResourceStatus updates the status of a resource
func (s *ResourceService) UpdateResourceStatus(id int, status string) error {
	// Allowed statuses
	allowedStatuses := map[string]bool{
		"PENDING":   true,
		"IN-REVIEW": true,
		"APPROVED":  true,
		"REJECTED":  true,
	}

	// Validate status
	if !allowedStatuses[status] {
		return errors.New("Invalid status value. Allowed values: PENDING, IN-REVIEW, APPROVED, REJECTED")
	}

	// Check if resource exists
	_, err := s.resourceRepo.GetResourceByID(id)
	if err != nil {
		return errors.New("Resource not found")
	}

	// Update resource status
	err = s.resourceRepo.UpdateResourceStatus(id, status)
	if err != nil {
		return errors.New("Update failed or resource not found")
	}

	return nil
}

// GetUserArticles fetches articles for a specific user
func (s *ResourceService) GetUserArticles(userID int) ([]models.Resource, error) {
	return s.resourceRepo.GetUserArticles(userID)
}

func (s *ResourceService) TotalUserArticles(userID int) (int, error) {
	return s.resourceRepo.TotalUserArticles(userID)
}

// GetRandomResources fetches random resources with pagination
// func (s *ResourceService) GetRandomResources(offset, itemsPerPage int) ([]models.Resource, int, error) {
// 	resources, totalCount, err := s.resourceRepo.GetRandomResources(offset, itemsPerPage)
// 	if err != nil {
// 		return nil, 0, errors.New("failed to fetch random resources")
// 	}
// 	return resources, totalCount, nil
// }