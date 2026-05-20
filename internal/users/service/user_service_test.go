package service

import (
	"context"
	"testing"
	"time"

	"github.com/donca/user-crud/internal/users/mocks/interfaces"
	"github.com/donca/user-crud/internal/users/models"
	kiterrors "github.com/donca/user-crud/pkg/kit/errors"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func sampleUser() *models.User {
	return &models.User{
		ID: "u1", Username: "alice", Email: "a@b.com",
		DisplayName: "Alice", CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
}

func TestGetByID_Success(t *testing.T) {
	repo := mocks.NewMockUserRepository(t)
	repo.EXPECT().FindByID(mock.Anything, "u1").Return(sampleUser(), nil)
	svc := NewUserService(repo)

	u, err := svc.GetByID(context.Background(), "u1")
	require.NoError(t, err)
	assert.Equal(t, "u1", u.ID)
}

func TestGetByID_NotFound(t *testing.T) {
	repo := mocks.NewMockUserRepository(t)
	repo.EXPECT().FindByID(mock.Anything, "x").Return(nil, nil)
	svc := NewUserService(repo)

	_, err := svc.GetByID(context.Background(), "x")
	assertDomain(t, err, kiterrors.CodeNotFound, "user not found")
}

func TestGetByID_RepositoryError(t *testing.T) {
	repo := mocks.NewMockUserRepository(t)
	repo.EXPECT().FindByID(mock.Anything, "u1").Return(nil, assert.AnError)
	svc := NewUserService(repo)

	_, err := svc.GetByID(context.Background(), "u1")
	assertDomain(t, err, kiterrors.CodeInternal, "failed to get user")
}

func TestGetProfile_Success(t *testing.T) {
	repo := mocks.NewMockUserRepository(t)
	repo.EXPECT().FindByID(mock.Anything, "u1").Return(sampleUser(), nil)
	svc := NewUserService(repo)

	p, err := svc.GetProfile(context.Background(), "u1")
	require.NoError(t, err)
	assert.Equal(t, "alice", p.Username)
}

func TestList_Success_DefaultPagination(t *testing.T) {
	repo := mocks.NewMockUserRepository(t)
	repo.EXPECT().List(mock.Anything, 20, 0).Return([]models.UserProfile{{ID: "u1"}}, nil)
	svc := NewUserService(repo)

	users, err := svc.List(context.Background(), 0, -1)
	require.NoError(t, err)
	assert.Len(t, users, 1)
}

func TestList_RepositoryError(t *testing.T) {
	repo := mocks.NewMockUserRepository(t)
	repo.EXPECT().List(mock.Anything, 20, 0).Return(nil, assert.AnError)
	svc := NewUserService(repo)

	_, err := svc.List(context.Background(), 0, 0)
	assertDomain(t, err, kiterrors.CodeInternal, "failed to list users")
}

func TestUpdate_Success(t *testing.T) {
	repo := mocks.NewMockUserRepository(t)
	name := "New"
	repo.EXPECT().Update(mock.Anything, "u1", mock.Anything, &name, mock.Anything, mock.Anything).
		Return(sampleUser(), nil)
	svc := NewUserService(repo)

	u, err := svc.Update(context.Background(), "u1", "u1", models.UpdateUserRequest{DisplayName: &name})
	require.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUpdate_Forbidden(t *testing.T) {
	svc := NewUserService(mocks.NewMockUserRepository(t))
	_, err := svc.Update(context.Background(), "u1", "u2", models.UpdateUserRequest{})
	assertDomain(t, err, kiterrors.CodeForbidden, "you can only update your own account")
}

func TestUpdate_NotFound(t *testing.T) {
	repo := mocks.NewMockUserRepository(t)
	repo.EXPECT().Update(mock.Anything, "u1", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil)
	svc := NewUserService(repo)

	_, err := svc.Update(context.Background(), "u1", "u1", models.UpdateUserRequest{})
	assertDomain(t, err, kiterrors.CodeNotFound, "user not found")
}

func TestDelete_Success(t *testing.T) {
	repo := mocks.NewMockUserRepository(t)
	repo.EXPECT().Delete(mock.Anything, "u1").Return(nil)
	svc := NewUserService(repo)
	require.NoError(t, svc.Delete(context.Background(), "u1", "u1"))
}

func TestDelete_Forbidden(t *testing.T) {
	svc := NewUserService(mocks.NewMockUserRepository(t))
	err := svc.Delete(context.Background(), "u1", "u2")
	assertDomain(t, err, kiterrors.CodeForbidden, "you can only delete your own account")
}

func TestDelete_NotFound(t *testing.T) {
	repo := mocks.NewMockUserRepository(t)
	repo.EXPECT().Delete(mock.Anything, "u1").Return(pgx.ErrNoRows)
	svc := NewUserService(repo)
	err := svc.Delete(context.Background(), "u1", "u1")
	assertDomain(t, err, kiterrors.CodeNotFound, "user not found")
}

func TestUpdateProfile_DelegatesToUpdate(t *testing.T) {
	repo := mocks.NewMockUserRepository(t)
	bio := "dev"
	repo.EXPECT().Update(mock.Anything, "u1", mock.Anything, mock.Anything, &bio, mock.Anything).
		Return(sampleUser(), nil)
	svc := NewUserService(repo)

	_, err := svc.UpdateProfile(context.Background(), "u1", models.UpdateProfileRequest{Bio: &bio})
	require.NoError(t, err)
}

func TestUpdate_RepositoryError(t *testing.T) {
	repo := mocks.NewMockUserRepository(t)
	repo.EXPECT().Update(mock.Anything, "u1", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, assert.AnError)
	svc := NewUserService(repo)

	_, err := svc.Update(context.Background(), "u1", "u1", models.UpdateUserRequest{})
	assertDomain(t, err, kiterrors.CodeInternal, "failed to update user")
}

func TestDelete_InternalError(t *testing.T) {
	repo := mocks.NewMockUserRepository(t)
	repo.EXPECT().Delete(mock.Anything, "u1").Return(assert.AnError)
	svc := NewUserService(repo)
	err := svc.Delete(context.Background(), "u1", "u1")
	assertDomain(t, err, kiterrors.CodeInternal, "failed to delete user")
}

func assertDomain(t *testing.T, err error, code, msg string) {
	t.Helper()
	var de *kiterrors.DomainError
	require.ErrorAs(t, err, &de)
	assert.Equal(t, code, de.Code)
	assert.Equal(t, msg, de.Message)
}
