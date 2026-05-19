package handler

import (
	"strconv"

	"github.com/donca/user-crud/config/generals/logger"
	"github.com/donca/user-crud/internal/users/interfaces"
	"github.com/donca/user-crud/internal/users/models"
	authmw "github.com/donca/user-crud/pkg/middleware/auth"
	kiterrors "github.com/donca/user-crud/pkg/kit/errors"
	"github.com/donca/user-crud/pkg/kit/wrapper"

	"github.com/labstack/echo/v4"
)

type UserHandler interface {
	GetByID(c echo.Context) error
	GetProfile(c echo.Context) error
	List(c echo.Context) error
	Update(c echo.Context) error
	UpdateProfile(c echo.Context) error
	Delete(c echo.Context) error
}

type userHandler struct {
	service interfaces.UserService
}

func NewUserHandler(service interfaces.UserService) UserHandler {
	return &userHandler{service: service}
}

// GetByID godoc
// @Summary      Get user by ID
// @Description  Returns full user data (authenticated owner only for email)
// @Tags         Users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  models.UserResponse
// @Failure      404  {object}  wrapper.Response[any]
// @Router       /v1/users/{id} [get]
func (h *userHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	user, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		resp, status := wrapper.FromError(err)
		return c.JSON(status, resp)
	}
	return wrapper.GenerateResponse(c, models.UserResponse{User: *user}, nil)
}

// GetProfile godoc
// @Summary      Get public profile
// @Tags         Users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  models.ProfileResponse
// @Router       /v1/users/{id}/profile [get]
func (h *userHandler) GetProfile(c echo.Context) error {
	profile, err := h.service.GetProfile(c.Request().Context(), c.Param("id"))
	if err != nil {
		resp, status := wrapper.FromError(err)
		return c.JSON(status, resp)
	}
	return wrapper.GenerateResponse(c, models.ProfileResponse{Profile: *profile}, nil)
}

// List godoc
// @Summary      List user profiles
// @Tags         Users
// @Produce      json
// @Param        limit   query     int  false  "Limit"
// @Param        offset  query     int  false  "Offset"
// @Success      200  {object}  models.UsersListResponse
// @Router       /v1/users [get]
func (h *userHandler) List(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	users, err := h.service.List(c.Request().Context(), limit, offset)
	if err != nil {
		resp, status := wrapper.FromError(err)
		return c.JSON(status, resp)
	}
	return wrapper.GenerateResponse(c, models.UsersListResponse{Users: users}, nil)
}

// Update godoc
// @Summary      Update user account
// @Tags         Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      string  true  "User ID"
// @Param        body  body      models.UpdateUserRequest  true  "Update payload"
// @Success      200  {object}  models.UserResponse
// @Router       /v1/users/{id} [put]
func (h *userHandler) Update(c echo.Context) error {
	actorID := authmw.UserIDFromContext(c)
	var req models.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		resp, status := wrapper.FromError(kiterrors.Validation("invalid request body"))
		return c.JSON(status, resp)
	}
	user, err := h.service.Update(c.Request().Context(), actorID, c.Param("id"), req)
	if err != nil {
		log := logger.FromCtx(c.Request().Context())
		log.Warn().Err(err).Msg("users: update rejected")
		resp, status := wrapper.FromError(err)
		return c.JSON(status, resp)
	}
	return wrapper.GenerateResponse(c, models.UserResponse{User: *user}, nil)
}

// UpdateProfile godoc
// @Summary      Update own profile
// @Tags         Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      models.UpdateProfileRequest  true  "Profile payload"
// @Success      200  {object}  models.UserResponse
// @Router       /v1/me/profile [patch]
func (h *userHandler) UpdateProfile(c echo.Context) error {
	var req models.UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		resp, status := wrapper.FromError(kiterrors.Validation("invalid request body"))
		return c.JSON(status, resp)
	}
	user, err := h.service.UpdateProfile(c.Request().Context(), authmw.UserIDFromContext(c), req)
	if err != nil {
		resp, status := wrapper.FromError(err)
		return c.JSON(status, resp)
	}
	return wrapper.GenerateResponse(c, models.UserResponse{User: *user}, nil)
}

// Delete godoc
// @Summary      Delete user account
// @Tags         Users
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  wrapper.Response[any]
// @Router       /v1/users/{id} [delete]
func (h *userHandler) Delete(c echo.Context) error {
	if err := h.service.Delete(c.Request().Context(), authmw.UserIDFromContext(c), c.Param("id")); err != nil {
		resp, status := wrapper.FromError(err)
		return c.JSON(status, resp)
	}
	return wrapper.GenerateResponse(c, map[string]string{"deleted": "true"}, nil)
}
