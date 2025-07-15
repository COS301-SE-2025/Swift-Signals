package postgres

import (
	"context"
	"database/sql"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/db"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
	"github.com/lib/pq"
)

type PostgresUserRepo struct {
	db *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) db.UserRepository {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) CreateUser(ctx context.Context, u *model.User) (*model.User, error) {
	query := `INSERT INTO users (uuid, name, email, password, is_admin, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`

	_, err := r.db.ExecContext(ctx, query, u.ID, u.Name, u.Email, u.Password, u.IsAdmin)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				// Unique constraint violation
				return nil, errs.NewAlreadyExistsError("email already exists", map[string]any{"email": u.Email})

			case "23503":
				// Foreign key violation (e.g. non-existent user_id in user_intersections)
				return nil, errs.NewDatabaseError("invalid reference to related resource", err, nil)

			case "23502":
				// Not-null constraint violation
				return nil, errs.NewDatabaseError("missing required field", err, nil)

			default:
				return nil, errs.NewInternalError("postgres error", err, map[string]any{"postgresErrCode": pqErr.Code, "postgresErrMessage": pqErr.Message})
			}
		}

		return nil, errs.NewInternalError("query execution failed", err, nil)
	}

	return u, err
}

func (r *PostgresUserRepo) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	query := `SELECT id, name, email, password, is_admin, created_at, updated_at 
	          FROM users 
	          WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	user := &model.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Get intersection IDs
	intIDs, err := r.GetIntersectionsByUserID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Populate User IDs
	user.IntersectionIDs = intIDs

	return user, nil
}

func (r *PostgresUserRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, name, email, password, is_admin, created_at, updated_at 
	          FROM users 
	          WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)

	user := &model.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// No user found with that email
			return nil, nil
		}
		return nil, err
	}

	// Get intersection IDs
	intIDs, err := r.GetIntersectionsByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	// Populate User IDs
	user.IntersectionIDs = intIDs

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

func (r *PostgresUserRepo) DeleteUser(ctx context.Context, id string) error {
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
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Get intersection IDs
		intIDs, err := r.GetIntersectionsByUserID(ctx, user.ID)
		if err != nil {
			return nil, err
		}

		// Populate User IDs
		user.IntersectionIDs = intIDs

		users = append(users, user)
	}
	return users, nil
}

func (r *PostgresUserRepo) AddIntersectionID(ctx context.Context, userID string, intID string) error {
	query := `INSERT INTO user_intersections (user_id, intersection_id) 
	          VALUES ($1, $2) 
	          ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, userID, intID)
	return err
}

func (r *PostgresUserRepo) GetIntersectionsByUserID(ctx context.Context, userID string) ([]string, error) {
	query := `SELECT intersection_id 
	          FROM user_intersections 
	          WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
