package service

import (
	"context"
	"os"
	"strings"

	"github.com/donca/user-crud/config/generals/logger"
	"github.com/donca/user-crud/internal/auth/interfaces"
	authmodels "github.com/donca/user-crud/internal/auth/models"
	userinterfaces "github.com/donca/user-crud/internal/users/interfaces"
	"github.com/donca/user-crud/pkg/kit/enums"
	kiterrors "github.com/donca/user-crud/pkg/kit/errors"
	"github.com/donca/user-crud/pkg/kit/jwt"

	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	users userinterfaces.UserRepository
}

func NewAuthService(users userinterfaces.UserRepository) interfaces.AuthService {
	return &authService{users: users}
}

func (s *authService) Register(ctx context.Context, req authmodels.RegisterRequest) (*authmodels.AuthResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}
	if existing, _ := s.users.FindByUsername(ctx, req.Username); existing != nil {
		return nil, kiterrors.AlreadyExists("username already taken")
	}
	if existing, _ := s.users.FindByEmail(ctx, req.Email); existing != nil {
		return nil, kiterrors.AlreadyExists("email already registered")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.FromCtx(ctx).Error().Err(err).Msg("auth: hash password failed")
		return nil, kiterrors.Internal("failed to register user")
	}
	displayName := req.DisplayName
	if displayName == "" {
		displayName = req.Username
	}
	user, err := s.users.Create(ctx, req.Username, strings.ToLower(req.Email), string(hash), displayName)
	if err != nil {
		logger.FromCtx(ctx).Error().Err(err).Msg("auth: create user failed")
		return nil, kiterrors.Internal("failed to register user")
	}
	token, err := jwt.GenerateToken(user.ID, os.Getenv(enums.JWTSecret))
	if err != nil {
		return nil, kiterrors.Internal("failed to generate token")
	}
	return &authmodels.AuthResponse{Token: token, User: *user}, nil
}

func (s *authService) Login(ctx context.Context, req authmodels.LoginRequest) (*authmodels.AuthResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, kiterrors.Validation("username and password are required")
	}
	user, err := s.users.FindByUsername(ctx, req.Username)
	if err != nil || user == nil {
		return nil, kiterrors.Unauthorized("invalid credentials")
	}
	hash, err := s.users.FindPasswordHash(ctx, user.ID)
	if err != nil || hash == "" {
		return nil, kiterrors.Unauthorized("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		return nil, kiterrors.Unauthorized("invalid credentials")
	}
	token, err := jwt.GenerateToken(user.ID, os.Getenv(enums.JWTSecret))
	if err != nil {
		return nil, kiterrors.Internal("failed to generate token")
	}
	return &authmodels.AuthResponse{Token: token, User: *user}, nil
}

func validateRegister(req authmodels.RegisterRequest) error {
	if len(req.Username) < 3 || len(req.Username) > 50 {
		return kiterrors.Validation("username must be 3-50 characters")
	}
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		return kiterrors.Validation("valid email is required")
	}
	if len(req.Password) < 8 {
		return kiterrors.Validation("password must be at least 8 characters")
	}
	return nil
}
