package repositories

import (
	"api-service/models"
	"log"
	"strconv"
	"gorm.io/gorm"
)

// BookmarkRepository handles DB interactions
type BookmarkRepository struct {
	db *gorm.DB
}

// NewBookmarkRepository initializes a new repository
func NewBookmarkRepository(db *gorm.DB) *BookmarkRepository {
	return &BookmarkRepository{db: db}
}

// Create inserts a bookmark into the database using GORM
func (br *BookmarkRepository) Create(userID, resourceID string) bool {
	// Convert userID from string to uint
	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		log.Println("Invalid user ID format:", err)
		return false
	}

	bookmark := models.Bookmark{
		UserId:     uint(userIDUint),
		ResourceId: resourceID,
	}

	// Save to DB using GORM
	if err := br.db.Create(&bookmark).Error; err != nil {
		log.Println("Error saving bookmark:", err)
		return false
	}
	return true
}

func (br *BookmarkRepository) FindByID(id uint) (*models.Bookmark, error) {
	var bookmark models.Bookmark
	result := br.db.First(&bookmark, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &bookmark, nil
}

func (br *BookmarkRepository) ExistsByID(id uint) (bool, error) {
	var count int64
	err := br.db.Model(&models.Bookmark{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (br *BookmarkRepository) DeleteByID(id uint) error {
	result := br.db.Delete(&models.Bookmark{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// FetchAll retrieves all bookmarks from the database
func (br *BookmarkRepository) FetchAll() ([]models.Bookmark, error) {
	var bookmarks []models.Bookmark
	err := br.db.Find(&bookmarks).Error
	if err != nil {
		return nil, err
	}
	return bookmarks, nil
}

// GetUserTotalBookmarks returns the total bookmark count for a given user
func (br *BookmarkRepository) GetUserTotalBookmarks(userID uint) (int64, error) {
	var total int64
	err := br.db.Raw(`
		SELECT COUNT(b.resource_id) 
		FROM resources r 
		INNER JOIN bookmark b ON r.id = b.resource_id 
		WHERE r.user_id = ?
	`, userID).Scan(&total).Error

	if err != nil {
		return 0, err
	}
	return total, nil
}