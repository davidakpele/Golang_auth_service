package controllers

import (
	"api-service/requests"
	"api-service/services"
	"net/http"
	"strconv"
	"api-service/models"
	"github.com/gin-gonic/gin"
)

type ReportController struct {
	reportService services.ReportService
}

func NewReportController(reportService services.ReportService) *ReportController {
	return &ReportController{reportService: reportService}
}

func (rc *ReportController) CreateReport(c *gin.Context) {
	var report requests.ReportRequest

	// Bind JSON request body
	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request body"})
		return
	}

	// Validate required fields
	if report.ResourceID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Resource ID is required"})
		return
	}
	if report.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Content is required"})
		return
	}
	if report.Fullname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Fullname is required"})
		return
	}
	if report.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Email is required"})
		return
	}

	// Call service to create report
	if err := rc.reportService.CreateReport(&report); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create report"})
		return
	}

	// Success response
	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Report has been submitted, we’ll carefully look into it. Thanks!"})
}

func (rc *ReportController) GetReportByID(c *gin.Context) {
	// Get ID parameter from URL
	reportID, err := strconv.Atoi(c.Param("id"))
	if err != nil || reportID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid report ID"})
		return
	}

	// Fetch report from service
	report, fetchErr := rc.reportService.GetReportByID(uint(reportID))
	if fetchErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Report not found"})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": report})
}

func (rc *ReportController) DeleteReportByID(c *gin.Context) {
	// Get ID parameter from URL
	reportID, err := strconv.Atoi(c.Param("id"))
	if err != nil || reportID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid report ID"})
		return
	}

	// Delete report via service
	err = rc.reportService.DeleteReportByID(uint(reportID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Report not found or could not be deleted"})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Report deleted successfully"})
}

func (rc *ReportController) GetAllReports(c *gin.Context) {
	reports, err := rc.reportService.GetAllReports()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch reports"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": reports})
}

func (rc *ReportController) UpdateReport(c *gin.Context) {
	reportID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Report ID must be an integer"})
		return
	}

	var report models.Reports
	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request data"})
		return
	}

	updatedReport, err := rc.reportService.UpdateReport(reportID, &report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to update report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedReport})
}