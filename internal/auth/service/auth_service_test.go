package service

import (
	"context"
	"os"
	"testing"
	"time"

	authmodels "github.com/donca/user-crud/internal/auth/models"
	usermocks "github.com/donca/user-crud/internal/users/mocks/interfaces"
	usermodels "github.com/donca/user-crud/internal/users/models"
	"github.com/donca/user-crud/pkg/kit/enums"
	kiterrors "github.com/donca/user-crud/pkg/kit/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestMain(m *testing.M) {
	_ = os.Setenv(enums.JWTSecret, "test-secret-key")
	os.Exit(m.Run())
}

func sampleUser() *usermodels.User {
	return &usermodels.User{
		ID: "u1", Username: "alice", Email: "alice@test.com",
		DisplayName: "Alice", CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
}

func validRegisterReq() authmodels.RegisterRequest {
	return authmodels.RegisterRequest{
		Username: "alice", Email: "alice@test.com", Password: "password123",
	}
}

func TestRegister_Success(t *testing.T) {
	repo := usermocks.NewMockUserRepository(t)
	repo.EXPECT().FindByUsername(mock.Anything, "alice").Return(nil, nil)
	repo.EXPECT().FindByEmail(mock.Anything, "alice@test.com").Return(nil, nil)
	repo.EXPECT().Create(mock.Anything, "alice", "alice@test.com", mock.AnythingOfType("string"), "alice").
		Return(sampleUser(), nil)
	svc := NewAuthService(repo)

	res, err := svc.Register(context.Background(), validRegisterReq())
	require.NoError(t, err)
	assert.NotEmpty(t, res.Token)
	assert.Equal(t, "u1", res.User.ID)
}

func TestRegister_DefaultDisplayName(t *testing.T) {
	repo := usermocks.NewMockUserRepository(t)
	repo.EXPECT().FindByUsername(mock.Anything, "bob").Return(nil, nil)
	repo.EXPECT().FindByEmail(mock.Anything, "bob@test.com").Return(nil, nil)
	repo.EXPECT().Create(mock.Anything, "bob", "bob@test.com", mock.AnythingOfType("string"), "bob").
		Return(sampleUser(), nil)
	svc := NewAuthService(repo)

	req := validRegisterReq()
	req.Username = "bob"
	req.Email = "bob@test.com"
	req.DisplayName = ""
	_, err := svc.Register(context.Background(), req)
	require.NoError(t, err)
}

func TestRegister_ValidationErrors(t *testing.T) {
	svc := NewAuthService(usermocks.NewMockUserRepository(t))
	cases := []struct {
		req  authmodels.RegisterRequest
		code string
		msg  string
	}{
		{authmodels.RegisterRequest{Username: "ab", Email: "a@b.com", Password: "password123"},
			kiterrors.CodeValidation, "username must be 3-50 characters"},
		{authmodels.RegisterRequest{Username: "alice", Email: "bad", Password: "password123"},
			kiterrors.CodeValidation, "valid email is required"},
		{authmodels.RegisterRequest{Username: "alice", Email: "a@b.com", Password: "short"},
			kiterrors.CodeValidation, "password must be at least 8 characters"},
	}
	for _, tc := range cases {
		_, err := svc.Register(context.Background(), tc.req)
		assertDomain(t, err, tc.code, tc.msg)
	}
}

func TestRegister_UsernameTaken(t *testing.T) {
	repo := usermocks.NewMockUserRepository(t)
	repo.EXPECT().FindByUsername(mock.Anything, "alice").Return(sampleUser(), nil)
	svc := NewAuthService(repo)

	_, err := svc.Register(context.Background(), validRegisterReq())
	assertDomain(t, err, kiterrors.CodeAlreadyExists, "username already taken")
}

func TestRegister_EmailTaken(t *testing.T) {
	repo := usermocks.NewMockUserRepository(t)
	repo.EXPECT().FindByUsername(mock.Anything, "alice").Return(nil, nil)
	repo.EXPECT().FindByEmail(mock.Anything, "alice@test.com").Return(sampleUser(), nil)
	svc := NewAuthService(repo)

	_, err := svc.Register(context.Background(), validRegisterReq())
	assertDomain(t, err, kiterrors.CodeAlreadyExists, "email already registered")
}

func TestRegister_CreateError(t *testing.T) {
	repo := usermocks.NewMockUserRepository(t)
	repo.EXPECT().FindByUsername(mock.Anything, "alice").Return(nil, nil)
	repo.EXPECT().FindByEmail(mock.Anything, "alice@test.com").Return(nil, nil)
	repo.EXPECT().Create(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, assert.AnError)
	svc := NewAuthService(repo)

	_, err := svc.Register(context.Background(), validRegisterReq())
	assertDomain(t, err, kiterrors.CodeInternal, "failed to register user")
}

func TestLogin_Validation(t *testing.T) {
	svc := NewAuthService(usermocks.NewMockUserRepository(t))
	_, err := svc.Login(context.Background(), authmodels.LoginRequest{})
	assertDomain(t, err, kiterrors.CodeValidation, "username and password are required")
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := usermocks.NewMockUserRepository(t)
	repo.EXPECT().FindByUsername(mock.Anything, "alice").Return(nil, nil)
	svc := NewAuthService(repo)

	_, err := svc.Login(context.Background(), authmodels.LoginRequest{Username: "alice", Password: "password123"})
	assertDomain(t, err, kiterrors.CodeUnauthorized, "invalid credentials")
}

func TestLogin_WrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("otherpass"), bcrypt.DefaultCost)
	repo := usermocks.NewMockUserRepository(t)
	repo.EXPECT().FindByUsername(mock.Anything, "alice").Return(sampleUser(), nil)
	repo.EXPECT().FindPasswordHash(mock.Anything, "u1").Return(string(hash), nil)
	svc := NewAuthService(repo)

	_, err := svc.Login(context.Background(), authmodels.LoginRequest{Username: "alice", Password: "password123"})
	assertDomain(t, err, kiterrors.CodeUnauthorized, "invalid credentials")
}

func TestLogin_NoPasswordHash(t *testing.T) {
	repo := usermocks.NewMockUserRepository(t)
	repo.EXPECT().FindByUsername(mock.Anything, "alice").Return(sampleUser(), nil)
	repo.EXPECT().FindPasswordHash(mock.Anything, "u1").Return("", nil)
	svc := NewAuthService(repo)

	_, err := svc.Login(context.Background(), authmodels.LoginRequest{Username: "alice", Password: "password123"})
	assertDomain(t, err, kiterrors.CodeUnauthorized, "invalid credentials")
}

func TestLogin_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	repo := usermocks.NewMockUserRepository(t)
	repo.EXPECT().FindByUsername(mock.Anything, "alice").Return(sampleUser(), nil)
	repo.EXPECT().FindPasswordHash(mock.Anything, "u1").Return(string(hash), nil)
	svc := NewAuthService(repo)

	res, err := svc.Login(context.Background(), authmodels.LoginRequest{Username: "alice", Password: "password123"})
	require.NoError(t, err)
	assert.NotEmpty(t, res.Token)
}

func assertDomain(t *testing.T, err error, code, msg string) {
	t.Helper()
	var de *kiterrors.DomainError
	require.ErrorAs(t, err, &de)
	assert.Equal(t, code, de.Code)
	assert.Equal(t, msg, de.Message)
}
