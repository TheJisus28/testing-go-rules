package repository

import (
	"context"
	"fmt"

	"github.com/donca/user-crud/config/storage"
	"github.com/donca/user-crud/internal/posts/interfaces"
	"github.com/donca/user-crud/internal/posts/models"

	"github.com/jackc/pgx/v5"
)

type postRepository struct {
	db *storage.PostgresPool
}

func NewPostRepository(db *storage.PostgresPool) interfaces.PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(ctx context.Context, authorID, content, visibility string) (*models.Post, error) {
	row := r.db.Pool.QueryRow(ctx, `
		INSERT INTO posts (author_id, content, visibility)
		VALUES ($1, $2, $3)
		RETURNING id, author_id, content, visibility, created_at, updated_at`,
		authorID, content, visibility)
	return scanPostBasic(row)
}

func (r *postRepository) FindByID(ctx context.Context, id string) (*models.Post, error) {
	row := r.db.Pool.QueryRow(ctx, `
		SELECT p.id, p.author_id, u.display_name, p.content, p.visibility, p.created_at, p.updated_at
		FROM posts p JOIN users u ON u.id = p.author_id
		WHERE p.id = $1`, id)
	p, err := scanPostWithAuthor(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return p, err
}

func (r *postRepository) Update(ctx context.Context, id, authorID string, content, visibility *string) (*models.Post, error) {
	row := r.db.Pool.QueryRow(ctx, `
		UPDATE posts SET
			content = COALESCE($3, content),
			visibility = COALESCE($4, visibility),
			updated_at = NOW()
		WHERE id = $1 AND author_id = $2
		RETURNING id, author_id, content, visibility, created_at, updated_at`,
		id, authorID, content, visibility)
	p, err := scanPostBasic(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return p, err
}

func (r *postRepository) Delete(ctx context.Context, id, authorID string) error {
	tag, err := r.db.Pool.Exec(ctx, `DELETE FROM posts WHERE id = $1 AND author_id = $2`, id, authorID)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *postRepository) ListByAuthor(ctx context.Context, authorID, viewerID string, limit, offset int) ([]models.Post, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT p.id, p.author_id, u.display_name, p.content, p.visibility, p.created_at, p.updated_at
		FROM posts p JOIN users u ON u.id = p.author_id
		WHERE p.author_id = $1
		  AND (
		    p.visibility = 'public'
		    OR p.author_id = $2
		    OR (
		      p.visibility = 'private'
		      AND EXISTS (
		        SELECT 1 FROM friendships f
		        WHERE f.status = 'accepted'
		          AND (
		            (f.requester_id = $2 AND f.addressee_id = p.author_id)
		            OR (f.addressee_id = $2 AND f.requester_id = p.author_id)
		          )
		      )
		    )
		  )
		ORDER BY p.created_at DESC LIMIT $3 OFFSET $4`,
		authorID, viewerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list posts by author: %w", err)
	}
	return scanPostRows(rows)
}

func (r *postRepository) ListFeed(ctx context.Context, viewerID string, limit, offset int) ([]models.Post, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT p.id, p.author_id, u.display_name, p.content, p.visibility, p.created_at, p.updated_at
		FROM posts p JOIN users u ON u.id = p.author_id
		WHERE
		  p.visibility = 'public'
		  OR p.author_id = $1
		  OR (
		    p.visibility = 'private'
		    AND EXISTS (
		      SELECT 1 FROM friendships f
		      WHERE f.status = 'accepted'
		        AND (
		          (f.requester_id = $1 AND f.addressee_id = p.author_id)
		          OR (f.addressee_id = $1 AND f.requester_id = p.author_id)
		        )
		    )
		  )
		ORDER BY p.created_at DESC LIMIT $2 OFFSET $3`,
		viewerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list feed: %w", err)
	}
	return scanPostRows(rows)
}

func scanPostBasic(row pgx.Row) (*models.Post, error) {
	var p models.Post
	if err := row.Scan(&p.ID, &p.AuthorID, &p.Content, &p.Visibility, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, fmt.Errorf("scan post: %w", err)
	}
	return &p, nil
}

func scanPostWithAuthor(row pgx.Row) (*models.Post, error) {
	var p models.Post
	if err := row.Scan(&p.ID, &p.AuthorID, &p.AuthorName, &p.Content, &p.Visibility, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, err
	}
	return &p, nil
}

func scanPostRows(rows pgx.Rows) ([]models.Post, error) {
	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.AuthorID, &p.AuthorName, &p.Content, &p.Visibility, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan post row: %w", err)
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}
