package controllers

import (
	"api-service/requests"
	"api-service/services"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"api-service/exceptions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"api-service/security"
)

type ResourceController struct {
	resourceService *services.ResourceService
}

func NewResourceController(resourceService *services.ResourceService) *ResourceController {
	return &ResourceController{resourceService: resourceService}
}

// CreateResource handles resource creation
func (rc *ResourceController) CreateResource(c *gin.Context) {
	// Parse form data
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10MB limit
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	// Extract form fields
	form := c.Request.MultipartForm
	request := requests.CreateResourceRequest{
		ResourceType:          form.Value["resourceType[]"],
		StartDate:             getFormValue(form, "startDate"),
		EndDate:               getFormValue(form, "endDate"),
		ResourceContactTitle:  getFormValue(form, "resourceContactTitle"),
		ResourceOrganization:  getFormValue(form, "resourceOrganization"),
		ResourceTitle:         getFormValue(form, "resourceTitle"),
		ResourceDescription:   getFormValue(form, "resourceDescription"),
		ResourceCategory:      form.Value["resourceCategory[]"],
		IdentityGroup:         form.Value["identityGroup[]"],
		RacialSphere:          form.Value["racialSphere[]"],
		Sustainable:           form.Value["sustainable[]"],
		TargetAudience:        getFormValue(form, "targetAudience"),
		YearResource:          getFormValue(form, "yearResource"),
		Status:                getFormValue(form, "status"),
		Weblink:               getFormValue(form, "weblink"),
	}

	// Validate required fields
	if err := validateResourceRequest(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "errors": err})
		return
	}

	// Handle file upload
	files := form.File["uploadedFile"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Resource Attachment is required!"})
		return
	}

	uploadedFiles := []requests.UploadedFile{}
	for _, file := range files {
		// Validate file type
		fileExt := strings.ToLower(filepath.Ext(file.Filename))
		allowedExtensions := map[string]bool{".pdf": true, ".docx": true, ".csv": true}

		if !allowedExtensions[fileExt] {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid file type. Only PDF, DOCX, and CSV are allowed."})
			return
		}

		// Save file
		uploadPath := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to upload file."})
			return
		}

		uploadedFiles = append(uploadedFiles, requests.UploadedFile{
			FileName: file.Filename,
			FilePath: uploadPath,
			FileSize: file.Size,
		})
	}

	// Pass data to service layer
	resourceID, err := rc.resourceService.CreateResource(request, uploadedFiles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create resource."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Resource created successfully", "resourceID": resourceID})
}

// GetAllResources handles the GET request for retrieving paginated resources
func (rc *ResourceController) GetAllResources(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "10")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset value"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit value"})
		return
	}

	resources, err := rc.resourceService.GetAllResources(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resources"})
		return
	}

	c.JSON(http.StatusOK, resources)
}

