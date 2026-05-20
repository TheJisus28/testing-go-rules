package router

import (
	"os"

	"github.com/donca/user-crud/config/generals/logger"
	authrouter "github.com/donca/user-crud/config/generals/router/auth"
	friendrouter "github.com/donca/user-crud/config/generals/router/friendships"
	"github.com/donca/user-crud/config/generals/router/health"
	postrouter "github.com/donca/user-crud/config/generals/router/posts"
	userrouter "github.com/donca/user-crud/config/generals/router/users"
	"github.com/donca/user-crud/pkg/kit/enums"
	"github.com/donca/user-crud/pkg/middleware/cors"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/donca/user-crud/docs"
)

type Router struct {
	echo *echo.Echo
}

func NewRouter(
	healthRouter *health.HealthRouter,
	authRouter *authrouter.AuthRouter,
	usersRouter *userrouter.UsersRouter,
	postsRouter *postrouter.PostsRouter,
	friendshipsRouter *friendrouter.FriendshipsRouter,
) *Router {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(cors.WithCors())
	e.Use(logger.Middleware())

	healthRouter.Register(e)
	authRouter.Register(e)
	usersRouter.Register(e)
	postsRouter.Register(e)
	friendshipsRouter.Register(e)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return &Router{echo: e}
}

func (r *Router) Start(addr string) error {
	if addr == "" {
		addr = ":" + os.Getenv(enums.AppPort)
	}
	logger.Get().Info().Str("addr", addr).Msg("router: starting server")
	return r.echo.Start(addr)
}
