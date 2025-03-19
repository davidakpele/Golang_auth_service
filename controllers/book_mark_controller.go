package controllers

import (
	"api-service/requests"
	"api-service/services"
	"api-service/security"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

// BookmarkController handles bookmark-related requests
type BookmarkController struct {
	bookmarkService services.BookmarkService
}

// NewBookmarkController initializes a new controller
func NewBookmarkController(bookmarkService services.BookmarkService) *BookmarkController {
	return &BookmarkController{bookmarkService: bookmarkService}
}

// Create handles bookmark creation
func (bc *BookmarkController) CreateBookmark(c *gin.Context) {
	var request requests.BookmarkRequest

	// Parse JSON request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid JSON format"})
		return
	}

	// Retrieve user ID from session
	sessionUserID, exists := c.Get("user_id")
	if !exists || sessionUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User is not authenticated"})
		return
	}

	// Convert ResourceID (uint) to string
	resourceIDStr := strconv.FormatUint(uint64(request.ResourceID), 10)

	// Call service to create the bookmark
	err := bc.bookmarkService.CreateBookmark(request.UserID, resourceIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Bookmark has been saved successfully."})
}

func (bc *BookmarkController) GetBookmarkByID(c *gin.Context) {
	// Get ID from URL
	idParam := c.Param("id")

	// Convert ID to uint
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Bookmark ID must be an integer"})
		return
	}

	// Call service
	bookmark, err := bc.bookmarkService.GetBookmarkByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Bookmark not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": bookmark})
}

func (bc *BookmarkController) DeleteBookmarkByID(c *gin.Context) {
	// Get ID from URL
	idParam := c.Param("id")

	// Convert ID to uint
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Bookmark ID must be an integer"})
		return
	}

	// Call service
	err = bc.bookmarkService.DeleteBookmarkByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Bookmark not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Bookmark deleted successfully"})
}

// GetAllBookmarks handles fetching all bookmarks
func (bc *BookmarkController) GetAllBookmarks(c *gin.Context) {
	bookmarks, err := bc.bookmarkService.GetAllBookmarks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch bookmarks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": bookmarks})
}

// GetTotalUserBookmarks handles fetching the total number of bookmarks for a user
func (bc *BookmarkController) GetTotalUserBookmarks(c *gin.Context) {
	// Decode user ID from JWT
	security := security.SecurityFilterChain{}
	user, err := security.IsValidToken(c)
	if err != nil || user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized access. Please login."})
		return
	}

	// Fetch total bookmarks count for the user
	total, err := bc.bookmarkService.GetTotalUserBookmarks(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch bookmark count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": total})
}