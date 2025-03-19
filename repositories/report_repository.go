package repositories

import (
	"api-service/requests"
	"gorm.io/gorm"
	"api-service/models"
)

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (rr *ReportRepository) CreateReport(reportRequest *requests.ReportRequest) error {
	// Convert DTO to model
	report := models.Reports{
		ResourceId: reportRequest.ResourceID,
		Content:    reportRequest.Content,
		Fullname:   reportRequest.Fullname,
		Email:      reportRequest.Email,
	}

	// Insert into database
	return rr.db.Create(&report).Error
}

func (rr *ReportRepository) GetReportByID(reportID uint) (*models.Reports, error) {
	var report models.Reports
	if err := rr.db.First(&report, reportID).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

func (rr *ReportRepository) DeleteReportByID(reportID uint) error {
	// Check if the report exists before deleting
	var report models.Reports
	if err := rr.db.First(&report, reportID).Error; err != nil {
		return err 
	}

	// Delete the report
	return rr.db.Delete(&report).Error
}

func (rr *ReportRepository) GetAllReports() ([]models.Reports, error) {
	var reports []models.Reports
	err := rr.db.Find(&reports).Error
	return reports, err
}

func (rr *ReportRepository) UpdateReport(reportID int, updatedReport *models.Reports) (*models.Reports, error) {
	var report models.Reports
	if err := rr.db.First(&report, reportID).Error; err != nil {
		return nil, err
	}

	report.ResourceId = updatedReport.ResourceId
	report.Content = updatedReport.Content
	report.Fullname = updatedReport.Fullname
	report.Email = updatedReport.Email

	if err := rr.db.Save(&report).Error; err != nil {
		return nil, err
	}

	return &report, nil
}