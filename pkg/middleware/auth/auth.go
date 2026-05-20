// Package auth provides JWT bearer authentication middleware for Echo.
package auth

import (
	"strings"

	kiterrors "github.com/donca/user-crud/pkg/kit/errors"
	"github.com/donca/user-crud/pkg/kit/wrapper"
	"github.com/donca/user-crud/pkg/kit/jwt"

	"github.com/labstack/echo/v4"
)

// UserIDKey is the Echo context key set after a valid JWT is parsed.
const UserIDKey = "user_id"

// Middleware validates Bearer tokens and exposes the authenticated user id on the context.
type Middleware struct {
	jwtSecret string
}

// NewMiddleware builds middleware that rejects missing or invalid Authorization headers.
func NewMiddleware(jwtSecret string) *Middleware {
	return &Middleware{jwtSecret: jwtSecret}
}

// RequireAuth aborts with 401 unless a valid JWT is present.
func (m *Middleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get(echo.HeaderAuthorization)
		if header == "" {
			resp, status := wrapper.FromError(kiterrors.Unauthorized("missing authorization header"))
			return c.JSON(status, resp)
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			resp, status := wrapper.FromError(kiterrors.Unauthorized("invalid authorization format"))
			return c.JSON(status, resp)
		}
		userID, err := jwt.ParseToken(parts[1], m.jwtSecret)
		if err != nil {
			resp, status := wrapper.FromError(kiterrors.Unauthorized("invalid or expired token"))
			return c.JSON(status, resp)
		}
		c.Set(UserIDKey, userID)
		return next(c)
	}
}

// UserIDFromContext returns the authenticated user id or an empty string when unauthenticated.
func UserIDFromContext(c echo.Context) string {
	id, _ := c.Get(UserIDKey).(string)
	return id
}
