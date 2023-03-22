package toolkit

import (
	"context"
	"net/http"

	"github.com/muidea/magicCommon/session"
	engine "github.com/muidea/magicEngine"

	"github.com/muidea/magicCas/pkg/common"
)

// RouteRegistry route registry
type RouteRegistry interface {
	SetApiVersion(version string)

	AddHandler(pattern, method string, handler func(context.Context, http.ResponseWriter, *http.Request))

	AddRoute(route engine.Route, filters ...engine.MiddleWareHandler)
}

// NewRouteRegistry create routeRegistry
func NewRouteRegistry(registry session.Registry, router engine.Router) (ret RouteRegistry) {
	ret = &routeRegistryImpl{sessionRegistry: registry, router: router}
	return
}

// routeRegistryImpl route registry
type routeRegistryImpl struct {
	sessionRegistry session.Registry
	router          engine.Router
}

func (s *routeRegistryImpl) SetApiVersion(version string) {
	s.router.SetApiVersion(version)
}

// AddHandler add route handler
func (s *routeRegistryImpl) AddHandler(
	pattern, method string,
	handler func(context.Context, http.ResponseWriter, *http.Request)) {

	s.router.AddRoute(engine.CreateRoute(pattern, method, handler), s)
}

func (s *routeRegistryImpl) AddRoute(route engine.Route, filters ...engine.MiddleWareHandler) {
	filters = append(filters, s)
	s.router.AddRoute(route, filters...)
}

func (s *routeRegistryImpl) Handle(ctx engine.RequestContext, res http.ResponseWriter, req *http.Request) {
	curSession := s.sessionRegistry.GetSession(res, req)
	ctx.Update(context.WithValue(ctx.Context(), common.AuthSession, curSession))
	ctx.Next()
}
