package repositories

import (
	"api-service/models"
	"errors"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (cr *CommentRepository) Create(comment models.Comments) (uint, error) {
	result := cr.db.Create(&comment)
	if result.Error != nil {
		return 0, result.Error
	}
	return comment.ID, nil
}

func (cr *CommentRepository) FindByID(commentID string) (*models.Comments, error) {
	var comment models.Comments
	if err := cr.db.First(&comment, "id = ?", commentID).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (cr *CommentRepository) Delete(commentID string) error {
	return cr.db.Delete(&models.Comments{}, "id = ?", commentID).Error
}

func (cr *CommentRepository) GetCommentsByResource(resourceID string) ([]map[string]interface{}, error) {
	var comments []models.Comments

	err := cr.db.Where("press_id = ?", resourceID).Order("created_at ASC").Find(&comments).Error
	if err != nil {
		return nil, err
	}

	// Create a map to store comments by ID
	commentMap := make(map[uint]map[string]interface{})

	for _, comment := range comments {
		// Convert comment struct to map and add an empty replies array
		commentMap[comment.ID] = map[string]interface{}{
			"id":        comment.ID,
			"user_id":   comment.UserId,
			"parent_id": comment.ParentId,
			"content":   comment.Content,
			"press_id":  comment.PressId,
			"created_at": comment.CreatedAt,
			"replies":   []map[string]interface{}{}, 
		}
	}

	// Build the nested structure
	var rootComments []map[string]interface{}
	for _, comment := range commentMap {
		parentID, ok := comment["parent_id"].(*uint)
		if ok && parentID != nil {
			if parentComment, exists := commentMap[*parentID]; exists {
				parentComment["replies"] = append(parentComment["replies"].([]map[string]interface{}), comment)
			}
		} else {
			rootComments = append(rootComments, comment)
		}
	}

	return rootComments, nil
}

// Get comment by ID
func (cr *CommentRepository) GetCommentByID(commentID uint) (*models.Comments, error) {
	var comment models.Comments
	if err := cr.db.First(&comment, commentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("comment not found")
		}
		return nil, err
	}
	return &comment, nil
}

// Update comment by ID
func (cr *CommentRepository) UpdateComment(commentID uint, newContent string) (bool, error) {
	result := cr.db.Model(&models.Comments{}).Where("id = ?", commentID).Update("content", newContent)
	if result.Error != nil {
		return false, errors.New("failed to update comment")
	}

	if result.RowsAffected == 0 {
		return false, errors.New("no changes made")
	}

	return true, nil
}


