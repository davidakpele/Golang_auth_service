package repositories

import (
	"encoding/json"
	"errors"
	"strings"
	"gorm.io/gorm"
	"api-service/models"
)

type ResourceRepository struct {
	db *gorm.DB
}

func NewResourceRepository(db *gorm.DB) *ResourceRepository {
	return &ResourceRepository{db: db}
}

// Comment represents a comment entity
type Comment struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	ResourceID    int       `json:"resource_id"`
	ParentID      *int      `json:"parent_id"`
	Content       string    `json:"content"`
	CreatedAt     string    `json:"created_at"`
	CommenterName string    `json:"commenter_name"`
	CommenterEmail string   `json:"commenter_email"`
	Replies       []Comment `json:"replies"`
}

// SaveResource inserts a new resource into the database
func (r *ResourceRepository) SaveResource(resource map[string]interface{}) (string, error) {
	// Insert resource into database
	result := r.db.Table("resources").Create(resource)
	if result.Error != nil {
		return "", result.Error
	}

	// Check if an ID was generated
	id, ok := resource["id"].(string)
	if !ok {
		return "", errors.New("failed to retrieve resource ID")
	}

	return id, nil
}

// GetAllResources retrieves paginated resources along with user details
func (r *ResourceRepository) GetAllResources(offset, itemsPerPage int) (map[string]interface{}, error) {
	var results []map[string]interface{}
	query := `
		SELECT r.*, u.id as user_id, u.email, u.fullname 
		FROM resources r 
		INNER JOIN users u ON r.user_id = u.id 
		ORDER BY r.resource_title ASC 
		LIMIT ? OFFSET ?
	`
	if err := r.db.Raw(query, itemsPerPage, offset).Scan(&results).Error; err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return map[string]interface{}{}, nil
	}

	fieldsToConvert := []string{
		"resource_category",
		"resource_type",
		"sustainable_development_goals",
		"racial_equality_consciousness",
		"resource_identity_group",
	}

	for i := range results {
		// Convert CSV fields to slice
		for _, field := range fieldsToConvert {
			if val, ok := results[i][field].(string); ok && val != "" {
				results[i][field] = strings.Split(strings.TrimSpace(val), ",")
			} else {
				results[i][field] = []string{}
			}
		}

		// Decode JSON fields
		jsonFields := []string{"resource_title", "resource_description"}
		for _, field := range jsonFields {
			if val, ok := results[i][field].(string); ok && val != "" {
				var decoded interface{}
				if err := json.Unmarshal([]byte(val), &decoded); err == nil {
					results[i][field] = decoded
				}
			}
		}

		// Format user data
		results[i]["user"] = map[string]interface{}{
			"id":       results[i]["user_id"],
			"email":    results[i]["email"],
			"fullname": results[i]["fullname"],
		}
		delete(results[i], "user_id")
		delete(results[i], "email")
		delete(results[i], "fullname")
	}

	// Get total count of resources
	var totalCount int64
	r.db.Table("resources").Count(&totalCount)

	response := map[string]interface{}{
		"resources":            results,
		"total_found":          len(results),
		"total_number_of_data": totalCount,
		"current_page":         (offset / itemsPerPage) + 1,
		"per_page":             itemsPerPage,
	}

	return response, nil
}

func (rr *ResourceRepository) DeleteResource(id int) error {
	// Check if the resource exists before deleting
	var resource models.Resource
	if err := rr.db.First(&resource, id).Error; err != nil {
		return err // If not found, return the error
	}

	// Delete the resource
	return rr.db.Delete(&resource).Error
}

