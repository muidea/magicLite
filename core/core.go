package core

import (
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/module"
	"github.com/muidea/magicCommon/task"
	engine "github.com/muidea/magicEngine"

	_ "github.com/muidea/magicLite/core/kernel/base"
)

// New 新建Core
func New(endpointName, listenPort string) (ret *Core, err error) {
	core := &Core{
		endpointName: endpointName,
		listenPort:   listenPort,
	}

	ret = core
	return
}

// Core Core对象
type Core struct {
	endpointName string
	listenPort   string

	routeRegistry engine.Router
	httpServer    engine.HTTPServer
}

// Startup 启动
func (s *Core) Startup(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine) {
	s.routeRegistry = engine.NewRouter()
	s.httpServer = engine.NewHTTPServer(s.listenPort)
	s.httpServer.Bind(s.routeRegistry)

	modules := module.GetModules()
	for _, val := range modules {

		module.BindRegistry(val, s.routeRegistry)

		module.Setup(val, s.endpointName, eventHub, backgroundRoutine)
	}
}

func (s *Core) Run() {
	s.httpServer.Run()
}

// Shutdown 销毁
func (s *Core) Shutdown() {
	modules := module.GetModules()
	for _, val := range modules {
		module.Teardown(val)
	}
}
