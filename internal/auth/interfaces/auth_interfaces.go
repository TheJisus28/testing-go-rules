package interfaces

import (
	"context"

	"github.com/donca/user-crud/internal/auth/models"
)

type AuthService interface {
	Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error)
	Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error)
}
