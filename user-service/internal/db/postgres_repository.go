package db

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
	_ "github.com/lib/pq"
)

type PostgresUserRepo struct {
	db *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) UserRepository {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) CreateUser(ctx context.Context, u *model.User) (*model.User, error) {
	query := `INSERT INTO users (name, email, password, is_admin, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, NOW(), NOW())`
	_, err := r.db.ExecContext(ctx, query, u.Name, u.Email, u.Password, u.IsAdmin)
	return u, err
}

func (r *PostgresUserRepo) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	query := `SELECT id, name, email, password, is_admin, created_at, updated_at 
	          FROM users 
	          WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	user := &model.User{}
	err := row.Scan(&id, &user.Name, &user.Email, &user.Password, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	user.ID = strconv.Itoa(id)

	// Get intersection IDs
	intIDs, err := r.GetIntersectionsByUserID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Populate User IDs
	user.IntersectionIDs = make([]int32, len(intIDs))
	for i, v := range intIDs {
		user.IntersectionIDs[i] = int32(v)
	}

	return user, nil
}

func (r *PostgresUserRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, name, email, password, is_admin, created_at, updated_at 
	          FROM users 
	          WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)

	user := &model.User{}
	var id int
	err := row.Scan(&id, &user.Name, &user.Email, &user.Password, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// No user found with that email
			return nil, nil
		}
		return nil, err
	}
	user.ID = strconv.Itoa(id)

	// Get intersection IDs
	intIDs, err := r.GetIntersectionsByUserID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Populate User IDs
	user.IntersectionIDs = make([]int32, len(intIDs))
	for i, v := range intIDs {
		user.IntersectionIDs[i] = int32(v)
	}

	return user, nil
}

func (r *PostgresUserRepo) UpdateUser(ctx context.Context, u *model.User) (*model.User, error) {
	query := `UPDATE users
	          SET name = $1, email = $2, password = $3, is_admin = $4, updated_at = NOW()
	          WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, u.Name, u.Email, u.Password, u.IsAdmin, u.ID)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *PostgresUserRepo) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users 
	          WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresUserRepo) ListUsers(ctx context.Context, limit, offset int) ([]*model.User, error) {
	query := `SELECT id, name, email, password, is_admin, created_at, updated_at 
	          FROM users 
	          ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		var id int
		err := rows.Scan(&id, &user.Name, &user.Email, &user.Password, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		user.ID = strconv.Itoa(id)

		// Get intersection IDs
		intIDs, err := r.GetIntersectionsByUserID(ctx, id)
		if err != nil {
			return nil, err
		}

		// Populate User IDs
		user.IntersectionIDs = make([]int32, len(intIDs))
		for i, v := range intIDs {
			user.IntersectionIDs[i] = int32(v)
		}

		users = append(users, user)
	}
	return users, nil
}

func (r *PostgresUserRepo) AddIntersectionID(ctx context.Context, userID int, intID int) error {
	query := `INSERT INTO user_intersections (user_id, intersection_id) 
	          VALUES ($1, $2) 
	          ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, userID, intID)
	return err
}

func (r *PostgresUserRepo) GetIntersectionsByUserID(ctx context.Context, userID int) ([]int, error) {
	query := `SELECT intersection_id 
	          FROM user_intersections 
	          WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
