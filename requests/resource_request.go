package requests

// CreateResourceRequest represents the request payload
type CreateResourceRequest struct {
	ResourceType          []string `json:"resourceType"`
	StartDate             string   `json:"startDate"`
	EndDate               string   `json:"endDate"`
	ResourceContactTitle  string   `json:"resourceContactTitle"`
	ResourceOrganization  string   `json:"resourceOrganization"`
	ResourceTitle         string   `json:"resourceTitle"`
	ResourceDescription   string   `json:"resourceDescription"`
	ResourceCategory      []string `json:"resourceCategory"`
	IdentityGroup         []string `json:"identityGroup"`
	RacialSphere          []string `json:"racialSphere"`
	Sustainable           []string `json:"sustainable"`
	TargetAudience        string   `json:"targetAudience"`
	YearResource          string   `json:"yearResource"`
	Status                string   `json:"status"`
	Weblink               string   `json:"weblink"`
}

// UploadedFile represents an uploaded file
type UploadedFile struct {
	FileName string
	FilePath string
	FileSize int64
}
