package models

import "github.com/donca/user-crud/internal/users/models"

type RegisterRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  models.User  `json:"user"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
