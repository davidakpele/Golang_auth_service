package services

import (
	"api-service/repositories"
	"errors"
	"api-service/models"
)

type BookmarkService struct {
	bookmarkRepo repositories.BookmarkRepository
}

// NewBookmarkService initializes a new service
func NewBookmarkService(bookmarkRepo repositories.BookmarkRepository) *BookmarkService {
	return &BookmarkService{bookmarkRepo: bookmarkRepo}
}

// CreateBookmark handles bookmark creation logic
func (bs *BookmarkService) CreateBookmark(userID, resourceID string) error {
	if resourceID == "" {
		return ErrInvalidResourceID
	}

	// Call repository to store the bookmark
	success := bs.bookmarkRepo.Create(userID, resourceID)
	if !success {
		return ErrFailedToSave
	}

	return nil
}

func (bs *BookmarkService) GetBookmarkByID(id uint) (*models.Bookmark, error) {
	bookmark, err := bs.bookmarkRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("bookmark not found")
	}
	return bookmark, nil
}

func (bs *BookmarkService) DeleteBookmarkByID(id uint) error {
	exists, err := bs.bookmarkRepo.ExistsByID(id)
	if err != nil {
		return errors.New("failed to check bookmark existence")
	}
	if !exists {
		return errors.New("bookmark not found")
	}

	return bs.bookmarkRepo.DeleteByID(id)
}

// GetAllBookmarks fetches all bookmarks from the database
func (bs *BookmarkService) GetAllBookmarks() ([]models.Bookmark, error) {
	return bs.bookmarkRepo.FetchAll()
}

func (bs *BookmarkService) GetTotalUserBookmarks(userID uint) (int64, error) {
	return bs.bookmarkRepo.GetUserTotalBookmarks(userID)
}

// Custom errors
var (
	ErrInvalidResourceID = errors.New("resource ID is required")
	ErrFailedToSave      = errors.New("failed to save bookmark")
)