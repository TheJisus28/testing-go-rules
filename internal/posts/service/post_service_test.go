package service

import (
	"context"
	"testing"
	"time"

	friendmocks "github.com/donca/user-crud/internal/friendships/mocks/interfaces"
	"github.com/donca/user-crud/internal/posts/mocks/interfaces"
	"github.com/donca/user-crud/internal/posts/models"
	kiterrors "github.com/donca/user-crud/pkg/kit/errors"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func samplePost() *models.Post {
	return &models.Post{
		ID: "p1", AuthorID: "u1", Content: "hello",
		Visibility: models.VisibilityPublic,
		CreatedAt:  time.Now(), UpdatedAt: time.Now(),
	}
}

func TestCreate_Success_PublicDefault(t *testing.T) {
	repo := mocks.NewMockPostRepository(t)
	repo.EXPECT().Create(mock.Anything, "u1", "hello", models.VisibilityPublic).Return(samplePost(), nil)
	svc := NewPostService(repo, friendmocks.NewMockFriendshipRepository(t))

	post, err := svc.Create(context.Background(), "u1", models.CreatePostRequest{Content: "hello"})
	require.NoError(t, err)
	assert.Equal(t, "p1", post.ID)
}

func TestCreate_Success_Private(t *testing.T) {
	repo := mocks.NewMockPostRepository(t)
	repo.EXPECT().Create(mock.Anything, "u1", "secret", models.VisibilityPrivate).Return(samplePost(), nil)
	svc := NewPostService(repo, friendmocks.NewMockFriendshipRepository(t))

	_, err := svc.Create(context.Background(), "u1", models.CreatePostRequest{
		Content: "secret", Visibility: models.VisibilityPrivate,
	})
	require.NoError(t, err)
}

func TestCreate_Validation(t *testing.T) {
	svc := NewPostService(mocks.NewMockPostRepository(t), friendmocks.NewMockFriendshipRepository(t))

	_, err := svc.Create(context.Background(), "u1", models.CreatePostRequest{})
	assertDomain(t, err, kiterrors.CodeValidation, "content is required")

	_, err = svc.Create(context.Background(), "u1", models.CreatePostRequest{
		Content: "x", Visibility: "friends-only",
	})
	assertDomain(t, err, kiterrors.CodeValidation, "visibility must be public or private")
}

func TestCreate_RepositoryError(t *testing.T) {
	repo := mocks.NewMockPostRepository(t)
	repo.EXPECT().Create(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError)
	svc := NewPostService(repo, friendmocks.NewMockFriendshipRepository(t))

	_, err := svc.Create(context.Background(), "u1", models.CreatePostRequest{Content: "hi"})
	assertDomain(t, err, kiterrors.CodeInternal, "failed to create post")
}

func TestGetByID_PublicSuccess(t *testing.T) {
	repo := mocks.NewMockPostRepository(t)
	repo.EXPECT().FindByID(mock.Anything, "p1").Return(samplePost(), nil)
	svc := NewPostService(repo, friendmocks.NewMockFriendshipRepository(t))

	post, err := svc.GetByID(context.Background(), "p1", "")
	require.NoError(t, err)
	assert.Equal(t, "p1", post.ID)
}

func TestGetByID_PrivateForbidden(t *testing.T) {
	p := samplePost()
	p.Visibility = models.VisibilityPrivate
	p.AuthorID = "u2"
	repo := mocks.NewMockPostRepository(t)
	repo.EXPECT().FindByID(mock.Anything, "p1").Return(p, nil)
	friends := friendmocks.NewMockFriendshipRepository(t)
	friends.EXPECT().AreFriends(mock.Anything, "u2", "u3").Return(false, nil)
	svc := NewPostService(repo, friends)

	_, err := svc.GetByID(context.Background(), "p1", "u3")
	assertDomain(t, err, kiterrors.CodeForbidden, "you cannot view this post")
}

func TestGetByID_PrivateAsAuthor(t *testing.T) {
	p := samplePost()
	p.Visibility = models.VisibilityPrivate
	repo := mocks.NewMockPostRepository(t)
	repo.EXPECT().FindByID(mock.Anything, "p1").Return(p, nil)
	svc := NewPostService(repo, friendmocks.NewMockFriendshipRepository(t))

	_, err := svc.GetByID(context.Background(), "p1", "u1")
	require.NoError(t, err)
}

func TestGetByID_PrivateAsFriend(t *testing.T) {
	p := samplePost()
	p.Visibility = models.VisibilityPrivate
	p.AuthorID = "u2"
	repo := mocks.NewMockPostRepository(t)
	repo.EXPECT().FindByID(mock.Anything, "p1").Return(p, nil)
	friends := friendmocks.NewMockFriendshipRepository(t)
	friends.EXPECT().AreFriends(mock.Anything, "u2", "u3").Return(true, nil)
	svc := NewPostService(repo, friends)

	_, err := svc.GetByID(context.Background(), "p1", "u3")
	require.NoError(t, err)
}

func TestGetByID_NotFound(t *testing.T) {
	repo := mocks.NewMockPostRepository(t)
	repo.EXPECT().FindByID(mock.Anything, "p1").Return(nil, nil)
	svc := NewPostService(repo, friendmocks.NewMockFriendshipRepository(t))

	_, err := svc.GetByID(context.Background(), "p1", "u1")
	assertDomain(t, err, kiterrors.CodeNotFound, "post not found")
}

func TestUpdate_Success(t *testing.T) {
	repo := mocks.NewMockPostRepository(t)
	repo.EXPECT().Update(mock.Anything, "p1", "u1", mock.Anything, mock.Anything).Return(samplePost(), nil)
	svc := NewPostService(repo, friendmocks.NewMockFriendshipRepository(t))

	_, err := svc.Update(context.Background(), "u1", "p1", models.UpdatePostRequest{})
	require.NoError(t, err)
}

func TestUpdate_InvalidVisibility(t *testing.T) {
	bad := "invalid"
	svc := NewPostService(mocks.NewMockPostRepository(t), friendmocks.NewMockFriendshipRepository(t))
	_, err := svc.Update(context.Background(), "u1", "p1", models.UpdatePostRequest{Visibility: &bad})
	assertDomain(t, err, kiterrors.CodeValidation, "visibility must be public or private")
}

func TestDelete_Success(t *testing.T) {
	repo := mocks.NewMockPostRepository(t)
	repo.EXPECT().Delete(mock.Anything, "p1", "u1").Return(nil)
	svc := NewPostService(repo, friendmocks.NewMockFriendshipRepository(t))
	require.NoError(t, svc.Delete(context.Background(), "u1", "p1"))
}

func TestDelete_NotFound(t *testing.T) {
	repo := mocks.NewMockPostRepository(t)
	repo.EXPECT().Delete(mock.Anything, "p1", "u1").Return(pgx.ErrNoRows)
	svc := NewPostService(repo, friendmocks.NewMockFriendshipRepository(t))
	err := svc.Delete(context.Background(), "u1", "p1")
	assertDomain(t, err, kiterrors.CodeNotFound, "post not found")
}

func TestFeed_Success(t *testing.T) {
	repo := mocks.NewMockPostRepository(t)
	repo.EXPECT().ListFeed(mock.Anything, "u1", 20, 0).Return([]models.Post{*samplePost()}, nil)
	svc := NewPostService(repo, friendmocks.NewMockFriendshipRepository(t))

	posts, err := svc.Feed(context.Background(), "u1", 0, 0)
	require.NoError(t, err)
	assert.Len(t, posts, 1)
}

func TestGetByID_RepositoryError(t *testing.T) {
	repo := mocks.NewMockPostRepository(t)
	repo.EXPECT().FindByID(mock.Anything, "p1").Return(nil, assert.AnError)
	svc := NewPostService(repo, friendmocks.NewMockFriendshipRepository(t))

	_, err := svc.GetByID(context.Background(), "p1", "u1")
	assertDomain(t, err, kiterrors.CodeInternal, "failed to get post")
}

func TestUpdate_NotFound(t *testing.T) {
	repo := mocks.NewMockPostRepository(t)
	repo.EXPECT().Update(mock.Anything, "p1", "u1", mock.Anything, mock.Anything).Return(nil, nil)
	svc := NewPostService(repo, friendmocks.NewMockFriendshipRepository(t))

	_, err := svc.Update(context.Background(), "u1", "p1", models.UpdatePostRequest{})
	assertDomain(t, err, kiterrors.CodeNotFound, "post not found")
}

func TestFeed_Error(t *testing.T) {
	repo := mocks.NewMockPostRepository(t)
	repo.EXPECT().ListFeed(mock.Anything, "u1", 20, 0).Return(nil, assert.AnError)
	svc := NewPostService(repo, friendmocks.NewMockFriendshipRepository(t))

	_, err := svc.Feed(context.Background(), "u1", 0, 0)
	assertDomain(t, err, kiterrors.CodeInternal, "failed to load feed")
}

func TestListByUser_RepositoryError(t *testing.T) {
	repo := mocks.NewMockPostRepository(t)
	repo.EXPECT().ListByAuthor(mock.Anything, "u2", "u1", 20, 0).Return(nil, assert.AnError)
	svc := NewPostService(repo, friendmocks.NewMockFriendshipRepository(t))

	_, err := svc.ListByUser(context.Background(), "u2", "u1", 0, 0)
	assertDomain(t, err, kiterrors.CodeInternal, "failed to list posts")
}

func assertDomain(t *testing.T, err error, code, msg string) {
	t.Helper()
	var de *kiterrors.DomainError
	require.ErrorAs(t, err, &de)
	assert.Equal(t, code, de.Code)
	assert.Equal(t, msg, de.Message)
}
