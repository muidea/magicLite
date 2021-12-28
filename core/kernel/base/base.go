package base

import (
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/module"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCas/toolkit"

	"github.com/muidea/magicLite/common"
	"github.com/muidea/magicLite/core/kernel/base/biz"
)

func init() {
	module.Register(New())
}

type Base struct {
	routeRegistry toolkit.RouteRegistry

	biz *biz.Base
}

func New() module.Module {
	return &Base{}
}

func (s *Base) ID() string {
	return common.BaseModule
}

func (s *Base) BindRegistry(routeRegistry toolkit.RouteRegistry) {
	s.routeRegistry = routeRegistry
}

func (s *Base) Setup(
	endpointName string,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine) {
	s.biz = biz.New(
		eventHub,
		backgroundRoutine,
	)
}

func (s *Base) Teardown() {

}
