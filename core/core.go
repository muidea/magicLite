package core

import (
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/session"
	"github.com/muidea/magicCommon/task"
	engine "github.com/muidea/magicEngine"

	_ "github.com/muidea/magicLite/core/kernel/base"
	"github.com/muidea/magicLite/module"
)

// New 新建Core
func New(endpointName string) (ret *Core, err error) {
	sessionRegistry := session.CreateRegistry(nil)
	eventHub := event.NewHub()
	backgroundRoutine := task.NewBackgroundRoutine()

	core := &Core{
		endpointName:      endpointName,
		sessionRegistry:   sessionRegistry,
		eventHub:          eventHub,
		backgroundRoutine: backgroundRoutine,
	}

	ret = core
	return
}

// Core Core对象
type Core struct {
	endpointName      string
	sessionRegistry   session.Registry
	eventHub          event.Hub
	backgroundRoutine *task.BackgroundRoutine
}

// Startup 启动
func (s *Core) Startup(router engine.Router) {
	moduleList := module.GetList()
	for _, val := range moduleList {
		val.BindBackgroundRoutine(s.backgroundRoutine)
		val.BindEventHub(s.eventHub)
		val.BindRegistry(s.sessionRegistry)
	}

	for _, val := range moduleList {
		val.Startup(s.endpointName, router)
	}
}

// Teardown 销毁
func (s *Core) Teardown() {
}
