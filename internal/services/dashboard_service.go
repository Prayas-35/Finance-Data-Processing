package services

import (
	"context"

	"github.com/Prayas-35/Finance-Data-Processing/internal/models"
	"github.com/Prayas-35/Finance-Data-Processing/internal/repositories"
)

type DashboardService struct {
	dash *repositories.DashboardRepository
}

func NewDashboardService(dash *repositories.DashboardRepository) *DashboardService {
	return &DashboardService{dash: dash}
}

func (s *DashboardService) Summary(ctx context.Context, f repositories.DashboardFilter) (models.Summary, error) {
	return s.dash.Summary(ctx, f)
}

func (s *DashboardService) CategoryTotals(ctx context.Context, f repositories.DashboardFilter) ([]models.CategoryTotal, error) {
	return s.dash.CategoryTotals(ctx, f)
}

func (s *DashboardService) TrendsMonthly(ctx context.Context, f repositories.DashboardFilter) ([]models.TrendPoint, error) {
	return s.dash.TrendsMonthly(ctx, f)
}

func (s *DashboardService) RecentTransactions(ctx context.Context, f repositories.DashboardFilter, limit int) ([]models.FinancialRecord, error) {
	return s.dash.RecentTransactions(ctx, f, limit)
}
