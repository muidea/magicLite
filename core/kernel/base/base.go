package base

import (
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/session"
	"github.com/muidea/magicCommon/task"
	engine "github.com/muidea/magicEngine"

	"github.com/muidea/magicLite/common"
	"github.com/muidea/magicLite/core/kernel/base/biz"
	"github.com/muidea/magicLite/module"
)

func init() {
	module.Register(New())
}

type Base struct {
	sessionRegistry   session.Registry
	eventHub          event.Hub
	backgroundRoutine *task.BackgroundRoutine

	biz *biz.Base
}

func New() module.Module {
	return &Base{}
}

func (s *Base) ID() string {
	return common.BaseModule
}

func (s *Base) BindEventHub(eventHub event.Hub) {
	s.eventHub = eventHub
}

func (s *Base) BindBackgroundRoutine(backgroundRoutine *task.BackgroundRoutine) {
	s.backgroundRoutine = backgroundRoutine
}

func (s *Base) BindRegistry(sessionRegistry session.Registry) {
	s.sessionRegistry = sessionRegistry
}

func (s *Base) Startup(endpointName string, router engine.Router) {
	router.SetApiVersion(common.ApiVersion)

	s.biz = biz.New(s.eventHub, s.backgroundRoutine)
}

func (s *Base) Shutdown() {

}
