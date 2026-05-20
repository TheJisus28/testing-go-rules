package handler

import (
	"strconv"

	"github.com/donca/user-crud/internal/posts/interfaces"
	"github.com/donca/user-crud/internal/posts/models"
	authmw "github.com/donca/user-crud/pkg/middleware/auth"
	kiterrors "github.com/donca/user-crud/pkg/kit/errors"
	"github.com/donca/user-crud/pkg/kit/wrapper"

	"github.com/labstack/echo/v4"
)

type PostHandler interface {
	Create(c echo.Context) error
	GetByID(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	ListByUser(c echo.Context) error
	Feed(c echo.Context) error
}

type postHandler struct {
	service interfaces.PostService
}

func NewPostHandler(service interfaces.PostService) PostHandler {
	return &postHandler{service: service}
}

func (h *postHandler) Create(c echo.Context) error {
	var req models.CreatePostRequest
	if err := c.Bind(&req); err != nil {
		resp, status := wrapper.FromError(kiterrors.Validation("invalid request body"))
		return c.JSON(status, resp)
	}
	post, err := h.service.Create(c.Request().Context(), authmw.UserIDFromContext(c), req)
	if err != nil {
		resp, status := wrapper.FromError(err)
		return c.JSON(status, resp)
	}
	return wrapper.GenerateResponse(c, models.PostResponse{Post: *post}, nil)
}

func (h *postHandler) GetByID(c echo.Context) error {
	viewerID := authmw.UserIDFromContext(c)
	post, err := h.service.GetByID(c.Request().Context(), c.Param("id"), viewerID)
	if err != nil {
		resp, status := wrapper.FromError(err)
		return c.JSON(status, resp)
	}
	return wrapper.GenerateResponse(c, models.PostResponse{Post: *post}, nil)
}

func (h *postHandler) Update(c echo.Context) error {
	var req models.UpdatePostRequest
	if err := c.Bind(&req); err != nil {
		resp, status := wrapper.FromError(kiterrors.Validation("invalid request body"))
		return c.JSON(status, resp)
	}
	post, err := h.service.Update(c.Request().Context(), authmw.UserIDFromContext(c), c.Param("id"), req)
	if err != nil {
		resp, status := wrapper.FromError(err)
		return c.JSON(status, resp)
	}
	return wrapper.GenerateResponse(c, models.PostResponse{Post: *post}, nil)
}

func (h *postHandler) Delete(c echo.Context) error {
	if err := h.service.Delete(c.Request().Context(), authmw.UserIDFromContext(c), c.Param("id")); err != nil {
		resp, status := wrapper.FromError(err)
		return c.JSON(status, resp)
	}
	return wrapper.GenerateResponse(c, map[string]string{"deleted": "true"}, nil)
}

func (h *postHandler) ListByUser(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	viewerID := authmw.UserIDFromContext(c)
	posts, err := h.service.ListByUser(c.Request().Context(), c.Param("userId"), viewerID, limit, offset)
	if err != nil {
		resp, status := wrapper.FromError(err)
		return c.JSON(status, resp)
	}
	return wrapper.GenerateResponse(c, models.PostsListResponse{Posts: posts}, nil)
}

func (h *postHandler) Feed(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	posts, err := h.service.Feed(c.Request().Context(), authmw.UserIDFromContext(c), limit, offset)
	if err != nil {
		resp, status := wrapper.FromError(err)
		return c.JSON(status, resp)
	}
	return wrapper.GenerateResponse(c, models.PostsListResponse{Posts: posts}, nil)
}
