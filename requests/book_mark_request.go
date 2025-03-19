package requests

type BookmarkRequest struct {
	ResourceID uint `json:"resource_id"`
	UserID     string `json:"userId"`
}