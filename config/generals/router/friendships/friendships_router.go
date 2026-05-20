package friendships

import (
	friendhandler "github.com/donca/user-crud/internal/friendships/handler"
	authmw "github.com/donca/user-crud/pkg/middleware/auth"

	"github.com/labstack/echo/v4"
)

type FriendshipsRouter struct {
	handler friendhandler.FriendshipHandler
	auth    *authmw.Middleware
}

func NewFriendshipsRouter(h friendhandler.FriendshipHandler, auth *authmw.Middleware) *FriendshipsRouter {
	return &FriendshipsRouter{handler: h, auth: auth}
}

func (r *FriendshipsRouter) Register(e *echo.Echo) {
	g := e.Group("/v1/friendships", r.auth.RequireAuth)
	g.POST("/requests", r.handler.SendRequest)
	g.GET("/requests/received", r.handler.ListPendingReceived)
	g.GET("/requests/sent", r.handler.ListPendingSent)
	g.POST("/requests/:id/accept", r.handler.Accept)
	g.POST("/requests/:id/reject", r.handler.Reject)
	g.GET("/friends", r.handler.ListFriends)
}
