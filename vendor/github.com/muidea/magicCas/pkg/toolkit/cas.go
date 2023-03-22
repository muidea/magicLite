package toolkit

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/session"

	engine "github.com/muidea/magicEngine"

	cClnt "github.com/muidea/magicCas/pkg/client"
	"github.com/muidea/magicCas/pkg/common"
)

// CasRegistry cas route registry
type CasRegistry interface {
	SetApiVersion(version string)

	AddHandler(pattern, method string, handler func(context.Context, http.ResponseWriter, *http.Request))

	AddRoute(route engine.Route, filters ...engine.MiddleWareHandler)
}

// NewCasRegistry create cas Registry
func NewCasRegistry(casService string, registry session.Registry, router engine.Router) (ret CasRegistry) {
	ret = &casRegistryImpl{casService: casService, sessionRegistry: registry, router: router}
	return
}

// casRegistryImpl cas route registry
type casRegistryImpl struct {
	casService      string
	sessionRegistry session.Registry
	router          engine.Router
}

func (s *casRegistryImpl) SetApiVersion(version string) {
	s.router.SetApiVersion(version)
}

// AddHandler add route handler
func (s *casRegistryImpl) AddHandler(
	pattern, method string,
	handler func(context.Context, http.ResponseWriter, *http.Request)) {

	s.router.AddRoute(engine.CreateRoute(pattern, method, handler), s)
}

func (s *casRegistryImpl) AddRoute(route engine.Route, filters ...engine.MiddleWareHandler) {
	filters = append(filters, s)
	s.router.AddRoute(route, filters...)
}

// Handle middleware handler
func (s *casRegistryImpl) Handle(ctx engine.RequestContext, res http.ResponseWriter, req *http.Request) {
	var curEntity *common.EntityView
	curSession := s.sessionRegistry.GetSession(res, req)
	authEntity, authOK := curSession.GetInt(common.AuthEntity)
	if authOK {
		func() {
			clnt := cClnt.NewClient(s.casService)
			defer clnt.Release()

			curNamespace, _ := curSession.GetString(common.AuthNamespace)
			curAuthorization, _ := curSession.GetString(session.Authorization)
			clnt.AttachAuthorization(curAuthorization)
			clnt.AttachNameSpace(curNamespace)

			entityPtr, entityErr := clnt.QueryEntity(int(authEntity))
			if entityErr != def.Success {
				return
			}

			curSession.SetOption(common.AuthEntityView, entityPtr)

			curEntity = entityPtr
		}()
	}

	if curEntity == nil {
		result := &def.Result{ErrorCode: def.InvalidAuthority, Reason: "请先登录系统"}
		block, err := json.Marshal(result)
		if err == nil {
			res.Write(block)
			return
		}

		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Update(context.WithValue(ctx.Context(), common.AuthSession, curSession))
	ctx.Next()
}
