package services

import (
	"api-service/repositories"
	"api-service/requests"
	"api-service/models"
)

type ReportService struct {
	reportRepo repositories.ReportRepository
}

func NewReportService(reportRepo repositories.ReportRepository) *ReportService {
	return &ReportService{reportRepo: reportRepo}
}

func (rs *ReportService) CreateReport(report *requests.ReportRequest) error {
	return rs.reportRepo.CreateReport(report)
}

func (rs *ReportService) GetReportByID(reportID uint) (*models.Reports, error) {
	return rs.reportRepo.GetReportByID(reportID)
}

func (rs *ReportService) DeleteReportByID(reportID uint) error {
	return rs.reportRepo.DeleteReportByID(reportID)
}

func (rs *ReportService) GetAllReports() ([]models.Reports, error) {
	return rs.reportRepo.GetAllReports()
}

func (rs *ReportService) UpdateReport(reportID int, updatedReport *models.Reports) (*models.Reports, error) {
	return rs.reportRepo.UpdateReport(reportID, updatedReport)
}