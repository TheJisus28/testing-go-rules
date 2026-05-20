// @title       SocialNet API
// @version     1.0
// @description Mini red social con usuarios, posts públicos/privados y solicitudes de amistad
// @BasePath    /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Token JWT. Formato: Bearer &lt;token&gt;
package main

import (
	"os"

	"github.com/donca/user-crud/config"
	"github.com/donca/user-crud/config/generals/injector"
	"github.com/donca/user-crud/config/generals/logger"
	"github.com/donca/user-crud/config/generals/router"
	"github.com/donca/user-crud/pkg/kit/enums"
)

func main() {
	config.LoadEnv()
	logger.Init(os.Getenv(enums.LogLevel))
	log := logger.Get()

	if err := injector.Container.Invoke(func(r *router.Router) {
		addr := ":" + os.Getenv(enums.AppPort)
		log.Fatal().Err(r.Start(addr)).Msg("server stopped")
	}); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
