package handler

import (
	"github.com/donca/user-crud/internal/friendships/interfaces"
	"github.com/donca/user-crud/internal/friendships/models"
	authmw "github.com/donca/user-crud/pkg/middleware/auth"
	kiterrors "github.com/donca/user-crud/pkg/kit/errors"
	"github.com/donca/user-crud/pkg/kit/wrapper"

	"github.com/labstack/echo/v4"
)

type FriendshipHandler interface {
	SendRequest(c echo.Context) error
	Accept(c echo.Context) error
	Reject(c echo.Context) error
	ListPendingReceived(c echo.Context) error
	ListPendingSent(c echo.Context) error
	ListFriends(c echo.Context) error
}

type friendshipHandler struct {
	service interfaces.FriendshipService
}

func NewFriendshipHandler(service interfaces.FriendshipService) FriendshipHandler {
	return &friendshipHandler{service: service}
}

// SendRequest godoc
// @Summary      Send a friend request
// @Tags         Friendships
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      models.FriendRequest  true  "Friend request"
// @Success      200   {object}  models.FriendshipResponse
// @Router       /v1/friendships/requests [post]
func (h *friendshipHandler) SendRequest(c echo.Context) error {
	var req models.FriendRequest
	if err := c.Bind(&req); err != nil {
		resp, status := wrapper.FromError(kiterrors.Validation("invalid request body"))
		return c.JSON(status, resp)
	}
	f, err := h.service.SendRequest(c.Request().Context(), authmw.UserIDFromContext(c), req)
	if err != nil {
		resp, status := wrapper.FromError(err)
		return c.JSON(status, resp)
	}
	return wrapper.GenerateResponse(c, models.FriendshipResponse{Friendship: *f}, nil)
}

// Accept godoc
// @Summary      Accept a friend request
// @Tags         Friendships
// @Security     BearerAuth
// @Param        id   path      string  true  "Request ID"
// @Success      200  {object}  models.FriendshipResponse
// @Router       /v1/friendships/requests/{id}/accept [post]
func (h *friendshipHandler) Accept(c echo.Context) error {
	f, err := h.service.Accept(c.Request().Context(), authmw.UserIDFromContext(c), c.Param("id"))
	if err != nil {
		resp, status := wrapper.FromError(err)
		return c.JSON(status, resp)
	}
	return wrapper.GenerateResponse(c, models.FriendshipResponse{Friendship: *f}, nil)
}

// Reject godoc
// @Summary      Reject a friend request
// @Tags         Friendships
// @Security     BearerAuth
// @Param        id   path      string  true  "Request ID"
// @Success      200  {object}  models.FriendshipResponse
// @Router       /v1/friendships/requests/{id}/reject [post]
func (h *friendshipHandler) Reject(c echo.Context) error {
	f, err := h.service.Reject(c.Request().Context(), authmw.UserIDFromContext(c), c.Param("id"))
	if err != nil {
		resp, status := wrapper.FromError(err)
		return c.JSON(status, resp)
	}
	return wrapper.GenerateResponse(c, models.FriendshipResponse{Friendship: *f}, nil)
}

// ListPendingReceived godoc
// @Summary      List pending received friend requests
// @Tags         Friendships
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  models.FriendshipsListResponse
// @Router       /v1/friendships/requests/received [get]
func (h *friendshipHandler) ListPendingReceived(c echo.Context) error {
	list, err := h.service.ListPendingReceived(c.Request().Context(), authmw.UserIDFromContext(c))
	if err != nil {
		resp, status := wrapper.FromError(err)
		return c.JSON(status, resp)
	}
	return wrapper.GenerateResponse(c, models.FriendshipsListResponse{Friendships: list}, nil)
}

// ListPendingSent godoc
// @Summary      List pending sent friend requests
// @Tags         Friendships
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  models.FriendshipsListResponse
// @Router       /v1/friendships/requests/sent [get]
func (h *friendshipHandler) ListPendingSent(c echo.Context) error {
	list, err := h.service.ListPendingSent(c.Request().Context(), authmw.UserIDFromContext(c))
	if err != nil {
		resp, status := wrapper.FromError(err)
		return c.JSON(status, resp)
	}
	return wrapper.GenerateResponse(c, models.FriendshipsListResponse{Friendships: list}, nil)
}

// ListFriends godoc
// @Summary      List accepted friends
// @Tags         Friendships
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  models.FriendsListResponse
// @Router       /v1/friendships/friends [get]
func (h *friendshipHandler) ListFriends(c echo.Context) error {
	friends, err := h.service.ListFriends(c.Request().Context(), authmw.UserIDFromContext(c))
	if err != nil {
		resp, status := wrapper.FromError(err)
		return c.JSON(status, resp)
	}
	return wrapper.GenerateResponse(c, models.FriendsListResponse{Friends: friends}, nil)
}