// GetLimitedResource fetches resources with pagination and applies necessary transformations
func (rr *ResourceRepository) GetLimitedResource(offset, itemsPerPage int) (map[string]interface{}, error) {
	var resources []models.Resource
	// Fetch the resources with pagination
	err := rr.db.Table("resources r").
		Joins("INNER JOIN users u ON r.user_id = u.id").
		Select("r.*, u.id as user_id, u.email, u.fullname").
		Order("r.resource_title ASC").
		Limit(itemsPerPage).
		Offset(offset).
		Find(&resources).Error

	if err != nil {
		return nil, err
	}

	// Process the results to convert fields as needed
	fieldsToConvert := []string{
		"resource_category_option_id",
		"resource_type_option_id",
		"sustainable_development_goals",
		"racial_equality_consciousness",
		"resource_identity_group_id",
	}

	for i := range resources {
		// Convert CSV fields to arrays
		for _, field := range fieldsToConvert {
			value := getFieldValue(&resources[i], field)
			if value != "" {
				// Split CSV string and filter empty values
				convertedValue := strings.FieldsFunc(value, func(r rune) bool {
					return r == ','
				})
				// Assign back to the struct
				setFieldValue(&resources[i], field, convertedValue)
			}
		}

		// Decode JSON fields (resource_title, resource_description)
		for _, field := range []string{"resource_title", "resource_description"} {
			value := getFieldValue(&resources[i], field)
			if value != "" {
				var decodedValue interface{}
				if err := json.Unmarshal([]byte(value), &decodedValue); err == nil {
					// Assign decoded value back
					setFieldValue(&resources[i], field, decodedValue)
				}
			}
		}
	}

	// Get the total number of resources
	var totalCount int64
	err = rr.db.Table("resources").Count(&totalCount).Error
	if err != nil {
		return nil, err
	}

	// Prepare the response
	response := map[string]interface{}{
		"resources":             resources,
		"total_found":           len(resources),
		"total_number_of_data": totalCount,
		"current_page":          (offset/itemsPerPage) + 1,
		"per_page":              itemsPerPage,
	}

	return response, nil
}

