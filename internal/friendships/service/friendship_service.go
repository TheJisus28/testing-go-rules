package service

import (
	"context"

	"github.com/donca/user-crud/config/generals/logger"
	userinterfaces "github.com/donca/user-crud/internal/users/interfaces"
	"github.com/donca/user-crud/internal/friendships/interfaces"
	"github.com/donca/user-crud/internal/friendships/models"
	kiterrors "github.com/donca/user-crud/pkg/kit/errors"
)

type friendshipService struct {
	repo  interfaces.FriendshipRepository
	users userinterfaces.UserRepository
}

func NewFriendshipService(repo interfaces.FriendshipRepository, users userinterfaces.UserRepository) interfaces.FriendshipService {
	return &friendshipService{repo: repo, users: users}
}

func (s *friendshipService) SendRequest(ctx context.Context, requesterID string, req models.FriendRequest) (*models.Friendship, error) {
	if req.AddresseeID == "" {
		return nil, kiterrors.Validation("addressee_id is required")
	}
	if requesterID == req.AddresseeID {
		return nil, kiterrors.Validation("cannot send friend request to yourself")
	}
	target, err := s.users.FindByID(ctx, req.AddresseeID)
	if err != nil || target == nil {
		return nil, kiterrors.NotFound("user not found")
	}
	existing, _ := s.repo.FindBetween(ctx, requesterID, req.AddresseeID)
	if existing != nil {
		switch existing.Status {
		case models.StatusAccepted:
			return nil, kiterrors.Conflict("already friends")
		case models.StatusPending:
			return nil, kiterrors.Conflict("friend request already pending")
		case models.StatusRejected:
			return nil, kiterrors.Conflict("friend request was rejected; cannot resend yet")
		}
	}
	f, err := s.repo.CreateRequest(ctx, requesterID, req.AddresseeID)
	if err != nil {
		logger.FromCtx(ctx).Error().Err(err).Msg("friendships: create request failed")
		return nil, kiterrors.Internal("failed to send friend request")
	}
	return f, nil
}

func (s *friendshipService) Accept(ctx context.Context, userID, requestID string) (*models.Friendship, error) {
	return s.respond(ctx, userID, requestID, models.StatusAccepted)
}

func (s *friendshipService) Reject(ctx context.Context, userID, requestID string) (*models.Friendship, error) {
	return s.respond(ctx, userID, requestID, models.StatusRejected)
}

func (s *friendshipService) respond(ctx context.Context, userID, requestID, status string) (*models.Friendship, error) {
	f, err := s.repo.UpdateStatus(ctx, requestID, userID, status)
	if err != nil {
		logger.FromCtx(ctx).Error().Err(err).Msg("friendships: update status failed")
		return nil, kiterrors.Internal("failed to update friend request")
	}
	if f == nil {
		return nil, kiterrors.NotFound("friend request not found or already handled")
	}
	return f, nil
}

func (s *friendshipService) ListPendingReceived(ctx context.Context, userID string) ([]models.Friendship, error) {
	list, err := s.repo.ListPendingReceived(ctx, userID)
	if err != nil {
		logger.FromCtx(ctx).Error().Err(err).Msg("friendships: list pending received failed")
		return nil, kiterrors.Internal("failed to list requests")
	}
	return list, nil
}

func (s *friendshipService) ListPendingSent(ctx context.Context, userID string) ([]models.Friendship, error) {
	list, err := s.repo.ListPendingSent(ctx, userID)
	if err != nil {
		logger.FromCtx(ctx).Error().Err(err).Msg("friendships: list pending sent failed")
		return nil, kiterrors.Internal("failed to list requests")
	}
	return list, nil
}

func (s *friendshipService) ListFriends(ctx context.Context, userID string) ([]models.FriendProfile, error) {
	friends, err := s.repo.ListFriends(ctx, userID)
	if err != nil {
		logger.FromCtx(ctx).Error().Err(err).Msg("friendships: list friends failed")
		return nil, kiterrors.Internal("failed to list friends")
	}
	return friends, nil
}
