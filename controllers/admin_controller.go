package controllers

import (
	// "api-service/requests"
	"api-service/services"
	"net/http"
	// "net/http"
	// "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AdminController struct {
	adminService services.AdminService
}

func NewAdminController(adminService services.AdminService) *AdminController {
	return &AdminController{adminService: adminService}
}

func (ac *AdminController) GetPendingArticles(c *gin.Context) {
	// Fetch pending articles
	resources, err := ac.adminService.GetPendingArticles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch pending articles."})
		return
	}

	// Check if resources are found
	if len(resources) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "No resources found."})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": resources})
}