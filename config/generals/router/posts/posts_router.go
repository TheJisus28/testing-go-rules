package posts

import (
	posthandler "github.com/donca/user-crud/internal/posts/handler"
	authmw "github.com/donca/user-crud/pkg/middleware/auth"

	"github.com/labstack/echo/v4"
)

type PostsRouter struct {
	handler posthandler.PostHandler
	auth    *authmw.Middleware
}

func NewPostsRouter(h posthandler.PostHandler, auth *authmw.Middleware) *PostsRouter {
	return &PostsRouter{handler: h, auth: auth}
}

func (r *PostsRouter) Register(e *echo.Echo) {
	g := e.Group("/v1/posts")
	protected := g.Group("", r.auth.RequireAuth)
	protected.GET("/feed", r.handler.Feed)
	protected.POST("", r.handler.Create)
	protected.PUT("/:id", r.handler.Update)
	protected.DELETE("/:id", r.handler.Delete)
	g.GET("/:id", r.handler.GetByID)

	e.GET("/v1/users/:userId/posts", r.handler.ListByUser)
}
