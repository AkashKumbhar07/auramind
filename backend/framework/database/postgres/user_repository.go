package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/AkashKumbhar07/auramind/backend/framework/database/interfaces"
	auerrors "github.com/AkashKumbhar07/auramind/backend/framework/errors"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *interfaces.User) error {
	query := `INSERT INTO users (id, email, username, password, created_at, updated_at)
	           VALUES ($1, $2, $3, $4, $5, $6)`

	now := time.Now().Unix()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.Username, user.Password, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return auerrors.Wrap(auerrors.KindInternal, "create user", err)
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*interfaces.User, error) {
	query := `SELECT id, email, username, password, created_at, updated_at FROM users WHERE id = $1`

	user := &interfaces.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, auerrors.NotFound("user not found")
	}
	if err != nil {
		return nil, auerrors.Wrap(auerrors.KindInternal, "get user by id", err)
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*interfaces.User, error) {
	query := `SELECT id, email, username, password, created_at, updated_at FROM users WHERE email = $1`

	user := &interfaces.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, auerrors.NotFound("user not found")
	}
	if err != nil {
		return nil, auerrors.Wrap(auerrors.KindInternal, "get user by email", err)
	}
	return user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*interfaces.User, error) {
	query := `SELECT id, email, username, password, created_at, updated_at FROM users WHERE username = $1`

	user := &interfaces.User{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Email, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, auerrors.NotFound("user not found")
	}
	if err != nil {
		return nil, auerrors.Wrap(auerrors.KindInternal, "get user by username", err)
	}
	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *interfaces.User) error {
	query := `UPDATE users SET email = $1, username = $2, password = $3, updated_at = $4 WHERE id = $5`

	user.UpdatedAt = time.Now().Unix()

	result, err := r.db.ExecContext(ctx, query,
		user.Email, user.Username, user.Password, user.UpdatedAt, user.ID,
	)
	if err != nil {
		return auerrors.Wrap(auerrors.KindInternal, "update user", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return auerrors.NotFound("user not found")
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return auerrors.Wrap(auerrors.KindInternal, "delete user", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return auerrors.NotFound("user not found")
	}
	return nil
}

func (r *UserRepository) List(ctx context.Context, filter map[string]any) ([]*interfaces.User, error) {
	query := `SELECT id, email, username, password, created_at, updated_at FROM users`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, auerrors.Wrap(auerrors.KindInternal, "list users", err)
	}
	defer rows.Close()

	var users []*interfaces.User
	for rows.Next() {
		u := &interfaces.User{}
		if err := rows.Scan(&u.ID, &u.Email, &u.Username, &u.Password, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, auerrors.Wrap(auerrors.KindInternal, "scan user", err)
		}
		users = append(users, u)
	}
	return users, nil
}