func (c *ResourceController) GetResourceByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id") // Extract ID from URL path
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, `{"status":"error", "message":"Invalid resource ID"}`, http.StatusBadRequest)
		return
	}

	resource, err := c.resourceService.GetResourceByID(id)
	if err != nil {
		http.Error(w, `{"status":"error", "message":"Resource not found"}`, http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"status":    "success",
		"resource":  resource,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (rc *ResourceController) DeleteResource(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is DELETE
	if r.Method != http.MethodDelete {
		exceptions.SendErrorResponse(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract resource ID from URL parameters
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists || idStr == "" {
		exceptions.SendErrorResponse(w, "Resource ID is required", http.StatusBadRequest)
		return
	}

	// Convert ID to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		exceptions.SendErrorResponse(w, "Resource ID must be an integer", http.StatusBadRequest)
		return
	}

	// Call service to delete the resource
	deleted, err := rc.resourceService.DeleteResource(id)
	if err != nil {
		exceptions.SendErrorResponse(w, "Deletion failed or resource not found", http.StatusNotFound)
		return
	}

	if deleted {
		exceptions.SendSuccessResponse(w, "Resource deleted successfully", http.StatusOK)
	} else {
		exceptions.SendErrorResponse(w, "Deletion failed or resource not found", http.StatusNotFound)
	}
}

func (rc *ResourceController) Listed(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is GET
	if r.Method == http.MethodGet {
		// Get query parameters
		page := r.URL.Query().Get("page")
		if page == "" {
			page = "100"
		}

		// Convert page to integer and calculate offset
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			http.Error(w, "Invalid page number", http.StatusBadRequest)
			return
		}

		itemsPerPage := 10
		offset := (pageInt - 1) * itemsPerPage

		// Call the service to get resources
		response, err := rc.resourceService.GetLimitedResource(offset, itemsPerPage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send the response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (rc *ResourceController) CreateViews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		exceptions.SendErrorResponse(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract `id` from URL params
	vars := mux.Vars(r) 
	resourceID, err := strconv.Atoi(vars["id"]) 
	if err != nil || resourceID <= 0 {
		exceptions.SendErrorResponse(w, "Invalid resource ID", http.StatusBadRequest)
		return
	}

	// Call service to create a view on the resource
	success, err := rc.resourceService.CreateViewOnResource(resourceID)
	if err != nil || !success {
		exceptions.SendErrorResponse(w, "Unable to process your view right now. Please try again later.", http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Your view has been recorded."})
}

func (rc *ResourceController) CreateLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		exceptions.SendErrorResponse(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read JSON body
	var requestBody map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		exceptions.SendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Extract and validate `resource_id`
	resourceID, ok := requestBody["resource_id"].(float64)
	if !ok || int(resourceID) <= 0 {
		exceptions.SendErrorResponse(w, "Resource ID must be an integer", http.StatusBadRequest)
		return
	}

	// Call Service Layer
	success, err := rc.resourceService.CreateLikeOnResource(int(resourceID))
	if err != nil || !success {
		exceptions.SendErrorResponse(w, "Unable to process your like right now. Please try again later.", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Your like has been recorded."})
}

// UpdateStatus handles HTTP requests to update a resource's status
func (rc *ResourceController) UpdateStatus(w http.ResponseWriter, r *http.Request) {

	// Ensure the request method is PUT
	if r.Method != http.MethodPut {
		exceptions.SendErrorResponse(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract resource ID from URL
	urlParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(urlParts) == 0 || urlParts[0] == "" {
		exceptions.SendErrorResponse(w, "Resource ID required", http.StatusNotFound)
		return
	}

	// Convert ID to integer
	id, err := strconv.Atoi(urlParts[0])
	if err != nil {
		exceptions.SendErrorResponse(w, "Resource ID must be an integer", http.StatusBadRequest)
		return
	}

	// Parse JSON body
	var requestBody struct {
		Status string `json:"status"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		exceptions.SendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate status
	if requestBody.Status == "" {
		exceptions.SendErrorResponse(w, "Status is required", http.StatusBadRequest)
		return
	}

	// Call Service Layer
	err = rc.resourceService.UpdateResourceStatus(id, requestBody.Status)
	if err != nil {
		exceptions.SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Success response
	exceptions.SendSuccessResponse(w, "Resource Status Successfully Updated", http.StatusOK)
}

// UserArticles retrieves articles for the authenticated user
func (rc *ResourceController) UserArticles(c *gin.Context) {
	security := security.SecurityFilterChain{}
	user, err := security.IsValidToken(c)
	if err != nil || user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized access. Please login."})
		return
	}

	// Get user articles from service
	articles, err := rc.resourceService.GetUserArticles(int(user.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch user articles"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"data": articles})
}

func (rc *ResourceController) CourtUserArticles(c *gin.Context) {
	// Extract user from security context (token validation)
	security := security.SecurityFilterChain{}
	user, err := security.IsValidToken(c)
	if err != nil || user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized access. Please login."})
		return
	}

	// Fetch user's total articles
	data, err := rc.resourceService.TotalUserArticles(int(user.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch user articles."})
		return
	}
	// Return JSON response
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// Utility function to get a single form value
func getFormValue(form *multipart.Form, key string) string {
	values, exists := form.Value[key]
	if exists && len(values) > 0 {
		return values[0]
	}
	return ""
}

// Validate required fields
func validateResourceRequest(req requests.CreateResourceRequest) map[string]string {
	errors := make(map[string]string)

	requiredFields := map[string]string{
		"resourceType":          "Resource Type is required.",
		"startDate":             "Start Date is required.",
		"endDate":               "End Date is required.",
		"resourceContactTitle":  "Resource Contact Title is required.",
		"resourceOrganization":  "Resource Organization is required.",
		"resourceTitle":         "Resource Title is required.",
		"resourceDescription":   "Resource Description is required.",
		"resourceCategory":      "Resource Category is required.",
		"identityGroup":         "Identity Group is required.",
		"racialSphere":          "Racial Sphere is required.",
		"sustainable":           "Sustainable Development Goal is required.",
		"targetAudience":        "Target Audience is required.",
		"yearResource":          "Year of Resource is required.",
		"status":                "Status is required.",
		"weblink":               "Weblink is required.",
	}

	// Validate fields
	for field, errMsg := range requiredFields {
		if field == "resourceType" || field == "resourceCategory" || field == "identityGroup" || field == "racialSphere" || field == "sustainable" {
			if len(req.ResourceType) == 0 {
				errors[field] = errMsg
			}
		} else if getFieldValue(req, field) == "" {
			errors[field] = errMsg
		}
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// Utility function to get struct field value dynamically
func getFieldValue(req requests.CreateResourceRequest, field string) string {
	switch field {
	case "startDate":
		return req.StartDate
	case "endDate":
		return req.EndDate
	case "resourceContactTitle":
		return req.ResourceContactTitle
	case "resourceOrganization":
		return req.ResourceOrganization
	case "resourceTitle":
		return req.ResourceTitle
	case "resourceDescription":
		return req.ResourceDescription
	case "targetAudience":
		return req.TargetAudience
	case "yearResource":
		return req.YearResource
	case "status":
		return req.Status
	case "weblink":
		return req.Weblink
	default:
		return ""
	}
}

