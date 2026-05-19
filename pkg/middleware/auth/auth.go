package auth

import (
	"strings"

	kiterrors "github.com/donca/user-crud/pkg/kit/errors"
	"github.com/donca/user-crud/pkg/kit/wrapper"
	"github.com/donca/user-crud/pkg/kit/jwt"

	"github.com/labstack/echo/v4"
)

const UserIDKey = "user_id"

type Middleware struct {
	jwtSecret string
}

func NewMiddleware(jwtSecret string) *Middleware {
	return &Middleware{jwtSecret: jwtSecret}
}

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

func UserIDFromContext(c echo.Context) string {
	id, _ := c.Get(UserIDKey).(string)
	return id
}
