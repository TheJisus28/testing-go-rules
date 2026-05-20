package auth

import (
	authhandler "github.com/donca/user-crud/internal/auth/handler"

	"github.com/labstack/echo/v4"
)

type AuthRouter struct {
	handler authhandler.AuthHandler
}

func NewAuthRouter(h authhandler.AuthHandler) *AuthRouter {
	return &AuthRouter{handler: h}
}

func (r *AuthRouter) Register(e *echo.Echo) {
	g := e.Group("/v1/auth")
	g.POST("/register", r.handler.Register)
	g.POST("/login", r.handler.Login)
}
