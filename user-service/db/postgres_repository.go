package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/COS301-SE-2025/Swift-Signals/user-service/models"
	_ "github.com/lib/pq"
)

type PostgresUserRepo struct {
	db *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) UserRepository {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) CreateUser(ctx context.Context, u *models.User) (*models.User, error) {
	query := `INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, u.ID, u.Name, u.Email, u.Password)
	return u, err
}

func (r *PostgresUserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	query := `SELECT id, name, email, password FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Name, &u.Email, &u.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u, nil
}
