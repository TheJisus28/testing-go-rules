package interfaces

import (
	"context"

	"github.com/donca/user-crud/internal/friendships/models"
)

type FriendshipRepository interface {
	CreateRequest(ctx context.Context, requesterID, addresseeID string) (*models.Friendship, error)
	FindByID(ctx context.Context, id string) (*models.Friendship, error)
	FindBetween(ctx context.Context, userA, userB string) (*models.Friendship, error)
	UpdateStatus(ctx context.Context, id, addresseeID, status string) (*models.Friendship, error)
	ListPendingReceived(ctx context.Context, userID string) ([]models.Friendship, error)
	ListPendingSent(ctx context.Context, userID string) ([]models.Friendship, error)
	ListFriends(ctx context.Context, userID string) ([]models.FriendProfile, error)
	AreFriends(ctx context.Context, userA, userB string) (bool, error)
}

type FriendshipService interface {
	SendRequest(ctx context.Context, requesterID string, req models.FriendRequest) (*models.Friendship, error)
	Accept(ctx context.Context, userID, requestID string) (*models.Friendship, error)
	Reject(ctx context.Context, userID, requestID string) (*models.Friendship, error)
	ListPendingReceived(ctx context.Context, userID string) ([]models.Friendship, error)
	ListPendingSent(ctx context.Context, userID string) ([]models.Friendship, error)
	ListFriends(ctx context.Context, userID string) ([]models.FriendProfile, error)
}
