package interfaces

import (
	"context"

	"github.com/donca/user-crud/internal/posts/models"
)

type PostRepository interface {
	Create(ctx context.Context, authorID, content, visibility string) (*models.Post, error)
	FindByID(ctx context.Context, id string) (*models.Post, error)
	Update(ctx context.Context, id, authorID string, content, visibility *string) (*models.Post, error)
	Delete(ctx context.Context, id, authorID string) error
	ListByAuthor(ctx context.Context, authorID, viewerID string, limit, offset int) ([]models.Post, error)
	ListFeed(ctx context.Context, viewerID string, limit, offset int) ([]models.Post, error)
}

type PostService interface {
	Create(ctx context.Context, authorID string, req models.CreatePostRequest) (*models.Post, error)
	GetByID(ctx context.Context, postID, viewerID string) (*models.Post, error)
	Update(ctx context.Context, authorID, postID string, req models.UpdatePostRequest) (*models.Post, error)
	Delete(ctx context.Context, authorID, postID string) error
	ListByUser(ctx context.Context, authorID, viewerID string, limit, offset int) ([]models.Post, error)
	Feed(ctx context.Context, viewerID string, limit, offset int) ([]models.Post, error)
}
