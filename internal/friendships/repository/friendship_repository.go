package repository

import (
	"context"
	"fmt"

	"github.com/donca/user-crud/config/storage"
	"github.com/donca/user-crud/internal/friendships/interfaces"
	"github.com/donca/user-crud/internal/friendships/models"

	"github.com/jackc/pgx/v5"
)

type friendshipRepository struct {
	db *storage.PostgresPool
}

func NewFriendshipRepository(db *storage.PostgresPool) interfaces.FriendshipRepository {
	return &friendshipRepository{db: db}
}

func (r *friendshipRepository) CreateRequest(ctx context.Context, requesterID, addresseeID string) (*models.Friendship, error) {
	row := r.db.Pool.QueryRow(ctx, `
		INSERT INTO friendships (requester_id, addressee_id, status)
		VALUES ($1, $2, 'pending')
		RETURNING id, requester_id, addressee_id, status, created_at, updated_at`,
		requesterID, addresseeID)
	return scanFriendship(row)
}

func (r *friendshipRepository) FindByID(ctx context.Context, id string) (*models.Friendship, error) {
	row := r.db.Pool.QueryRow(ctx, `
		SELECT id, requester_id, addressee_id, status, created_at, updated_at
		FROM friendships WHERE id = $1`, id)
	f, err := scanFriendship(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return f, err
}

// FindBetween looks up the single friendship row for an unordered user pair, if any.
func (r *friendshipRepository) FindBetween(ctx context.Context, userA, userB string) (*models.Friendship, error) {
	row := r.db.Pool.QueryRow(ctx, `
		SELECT id, requester_id, addressee_id, status, created_at, updated_at
		FROM friendships
		WHERE (requester_id = $1 AND addressee_id = $2)
		   OR (requester_id = $2 AND addressee_id = $1)`,
		userA, userB)
	f, err := scanFriendship(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return f, err
}

func (r *friendshipRepository) UpdateStatus(ctx context.Context, id, addresseeID, status string) (*models.Friendship, error) {
	row := r.db.Pool.QueryRow(ctx, `
		UPDATE friendships SET status = $3, updated_at = NOW()
		WHERE id = $1 AND addressee_id = $2 AND status = 'pending'
		RETURNING id, requester_id, addressee_id, status, created_at, updated_at`,
		id, addresseeID, status)
	f, err := scanFriendship(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return f, err
}

func (r *friendshipRepository) ListPendingReceived(ctx context.Context, userID string) ([]models.Friendship, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, requester_id, addressee_id, status, created_at, updated_at
		FROM friendships WHERE addressee_id = $1 AND status = 'pending'
		ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("list pending received: %w", err)
	}
	return scanFriendships(rows)
}

func (r *friendshipRepository) ListPendingSent(ctx context.Context, userID string) ([]models.Friendship, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, requester_id, addressee_id, status, created_at, updated_at
		FROM friendships WHERE requester_id = $1 AND status = 'pending'
		ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("list pending sent: %w", err)
	}
	return scanFriendships(rows)
}

func (r *friendshipRepository) ListFriends(ctx context.Context, userID string) ([]models.FriendProfile, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT u.id, u.username, u.display_name, u.avatar_url
		FROM friendships f
		JOIN users u ON u.id = CASE
			WHEN f.requester_id = $1 THEN f.addressee_id
			ELSE f.requester_id
		END
		WHERE f.status = 'accepted' AND (f.requester_id = $1 OR f.addressee_id = $1)
		ORDER BY u.display_name`, userID)
	if err != nil {
		return nil, fmt.Errorf("list friends: %w", err)
	}
	defer rows.Close()
	var friends []models.FriendProfile
	for rows.Next() {
		var f models.FriendProfile
		if err := rows.Scan(&f.ID, &f.Username, &f.DisplayName, &f.AvatarURL); err != nil {
			return nil, fmt.Errorf("scan friend: %w", err)
		}
		friends = append(friends, f)
	}
	return friends, rows.Err()
}

// AreFriends reports whether an accepted friendship exists in either direction.
func (r *friendshipRepository) AreFriends(ctx context.Context, userA, userB string) (bool, error) {
	var exists bool
	err := r.db.Pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM friendships
			WHERE status = 'accepted'
			  AND (
			    (requester_id = $1 AND addressee_id = $2)
			    OR (requester_id = $2 AND addressee_id = $1)
			  )
		)`, userA, userB).Scan(&exists)
	return exists, err
}

func scanFriendship(row pgx.Row) (*models.Friendship, error) {
	var f models.Friendship
	if err := row.Scan(&f.ID, &f.RequesterID, &f.AddresseeID, &f.Status, &f.CreatedAt, &f.UpdatedAt); err != nil {
		return nil, err
	}
	return &f, nil
}

func scanFriendships(rows pgx.Rows) ([]models.Friendship, error) {
	defer rows.Close()
	var list []models.Friendship
	for rows.Next() {
		var f models.Friendship
		if err := rows.Scan(&f.ID, &f.RequesterID, &f.AddresseeID, &f.Status, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan friendship: %w", err)
		}
		list = append(list, f)
	}
	return list, rows.Err()
}
