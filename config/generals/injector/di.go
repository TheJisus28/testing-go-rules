package injector

import (
	"os"

	"github.com/donca/user-crud/config/generals/logger"
	"github.com/donca/user-crud/config/generals/router"
	authrouter "github.com/donca/user-crud/config/generals/router/auth"
	friendrouter "github.com/donca/user-crud/config/generals/router/friendships"
	"github.com/donca/user-crud/config/generals/router/health"
	postrouter "github.com/donca/user-crud/config/generals/router/posts"
	userrouter "github.com/donca/user-crud/config/generals/router/users"
	"github.com/donca/user-crud/config/storage"
	authhandler "github.com/donca/user-crud/internal/auth/handler"
	authservice "github.com/donca/user-crud/internal/auth/service"
	friendhandler "github.com/donca/user-crud/internal/friendships/handler"
	friendrepo "github.com/donca/user-crud/internal/friendships/repository"
	friendservice "github.com/donca/user-crud/internal/friendships/service"
	posthandler "github.com/donca/user-crud/internal/posts/handler"
	postrepo "github.com/donca/user-crud/internal/posts/repository"
	postservice "github.com/donca/user-crud/internal/posts/service"
	userhandler "github.com/donca/user-crud/internal/users/handler"
	userrepo "github.com/donca/user-crud/internal/users/repository"
	userservice "github.com/donca/user-crud/internal/users/service"
	"github.com/donca/user-crud/pkg/kit/enums"
	authmw "github.com/donca/user-crud/pkg/middleware/auth"

	"go.uber.org/dig"
)

var Container *dig.Container

func init() {
	Container = dig.New()
	check(Container.Provide(storage.NewPostgresPool))
	check(Container.Provide(userrepo.NewUserRepository))
	check(Container.Provide(friendrepo.NewFriendshipRepository))
	check(Container.Provide(postrepo.NewPostRepository))
	check(Container.Provide(userservice.NewUserService))
	check(Container.Provide(friendservice.NewFriendshipService))
	check(Container.Provide(postservice.NewPostService))
	check(Container.Provide(authservice.NewAuthService))
	check(Container.Provide(userhandler.NewUserHandler))
	check(Container.Provide(friendhandler.NewFriendshipHandler))
	check(Container.Provide(posthandler.NewPostHandler))
	check(Container.Provide(authhandler.NewAuthHandler))
	check(Container.Provide(provideAuthMiddleware))
	check(Container.Provide(health.NewHealthRouter))
	check(Container.Provide(authrouter.NewAuthRouter))
	check(Container.Provide(userrouter.NewUsersRouter))
	check(Container.Provide(postrouter.NewPostsRouter))
	check(Container.Provide(friendrouter.NewFriendshipsRouter))
	check(Container.Provide(router.NewRouter))
}

func provideAuthMiddleware() *authmw.Middleware {
	return authmw.NewMiddleware(os.Getenv(enums.JWTSecret))
}

func check(err error) {
	if err != nil {
		logger.Get().Fatal().Err(err).Msg("injector: provide failed")
	}
}
