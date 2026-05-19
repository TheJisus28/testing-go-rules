package interfaces

import (
	"context"

	"github.com/donca/user-crud/internal/users/models"
)

type UserRepository interface {
	Create(ctx context.Context, username, email, passwordHash, displayName string) (*models.User, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindPasswordHash(ctx context.Context, id string) (string, error)
	Update(ctx context.Context, id string, email, displayName, bio, avatarURL *string) (*models.User, error)
	UpdateProfile(ctx context.Context, id string, displayName, bio, avatarURL *string) (*models.User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]models.UserProfile, error)
}

type UserService interface {
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetProfile(ctx context.Context, id string) (*models.UserProfile, error)
	List(ctx context.Context, limit, offset int) ([]models.UserProfile, error)
	Update(ctx context.Context, actorID, targetID string, req models.UpdateUserRequest) (*models.User, error)
	UpdateProfile(ctx context.Context, userID string, req models.UpdateProfileRequest) (*models.User, error)
	Delete(ctx context.Context, actorID, targetID string) error
}
