package controllers

import (
	"api-service/security"
	"api-service/services"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"

)

type CommentController struct {
	commentService services.CommentService
}

func NewCommentController(commentService services.CommentService) *CommentController {
	return &CommentController{commentService: commentService}
}

func (cc *CommentController) CreateComment(c *gin.Context) {
	security := security.SecurityFilterChain{}
  	user, err := security.IsValidToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized access. Please Login."})
		return
	}

	var request struct {
		ResourceID string `json:"resource_id"`
		ParentID   string `json:"parent_id,omitempty"`
		Content    string `json:"content"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request format."})
		return
	}

	if request.ResourceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Resource ID is required."})
		return
	}
	if request.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Comment content is required."})
		return
	}

	commentID, err := cc.commentService.CreateComment(user.ID, request.ResourceID, request.ParentID, request.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to post comment."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Comment posted.", "comment_id": commentID})
}

func (cc *CommentController) DeleteComment(c *gin.Context) {
	commentID := c.Param("id")
	security := security.SecurityFilterChain{}
  	user, err := security.IsValidToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized access. Please Login."})
		return
	}

	err = cc.commentService.DeleteComment(user.ID, commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to delete comment."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Comment deleted."})
}

func (cc *CommentController) GetCommentsByResource(c *gin.Context) {
	resourceID := c.Param("resourceId")

	comments, err := cc.commentService.GetCommentsByResource(resourceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "comments": comments})
}

func (cc *CommentController) UpdateComment(c *gin.Context) {
	// Get ID parameter from URL
	commentID, err := strconv.Atoi(c.Param("id"))
	if err != nil || commentID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Comment ID must be a valid integer",
		})
		return
	}

	// Parse request body
	var updatedComment map[string]interface{}
	if err := c.ShouldBindJSON(&updatedComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request body",
		})
		return
	}

	// Call service to update comment 
	success, updateErr := cc.commentService.UpdateComment(uint(commentID), updatedComment)
	if updateErr != nil || !success {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to update comment",
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Comment has been updated successfully!",
	})
}