func (rr *ResourceRepository) CreateViewOnResource(resourceID int) (bool, error) {
	var resource models.Resource

	// Fetch the resource by ID
	err := rr.db.First(&resource, resourceID).Error
	if err != nil {
		return false, err
	}

	// Increment the views count
	resource.Views++

	// Save the updated resource
	err = rr.db.Save(&resource).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

// getFieldValue dynamically accesses the value of a field in the Resource struct
func getFieldValue(resource *models.Resource, fieldName string) string {
	switch fieldName {
	case "resource_category_option_id":
		return resource.ResourceCategoryOptionId
	case "resource_type_option_id":
		return resource.ResourceTypeOptionId
	case "sustainable_development_goals":
		return resource.SustainableDevelopmentGoals
	case "racial_equality_consciousness":
		return resource.RacialEqualityConsciousness
	case "resource_identity_group_id":
		return resource.ResourceIdentityGroupId
	case "resource_title":
		return resource.ResourceTitle
	case "resource_description":
		return resource.ResourceDescription
	default:
		return ""
	}
}

// setFieldValue dynamically sets the value of a field in the Resource struct
func setFieldValue(resource *models.Resource, fieldName string, value interface{}) {
	switch fieldName {
	case "resource_category_option_id":
		// Join []string into a comma-separated string and assign
		resource.ResourceCategoryOptionId = strings.Join(value.([]string), ",")
	case "resource_type_option_id":
		resource.ResourceTypeOptionId = strings.Join(value.([]string), ",")
	case "sustainable_development_goals":
		resource.SustainableDevelopmentGoals = strings.Join(value.([]string), ",")
	case "racial_equality_consciousness":
		resource.RacialEqualityConsciousness = strings.Join(value.([]string), ",")
	case "resource_identity_group_id":
		resource.ResourceIdentityGroupId = strings.Join(value.([]string), ",")
	case "resource_title":
		resource.ResourceTitle = value.(string)
	case "resource_description":
		resource.ResourceDescription = value.(string)
	}
}

func (repo *ResourceRepository) CreateLikeOnResource(resourceID int) (bool, error) {
	// Use GORM's Update feature
	result := repo.db.Model(&models.Resource{}).Where("id = ?", resourceID).Update("likes", gorm.Expr("likes + 1"))
	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (r *ResourceRepository) GetResourceByID(id int) (map[string]interface{}, error) {
    var resource models.Resource
    var comments []Comment

    // Fetch resource details
    err := r.db.
        Table("resources").
        Select(`resources.*, 
                users.id as user_id, 
                users.fullname, 
                users.email, 
                category_options.title as category_title`).
        Joins("INNER JOIN users ON resources.user_id = users.id").
        Joins("LEFT JOIN category_options ON resources.category_opt_id = category_options.id").
        Where("resources.id = ?", id).
        First(&resource).Error

    if err != nil {
        return map[string]interface{}{
            "status":    "success",
            "resources": nil,
            "comments":  []Comment{},
        }, nil
    }

    // Fetch all comments (including replies) for this resource
    err = r.db.
        Table("comments").
        Select("comments.*, users.fullname as commenter_name, users.email as commenter_email").
        Joins("INNER JOIN users ON comments.user_id = users.id").
        Where("comments.resource_id = ?", id).
        Order("comments.created_at ASC").
        Find(&comments).Error

    if err != nil {
        return nil, err
    }

    // Convert comments into a nested structure
    commentMap := make(map[int]*Comment)
    var rootComments []Comment

    for i := range comments {
        comment := &comments[i]
        comment.Replies = []Comment{} // Initialize Replies slice
        commentMap[comment.ID] = comment
    }

    for i := range comments {
        comment := &comments[i]
        if comment.ParentID != nil {
            parent, exists := commentMap[*comment.ParentID]
            if exists {
                parent.Replies = append(parent.Replies, *comment)
            }
        } else {
            rootComments = append(rootComments, *comment)
        }
    }

    // Update profile views
    r.db.Exec("UPDATE users SET views = views + 1 WHERE id = ?", resource.UserID)

    return map[string]interface{}{
        "status":    "success",
        "resources": resource,
        "comments":  rootComments, // Nested structure
    }, nil
}

// UpdateResourceStatus updates the status of a resource in the database
func (r *ResourceRepository) UpdateResourceStatus(id int, status string) error {
	// Update status where ID matches
	result := r.db.Model(&models.Resource{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no rows updated, resource may not exist")
	}

	return nil
}

// GetUserArticles retrieves all articles for a given user ID
func (r *ResourceRepository) GetUserArticles(userID int) ([]models.Resource, error) {
	var articles []models.Resource
	err := r.db.Where("user_id = ?", userID).Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (r *ResourceRepository) TotalUserArticles(userID int) (int, error) {
	var count int64
	err := r.db.Model(&models.Resource{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// func (r *ResourceRepository) GetRandomResources(offset, limit int) ([]models.Resource, int, error) {
// 	var resources []models.Resource
// 	var totalCount int64

// 	// Count total records
// 	if err := r.db.Model(&models.Resource{}).Count(&totalCount).Error; err != nil {
// 		return nil, 0, err
// 	}

// 	// Fetch resources with random order
// 	err := r.db.Preload("User").
// 		Order("RANDOM()").
// 		Offset(offset).
// 		Limit(limit).
// 		Find(&resources).Error
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	// Convert fields (matching PHP behavior)
// 	for i := range resources {
// 		// Convert CSV to array
// 		resources[i].ResourceCategoryOptionId = parseCSV(resources[i].Category)
// 		resources[i].SDGs = parseCSV(resources[i].SDGs)
// 		resources[i].RacialEquality = parseCSV(resources[i].RacialEquality)
// 		resources[i].IdentityGroup = parseCSV(resources[i].IdentityGroup)

// 		// Decode JSON fields
// 		resources[i].Title = parseJSON(resources[i].Title)
// 		resources[i].Description = parseJSON(resources[i].Description)
// 	}

// 	return resources, int(totalCount), nil
// }

// // Helper function to parse CSV fields
// func parseCSV(data string) []string {
// 	if data == "" {
// 		return []string{}
// 	}
// 	return strings.Split(data, ",")
// }

// // Helper function to parse JSON fields
// func parseJSON(data string) string {
// 	var result string
// 	err := json.Unmarshal([]byte(data), &result)
// 	if err != nil {
// 		return data // Return original if not valid JSON
// 	}
// 	return result
// }
