package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/Prayas-35/Finance-Data-Processing/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.UserWithPassword, error) {
	query := `
SELECT id, name, email, role, is_active, password_hash, created_at, updated_at
FROM users
WHERE email = $1 AND is_active = TRUE
`

	var user models.UserWithPassword
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.IsActive,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]models.User, error) {
	query := `
SELECT id, name, email, role, is_active, created_at, updated_at
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		if err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (r *UserRepository) Create(ctx context.Context, name, email, passwordHash string, role models.Role) (string, error) {
	query := `
INSERT INTO users (name, email, password_hash, role)
VALUES ($1, $2, $3, $4)
RETURNING id
`

	var id string
	if err := r.pool.QueryRow(ctx, query, name, email, passwordHash, role).Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func (r *UserRepository) UpdateActive(ctx context.Context, id string, active bool) error {
	query := `
UPDATE users
SET is_active = $1, updated_at = NOW()
WHERE id = $2
`

	cmd, err := r.pool.Exec(ctx, query, active, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *UserRepository) Update(ctx context.Context, id, name string, role models.Role) error {
	query := `
UPDATE users
SET name = $1, role = $2, updated_at = NOW()
WHERE id = $3
`

	cmd, err := r.pool.Exec(ctx, query, name, role, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *UserRepository) SeedAdmin(ctx context.Context, email, name, passwordHash string) error {
	query := `
INSERT INTO users (name, email, password_hash, role)
VALUES ($1, $2, $3, 'admin')
ON CONFLICT (email) DO NOTHING
`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.pool.Exec(ctx, query, name, email, passwordHash)
	return err
}
