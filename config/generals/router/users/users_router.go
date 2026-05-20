package users

import (
	userhandler "github.com/donca/user-crud/internal/users/handler"
	authmw "github.com/donca/user-crud/pkg/middleware/auth"

	"github.com/labstack/echo/v4"
)

type UsersRouter struct {
	handler userhandler.UserHandler
	auth    *authmw.Middleware
}

func NewUsersRouter(h userhandler.UserHandler, auth *authmw.Middleware) *UsersRouter {
	return &UsersRouter{handler: h, auth: auth}
}

func (r *UsersRouter) Register(e *echo.Echo) {
	g := e.Group("/v1/users")
	g.GET("", r.handler.List)
	g.GET("/:id/profile", r.handler.GetProfile)
	g.GET("/:id", r.handler.GetByID)

	protected := g.Group("", r.auth.RequireAuth)
	protected.PUT("/:id", r.handler.Update)
	protected.DELETE("/:id", r.handler.Delete)

	me := e.Group("/v1/me", r.auth.RequireAuth)
	me.PATCH("/profile", r.handler.UpdateProfile)
}
