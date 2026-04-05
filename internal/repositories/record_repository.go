package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Prayas-35/Finance-Data-Processing/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RecordFilter struct {
	UserID   string
	Role     string
	FromDate *time.Time
	ToDate   *time.Time
	Category string
	Type     string
	Limit    int
	Offset   int
}

type RecordRepository struct {
	pool *pgxpool.Pool
}

func NewRecordRepository(pool *pgxpool.Pool) *RecordRepository {
	return &RecordRepository{pool: pool}
}

func (r *RecordRepository) Create(ctx context.Context, rec models.FinancialRecord) (string, error) {
	query := `
INSERT INTO financial_records (user_id, amount, type, category, date, notes)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id
`

	var id string
	err := r.pool.QueryRow(ctx, query, rec.UserID, rec.Amount, rec.Type, rec.Category, rec.Date, rec.Notes).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *RecordRepository) List(ctx context.Context, f RecordFilter) ([]models.FinancialRecord, error) {
	args := []interface{}{}
	where := []string{"deleted_at IS NULL"}

	if f.Role != "admin" {
		args = append(args, f.UserID)
		where = append(where, fmt.Sprintf("user_id = $%d", len(args)))
	}
	if f.FromDate != nil {
		args = append(args, *f.FromDate)
		where = append(where, fmt.Sprintf("date >= $%d", len(args)))
	}
	if f.ToDate != nil {
		args = append(args, *f.ToDate)
		where = append(where, fmt.Sprintf("date <= $%d", len(args)))
	}
	if f.Category != "" {
		args = append(args, f.Category)
		where = append(where, fmt.Sprintf("category = $%d", len(args)))
	}
	if f.Type != "" {
		args = append(args, f.Type)
		where = append(where, fmt.Sprintf("type = $%d", len(args)))
	}

	args = append(args, f.Limit)
	limitPos := len(args)
	args = append(args, f.Offset)
	offsetPos := len(args)

	query := fmt.Sprintf(`
SELECT id, user_id, amount::text, type, category, date, COALESCE(notes, ''), created_at, updated_at
FROM financial_records
WHERE %s
ORDER BY date DESC, created_at DESC
LIMIT $%d OFFSET $%d
`, strings.Join(where, " AND "), limitPos, offsetPos)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make([]models.FinancialRecord, 0)
	for rows.Next() {
		var rec models.FinancialRecord
		if err = rows.Scan(&rec.ID, &rec.UserID, &rec.Amount, &rec.Type, &rec.Category, &rec.Date, &rec.Notes, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	return records, rows.Err()
}

func (r *RecordRepository) GetByID(ctx context.Context, id, requesterID, role string) (*models.FinancialRecord, error) {
	query := `
SELECT id, user_id, amount::text, type, category, date, COALESCE(notes, ''), created_at, updated_at
FROM financial_records
WHERE id = $1 AND deleted_at IS NULL
`

	var rec models.FinancialRecord
	err := r.pool.QueryRow(ctx, query, id).Scan(&rec.ID, &rec.UserID, &rec.Amount, &rec.Type, &rec.Category, &rec.Date, &rec.Notes, &rec.CreatedAt, &rec.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if role != "admin" && rec.UserID != requesterID {
		return nil, pgx.ErrNoRows
	}
	return &rec, nil
}

func (r *RecordRepository) Update(ctx context.Context, id, requesterID, role string, rec models.FinancialRecord) error {
	query := `
UPDATE financial_records
SET amount = $1, type = $2, category = $3, date = $4, notes = $5, updated_at = NOW()
WHERE id = $6 AND deleted_at IS NULL
`
	args := []interface{}{rec.Amount, rec.Type, rec.Category, rec.Date, rec.Notes, id}

	if role != "admin" {
		query += " AND user_id = $7"
		args = append(args, requesterID)
	}

	cmd, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("record not found")
	}
	return nil
}

func (r *RecordRepository) SoftDelete(ctx context.Context, id, requesterID, role string) error {
	query := `
UPDATE financial_records
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
`
	args := []interface{}{id}
	if role != "admin" {
		query += " AND user_id = $2"
		args = append(args, requesterID)
	}

	cmd, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("record not found")
	}
	return nil
}
