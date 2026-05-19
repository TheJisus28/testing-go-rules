package handler

import (
	"github.com/donca/user-crud/internal/auth/interfaces"
	"github.com/donca/user-crud/internal/auth/models"
	kiterrors "github.com/donca/user-crud/pkg/kit/errors"
	"github.com/donca/user-crud/pkg/kit/wrapper"

	"github.com/labstack/echo/v4"
)

type AuthHandler interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
}

type authHandler struct {
	service interfaces.AuthService
}

func NewAuthHandler(service interfaces.AuthService) AuthHandler {
	return &authHandler{service: service}
}

// Register godoc
// @Summary      Register a new user
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.RegisterRequest  true  "Registration"
// @Success      201  {object}  models.AuthResponse
// @Router       /v1/auth/register [post]
func (h *authHandler) Register(c echo.Context) error {
	var req models.RegisterRequest
	if err := c.Bind(&req); err != nil {
		resp, status := wrapper.FromError(kiterrors.Validation("invalid request body"))
		return c.JSON(status, resp)
	}
	result, err := h.service.Register(c.Request().Context(), req)
	if err != nil {
		resp, status := wrapper.FromError(err)
		return c.JSON(status, resp)
	}
	return c.JSON(201, wrapper.Response[models.AuthResponse]{
		Status:       "success",
		StatusReason: "ok",
		Data:         *result,
	})
}

// Login godoc
// @Summary      Login
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.LoginRequest  true  "Credentials"
// @Success      200  {object}  models.AuthResponse
// @Router       /v1/auth/login [post]
func (h *authHandler) Login(c echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		resp, status := wrapper.FromError(kiterrors.Validation("invalid request body"))
		return c.JSON(status, resp)
	}
	result, err := h.service.Login(c.Request().Context(), req)
	if err != nil {
		resp, status := wrapper.FromError(err)
		return c.JSON(status, resp)
	}
	return wrapper.GenerateResponse(c, *result, nil)
}
