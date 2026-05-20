package service

import (
	"context"
	"testing"
	"time"

	"github.com/donca/user-crud/internal/friendships/mocks/interfaces"
	"github.com/donca/user-crud/internal/friendships/models"
	usermocks "github.com/donca/user-crud/internal/users/mocks/interfaces"
	usermodels "github.com/donca/user-crud/internal/users/models"
	kiterrors "github.com/donca/user-crud/pkg/kit/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func sampleFriendship(status string) *models.Friendship {
	return &models.Friendship{
		ID: "f1", RequesterID: "u1", AddresseeID: "u2",
		Status: status, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
}

func sampleUser() *usermodels.User {
	return &usermodels.User{ID: "u2", Username: "bob"}
}

func TestSendRequest_Success(t *testing.T) {
	repo := mocks.NewMockFriendshipRepository(t)
	users := usermocks.NewMockUserRepository(t)
	users.EXPECT().FindByID(mock.Anything, "u2").Return(sampleUser(), nil)
	repo.EXPECT().FindBetween(mock.Anything, "u1", "u2").Return(nil, nil)
	repo.EXPECT().CreateRequest(mock.Anything, "u1", "u2").Return(sampleFriendship(models.StatusPending), nil)
	svc := NewFriendshipService(repo, users)

	f, err := svc.SendRequest(context.Background(), "u1", models.FriendRequest{AddresseeID: "u2"})
	require.NoError(t, err)
	assert.Equal(t, models.StatusPending, f.Status)
}

func TestSendRequest_Validation(t *testing.T) {
	svc := NewFriendshipService(mocks.NewMockFriendshipRepository(t), usermocks.NewMockUserRepository(t))

	_, err := svc.SendRequest(context.Background(), "u1", models.FriendRequest{})
	assertDomain(t, err, kiterrors.CodeValidation, "addressee_id is required")

	_, err = svc.SendRequest(context.Background(), "u1", models.FriendRequest{AddresseeID: "u1"})
	assertDomain(t, err, kiterrors.CodeValidation, "cannot send friend request to yourself")
}

func TestSendRequest_UserNotFound(t *testing.T) {
	users := usermocks.NewMockUserRepository(t)
	users.EXPECT().FindByID(mock.Anything, "u2").Return(nil, nil)
	svc := NewFriendshipService(mocks.NewMockFriendshipRepository(t), users)

	_, err := svc.SendRequest(context.Background(), "u1", models.FriendRequest{AddresseeID: "u2"})
	assertDomain(t, err, kiterrors.CodeNotFound, "user not found")
}

func TestSendRequest_Conflicts(t *testing.T) {
	users := usermocks.NewMockUserRepository(t)
	users.EXPECT().FindByID(mock.Anything, "u2").Return(sampleUser(), nil)
	cases := []struct {
		status string
		msg    string
	}{
		{models.StatusAccepted, "already friends"},
		{models.StatusPending, "friend request already pending"},
		{models.StatusRejected, "friend request was rejected; cannot resend yet"},
	}
	for _, tc := range cases {
		repo := mocks.NewMockFriendshipRepository(t)
		repo.EXPECT().FindBetween(mock.Anything, "u1", "u2").Return(sampleFriendship(tc.status), nil)
		svc := NewFriendshipService(repo, users)
		_, err := svc.SendRequest(context.Background(), "u1", models.FriendRequest{AddresseeID: "u2"})
		assertDomain(t, err, kiterrors.CodeConflict, tc.msg)
	}
}

func TestSendRequest_CreateError(t *testing.T) {
	repo := mocks.NewMockFriendshipRepository(t)
	users := usermocks.NewMockUserRepository(t)
	users.EXPECT().FindByID(mock.Anything, "u2").Return(sampleUser(), nil)
	repo.EXPECT().FindBetween(mock.Anything, "u1", "u2").Return(nil, nil)
	repo.EXPECT().CreateRequest(mock.Anything, "u1", "u2").Return(nil, assert.AnError)
	svc := NewFriendshipService(repo, users)

	_, err := svc.SendRequest(context.Background(), "u1", models.FriendRequest{AddresseeID: "u2"})
	assertDomain(t, err, kiterrors.CodeInternal, "failed to send friend request")
}

func TestAccept_Success(t *testing.T) {
	repo := mocks.NewMockFriendshipRepository(t)
	repo.EXPECT().UpdateStatus(mock.Anything, "f1", "u2", models.StatusAccepted).
		Return(sampleFriendship(models.StatusAccepted), nil)
	svc := NewFriendshipService(repo, usermocks.NewMockUserRepository(t))

	f, err := svc.Accept(context.Background(), "u2", "f1")
	require.NoError(t, err)
	assert.Equal(t, models.StatusAccepted, f.Status)
}

func TestReject_NotFound(t *testing.T) {
	repo := mocks.NewMockFriendshipRepository(t)
	repo.EXPECT().UpdateStatus(mock.Anything, "f1", "u2", models.StatusRejected).Return(nil, nil)
	svc := NewFriendshipService(repo, usermocks.NewMockUserRepository(t))

	_, err := svc.Reject(context.Background(), "u2", "f1")
	assertDomain(t, err, kiterrors.CodeNotFound, "friend request not found or already handled")
}

func TestListPendingReceived_Success(t *testing.T) {
	repo := mocks.NewMockFriendshipRepository(t)
	repo.EXPECT().ListPendingReceived(mock.Anything, "u2").
		Return([]models.Friendship{*sampleFriendship(models.StatusPending)}, nil)
	svc := NewFriendshipService(repo, usermocks.NewMockUserRepository(t))

	list, err := svc.ListPendingReceived(context.Background(), "u2")
	require.NoError(t, err)
	assert.Len(t, list, 1)
}

func TestListPendingSent_Error(t *testing.T) {
	repo := mocks.NewMockFriendshipRepository(t)
	repo.EXPECT().ListPendingSent(mock.Anything, "u1").Return(nil, assert.AnError)
	svc := NewFriendshipService(repo, usermocks.NewMockUserRepository(t))

	_, err := svc.ListPendingSent(context.Background(), "u1")
	assertDomain(t, err, kiterrors.CodeInternal, "failed to list requests")
}

func TestListFriends_Success(t *testing.T) {
	repo := mocks.NewMockFriendshipRepository(t)
	repo.EXPECT().ListFriends(mock.Anything, "u1").
		Return([]models.FriendProfile{{ID: "u2", Username: "bob"}}, nil)
	svc := NewFriendshipService(repo, usermocks.NewMockUserRepository(t))

	friends, err := svc.ListFriends(context.Background(), "u1")
	require.NoError(t, err)
	assert.Len(t, friends, 1)
}

func assertDomain(t *testing.T, err error, code, msg string) {
	t.Helper()
	var de *kiterrors.DomainError
	require.ErrorAs(t, err, &de)
	assert.Equal(t, code, de.Code)
	assert.Equal(t, msg, de.Message)
}
