package services

import (
	"api-service/repositories"
	"errors"
	"api-service/models"
)

type CommentService struct {
	commentRepo repositories.CommentRepository
}

// NewCommentService initializes a new service
func NewCommentService(commentRepo repositories.CommentRepository) *CommentService {
	return &CommentService{commentRepo: commentRepo}
}

func (cs *CommentService) CreateComment(userID uint, pressID, parentID, content string) (uint, error) {
	comment := models.Comments{
		UserId:   userID,
		ParentId: parentID,
		Content:  content,
		PressId:  pressID,
	}

	return cs.commentRepo.Create(comment)
}

func (cs *CommentService) DeleteComment(userID uint, commentID string) error {
	comment, err := cs.commentRepo.FindByID(commentID)
	if err != nil {
		return errors.New("comment not found")
	}

	if comment.UserId != userID {
		return errors.New("unauthorized to delete this comment")
	}

	return cs.commentRepo.Delete(commentID)
}

func (cs *CommentService) GetCommentsByResource(resourceID string) ([]map[string]interface{}, error) {
	return cs.commentRepo.GetCommentsByResource(resourceID)
}

func (cs *CommentService) UpdateComment(commentID uint, updatedComment map[string]interface{}) (bool, error) {
	content, ok := updatedComment["content"].(string)
	if !ok {
		return false, errors.New("invalid content format")
	}

	success, err := cs.commentRepo.UpdateComment(commentID, content)
	if err != nil || !success {
		return false, errors.New("failed to update comment")
	}

	return true, nil
}

