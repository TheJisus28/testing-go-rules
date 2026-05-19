package service

import (
	"context"

	"github.com/donca/user-crud/config/generals/logger"
	"github.com/donca/user-crud/internal/users/interfaces"
	"github.com/donca/user-crud/internal/users/models"
	kiterrors "github.com/donca/user-crud/pkg/kit/errors"

	"github.com/jackc/pgx/v5"
)

type userService struct {
	repo interfaces.UserRepository
}

func NewUserService(repo interfaces.UserRepository) interfaces.UserService {
	return &userService{repo: repo}
}

func (s *userService) GetByID(ctx context.Context, id string) (*models.User, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		logger.FromCtx(ctx).Error().Err(err).Str("user_id", id).Msg("users: find failed")
		return nil, kiterrors.Internal("failed to get user")
	}
	if u == nil {
		return nil, kiterrors.NotFound("user not found")
	}
	return u, nil
}

func (s *userService) GetProfile(ctx context.Context, id string) (*models.UserProfile, error) {
	u, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toProfile(u), nil
}

func (s *userService) List(ctx context.Context, limit, offset int) ([]models.UserProfile, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	users, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		logger.FromCtx(ctx).Error().Err(err).Msg("users: list failed")
		return nil, kiterrors.Internal("failed to list users")
	}
	return users, nil
}

func (s *userService) Update(ctx context.Context, actorID, targetID string, req models.UpdateUserRequest) (*models.User, error) {
	if actorID != targetID {
		return nil, kiterrors.Forbidden("you can only update your own account")
	}
	u, err := s.repo.Update(ctx, targetID, req.Email, req.DisplayName, req.Bio, req.AvatarURL)
	if err != nil {
		logger.FromCtx(ctx).Error().Err(err).Str("user_id", targetID).Msg("users: update failed")
		return nil, kiterrors.Internal("failed to update user")
	}
	if u == nil {
		return nil, kiterrors.NotFound("user not found")
	}
	return u, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID string, req models.UpdateProfileRequest) (*models.User, error) {
	return s.Update(ctx, userID, userID, models.UpdateUserRequest{
		DisplayName: req.DisplayName,
		Bio:         req.Bio,
		AvatarURL:   req.AvatarURL,
	})
}

func (s *userService) Delete(ctx context.Context, actorID, targetID string) error {
	if actorID != targetID {
		return kiterrors.Forbidden("you can only delete your own account")
	}
	if err := s.repo.Delete(ctx, targetID); err != nil {
		if err == pgx.ErrNoRows {
			return kiterrors.NotFound("user not found")
		}
		logger.FromCtx(ctx).Error().Err(err).Str("user_id", targetID).Msg("users: delete failed")
		return kiterrors.Internal("failed to delete user")
	}
	return nil
}

func toProfile(u *models.User) *models.UserProfile {
	return &models.UserProfile{
		ID:          u.ID,
		Username:    u.Username,
		DisplayName: u.DisplayName,
		Bio:         u.Bio,
		AvatarURL:   u.AvatarURL,
		CreatedAt:   u.CreatedAt,
	}
}
