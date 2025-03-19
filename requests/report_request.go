package requests

type ReportRequest struct {
	ResourceID uint `json:"resource_id"`
	Content     string `json:"content"`
	Fullname     string `json:"fullname"`
	Email     string `json:"email"`
}