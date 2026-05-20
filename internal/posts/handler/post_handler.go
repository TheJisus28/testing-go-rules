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

// Create godoc
// @Summary      Create a post
// @Tags         Posts
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      models.CreatePostRequest  true  "Post payload"
// @Success      200   {object}  models.PostResponse
// @Failure      400   {object}  wrapper.Response[any]
// @Router       /v1/posts [post]
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

// GetByID godoc
// @Summary      Get post by ID
// @Description  Respects public/private visibility for the viewer
// @Tags         Posts
// @Produce      json
// @Param        id   path      string  true  "Post ID"
// @Success      200  {object}  models.PostResponse
// @Failure      403  {object}  wrapper.Response[any]
// @Failure      404  {object}  wrapper.Response[any]
// @Router       /v1/posts/{id} [get]
func (h *postHandler) GetByID(c echo.Context) error {
	viewerID := authmw.UserIDFromContext(c)
	post, err := h.service.GetByID(c.Request().Context(), c.Param("id"), viewerID)
	if err != nil {
		resp, status := wrapper.FromError(err)
		return c.JSON(status, resp)
	}
	return wrapper.GenerateResponse(c, models.PostResponse{Post: *post}, nil)
}

// Update godoc
// @Summary      Update a post
// @Tags         Posts
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      string  true  "Post ID"
// @Param        body  body      models.UpdatePostRequest  true  "Update payload"
// @Success      200   {object}  models.PostResponse
// @Router       /v1/posts/{id} [put]
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

// Delete godoc
// @Summary      Delete a post
// @Tags         Posts
// @Security     BearerAuth
// @Param        id   path      string  true  "Post ID"
// @Success      200  {object}  wrapper.Response[any]
// @Router       /v1/posts/{id} [delete]
func (h *postHandler) Delete(c echo.Context) error {
	if err := h.service.Delete(c.Request().Context(), authmw.UserIDFromContext(c), c.Param("id")); err != nil {
		resp, status := wrapper.FromError(err)
		return c.JSON(status, resp)
	}
	return wrapper.GenerateResponse(c, map[string]string{"deleted": "true"}, nil)
}

// ListByUser godoc
// @Summary      List posts on a user wall
// @Tags         Posts
// @Produce      json
// @Param        userId  path      string  true  "User ID"
// @Param        limit   query     int     false  "Limit"
// @Param        offset  query     int     false  "Offset"
// @Success      200     {object}  models.PostsListResponse
// @Router       /v1/users/{userId}/posts [get]
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

// Feed godoc
// @Summary      Get personalized feed
// @Tags         Posts
// @Security     BearerAuth
// @Produce      json
// @Param        limit   query     int  false  "Limit"
// @Param        offset  query     int  false  "Offset"
// @Success      200     {object}  models.PostsListResponse
// @Router       /v1/posts/feed [get]
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
