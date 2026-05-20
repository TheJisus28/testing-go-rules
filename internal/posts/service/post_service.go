package service

import (
	"context"

	"github.com/donca/user-crud/config/generals/logger"
	friendinterfaces "github.com/donca/user-crud/internal/friendships/interfaces"
	"github.com/donca/user-crud/internal/posts/interfaces"
	"github.com/donca/user-crud/internal/posts/models"
	kiterrors "github.com/donca/user-crud/pkg/kit/errors"

	"github.com/jackc/pgx/v5"
)

type postService struct {
	repo      interfaces.PostRepository
	friendships friendinterfaces.FriendshipRepository
}

// NewPostService wires post persistence with friendship checks for private visibility.
func NewPostService(repo interfaces.PostRepository, friendships friendinterfaces.FriendshipRepository) interfaces.PostService {
	return &postService{repo: repo, friendships: friendships}
}

func (s *postService) Create(ctx context.Context, authorID string, req models.CreatePostRequest) (*models.Post, error) {
	if req.Content == "" {
		return nil, kiterrors.Validation("content is required")
	}
	visibility := req.Visibility
	if visibility == "" {
		visibility = models.VisibilityPublic
	}
	if visibility != models.VisibilityPublic && visibility != models.VisibilityPrivate {
		return nil, kiterrors.Validation("visibility must be public or private")
	}
	post, err := s.repo.Create(ctx, authorID, req.Content, visibility)
	if err != nil {
		logger.FromCtx(ctx).Error().Err(err).Msg("posts: create failed")
		return nil, kiterrors.Internal("failed to create post")
	}
	return post, nil
}

func (s *postService) GetByID(ctx context.Context, postID, viewerID string) (*models.Post, error) {
	post, err := s.repo.FindByID(ctx, postID)
	if err != nil {
		logger.FromCtx(ctx).Error().Err(err).Msg("posts: find failed")
		return nil, kiterrors.Internal("failed to get post")
	}
	if post == nil {
		return nil, kiterrors.NotFound("post not found")
	}
	if !s.canView(ctx, post, viewerID) {
		return nil, kiterrors.Forbidden("you cannot view this post")
	}
	return post, nil
}

func (s *postService) Update(ctx context.Context, authorID, postID string, req models.UpdatePostRequest) (*models.Post, error) {
	if req.Visibility != nil && *req.Visibility != models.VisibilityPublic && *req.Visibility != models.VisibilityPrivate {
		return nil, kiterrors.Validation("visibility must be public or private")
	}
	post, err := s.repo.Update(ctx, postID, authorID, req.Content, req.Visibility)
	if err != nil {
		logger.FromCtx(ctx).Error().Err(err).Msg("posts: update failed")
		return nil, kiterrors.Internal("failed to update post")
	}
	if post == nil {
		return nil, kiterrors.NotFound("post not found")
	}
	return post, nil
}

func (s *postService) Delete(ctx context.Context, authorID, postID string) error {
	if err := s.repo.Delete(ctx, postID, authorID); err != nil {
		if err == pgx.ErrNoRows {
			return kiterrors.NotFound("post not found")
		}
		logger.FromCtx(ctx).Error().Err(err).Msg("posts: delete failed")
		return kiterrors.Internal("failed to delete post")
	}
	return nil
}

func (s *postService) ListByUser(ctx context.Context, authorID, viewerID string, limit, offset int) ([]models.Post, error) {
	limit, offset = paginate(limit, offset)
	posts, err := s.repo.ListByAuthor(ctx, authorID, viewerID, limit, offset)
	if err != nil {
		logger.FromCtx(ctx).Error().Err(err).Msg("posts: list by user failed")
		return nil, kiterrors.Internal("failed to list posts")
	}
	return posts, nil
}

func (s *postService) Feed(ctx context.Context, viewerID string, limit, offset int) ([]models.Post, error) {
	limit, offset = paginate(limit, offset)
	posts, err := s.repo.ListFeed(ctx, viewerID, limit, offset)
	if err != nil {
		logger.FromCtx(ctx).Error().Err(err).Msg("posts: feed failed")
		return nil, kiterrors.Internal("failed to load feed")
	}
	return posts, nil
}

// canView mirrors feed filtering: public posts, own posts, or private posts visible to friends.
func (s *postService) canView(ctx context.Context, post *models.Post, viewerID string) bool {
	if post.Visibility == models.VisibilityPublic {
		return true
	}
	if post.AuthorID == viewerID {
		return true
	}
	if viewerID == "" {
		return false
	}
	ok, _ := s.friendships.AreFriends(ctx, post.AuthorID, viewerID)
	return ok
}

func paginate(limit, offset int) (int, int) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
