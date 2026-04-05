package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Prayas-35/Finance-Data-Processing/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardFilter struct {
	UserID   string
	Role     string
	FromDate *time.Time
	ToDate   *time.Time
}

type DashboardRepository struct {
	pool *pgxpool.Pool
}

func NewDashboardRepository(pool *pgxpool.Pool) *DashboardRepository {
	return &DashboardRepository{pool: pool}
}

func (r *DashboardRepository) Summary(ctx context.Context, f DashboardFilter) (models.Summary, error) {
	where, args := dashboardWhereClause(f)
	query := fmt.Sprintf(`
		SELECT
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0)::text,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0)::text,
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE -amount END), 0)::text
		FROM financial_records
		WHERE %s
	`, where)

	var s models.Summary
	err := r.pool.QueryRow(ctx, query, args...).Scan(&s.TotalIncome, &s.TotalExpenses, &s.NetBalance)
	return s, err
}

func (r *DashboardRepository) CategoryTotals(ctx context.Context, f DashboardFilter) ([]models.CategoryTotal, error) {
	where, args := dashboardWhereClause(f)
	query := fmt.Sprintf(`
		SELECT category, COALESCE(SUM(amount), 0)::text
		FROM financial_records
		WHERE %s
		GROUP BY category
		ORDER BY category ASC
	`, where)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]models.CategoryTotal, 0)
	for rows.Next() {
		var c models.CategoryTotal
		if err = rows.Scan(&c.Category, &c.Total); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, rows.Err()
}

func (r *DashboardRepository) TrendsMonthly(ctx context.Context, f DashboardFilter) ([]models.TrendPoint, error) {
	where, args := dashboardWhereClause(f)
	query := fmt.Sprintf(`
		SELECT
			TO_CHAR(DATE_TRUNC('month', date), 'YYYY-MM') AS period,
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0)::text AS income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0)::text AS expense
		FROM financial_records
		WHERE %s
		GROUP BY DATE_TRUNC('month', date)
		ORDER BY DATE_TRUNC('month', date) ASC
	`, where)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]models.TrendPoint, 0)
	for rows.Next() {
		var t models.TrendPoint
		if err = rows.Scan(&t.Period, &t.Income, &t.Expense); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, rows.Err()
}

func (r *DashboardRepository) RecentTransactions(ctx context.Context, f DashboardFilter, limit int) ([]models.FinancialRecord, error) {
	where, args := dashboardWhereClause(f)
	args = append(args, limit)
	limitPos := len(args)
	query := fmt.Sprintf(`
		SELECT id, user_id, amount::text, type, category, date, COALESCE(notes, ''), created_at, updated_at
		FROM financial_records
		WHERE %s
		ORDER BY date DESC, created_at DESC
		LIMIT $%d
	`, where, limitPos)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]models.FinancialRecord, 0)
	for rows.Next() {
		var rec models.FinancialRecord
		if err = rows.Scan(&rec.ID, &rec.UserID, &rec.Amount, &rec.Type, &rec.Category, &rec.Date, &rec.Notes, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, rec)
	}
	return result, rows.Err()
}

func dashboardWhereClause(f DashboardFilter) (string, []interface{}) {
	args := make([]interface{}, 0)
	parts := []string{"deleted_at IS NULL"}

	if f.Role != "admin" {
		args = append(args, f.UserID)
		parts = append(parts, fmt.Sprintf("user_id = $%d", len(args)))
	}
	if f.FromDate != nil {
		args = append(args, *f.FromDate)
		parts = append(parts, fmt.Sprintf("date >= $%d", len(args)))
	}
	if f.ToDate != nil {
		args = append(args, *f.ToDate)
		parts = append(parts, fmt.Sprintf("date <= $%d", len(args)))
	}

	return strings.Join(parts, " AND "), args
}
