package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Prayas-35/Finance-Data-Processing/internal/models"
	"github.com/Prayas-35/Finance-Data-Processing/internal/repositories"
)

type RecordService struct {
	records *repositories.RecordRepository
}

func NewRecordService(records *repositories.RecordRepository) *RecordService {
	return &RecordService{records: records}
}

func (s *RecordService) Create(ctx context.Context, rec models.FinancialRecord) (string, error) {
	if rec.UserID == "" || rec.Amount == "" || rec.Category == "" || rec.Type == "" {
		return "", fmt.Errorf("user_id, amount, category and type are required")
	}
	if rec.Date.IsZero() {
		rec.Date = time.Now().UTC()
	}
	if rec.Type != models.RecordTypeIncome && rec.Type != models.RecordTypeExpense {
		return "", fmt.Errorf("invalid record type")
	}
	return s.records.Create(ctx, rec)
}

func (s *RecordService) List(ctx context.Context, filter repositories.RecordFilter) ([]models.FinancialRecord, error) {
	return s.records.List(ctx, filter)
}

func (s *RecordService) GetByID(ctx context.Context, id, requesterID, role string) (*models.FinancialRecord, error) {
	return s.records.GetByID(ctx, id, requesterID, role)
}

func (s *RecordService) Update(ctx context.Context, id, requesterID, role string, rec models.FinancialRecord) error {
	if rec.Type != models.RecordTypeIncome && rec.Type != models.RecordTypeExpense {
		return fmt.Errorf("invalid record type")
	}
	return s.records.Update(ctx, id, requesterID, role, rec)
}

func (s *RecordService) SoftDelete(ctx context.Context, id, requesterID, role string) error {
	return s.records.SoftDelete(ctx, id, requesterID, role)
}
