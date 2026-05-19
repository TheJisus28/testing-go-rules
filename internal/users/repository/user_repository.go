package repository

import (
	"context"
	"fmt"

	"github.com/donca/user-crud/config/storage"
	"github.com/donca/user-crud/internal/users/interfaces"
	"github.com/donca/user-crud/internal/users/models"

	"github.com/jackc/pgx/v5"
)

type userRepository struct {
	db *storage.PostgresPool
}

func NewUserRepository(db *storage.PostgresPool) interfaces.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, username, email, passwordHash, displayName string) (*models.User, error) {
	row := r.db.Pool.QueryRow(ctx, `
		INSERT INTO users (username, email, password_hash, display_name)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, email, display_name, bio, avatar_url, created_at, updated_at`,
		username, email, passwordHash, displayName)
	return scanUser(row)
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	row := r.db.Pool.QueryRow(ctx, `
		SELECT id, username, email, display_name, bio, avatar_url, created_at, updated_at
		FROM users WHERE id = $1`, id)
	u, err := scanUser(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	row := r.db.Pool.QueryRow(ctx, `
		SELECT id, username, email, display_name, bio, avatar_url, created_at, updated_at
		FROM users WHERE username = $1`, username)
	u, err := scanUser(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	row := r.db.Pool.QueryRow(ctx, `
		SELECT id, username, email, display_name, bio, avatar_url, created_at, updated_at
		FROM users WHERE email = $1`, email)
	u, err := scanUser(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (r *userRepository) FindPasswordHash(ctx context.Context, id string) (string, error) {
	var hash string
	err := r.db.Pool.QueryRow(ctx, `SELECT password_hash FROM users WHERE id = $1`, id).Scan(&hash)
	if err == pgx.ErrNoRows {
		return "", nil
	}
	return hash, err
}

func (r *userRepository) Update(ctx context.Context, id string, email, displayName, bio, avatarURL *string) (*models.User, error) {
	row := r.db.Pool.QueryRow(ctx, `
		UPDATE users SET
			email = COALESCE($2, email),
			display_name = COALESCE($3, display_name),
			bio = COALESCE($4, bio),
			avatar_url = COALESCE($5, avatar_url),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, username, email, display_name, bio, avatar_url, created_at, updated_at`,
		id, email, displayName, bio, avatarURL)
	u, err := scanUser(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (r *userRepository) UpdateProfile(ctx context.Context, id string, displayName, bio, avatarURL *string) (*models.User, error) {
	return r.Update(ctx, id, nil, displayName, bio, avatarURL)
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.db.Pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]models.UserProfile, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, username, display_name, bio, avatar_url, created_at
		FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()
	var users []models.UserProfile
	for rows.Next() {
		var u models.UserProfile
		if err := rows.Scan(&u.ID, &u.Username, &u.DisplayName, &u.Bio, &u.AvatarURL, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan user profile: %w", err)
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func scanUser(row pgx.Row) (*models.User, error) {
	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.DisplayName, &u.Bio, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("scan user: %w", err)
	}
	return &u, nil
}
