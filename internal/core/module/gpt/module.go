package gpt

import (
	"github.com/muidea/magicCas/pkg/toolkit"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/module"
	"github.com/muidea/magicCommon/task"
	"github.com/muidea/magicLite/internal/core/module/gpt/biz"
	"github.com/muidea/magicLite/internal/core/module/gpt/service"
	"github.com/muidea/magicLite/pkg/common"
)

func init() {
	module.Register(New())
}

type GPT struct {
	routeRegistry toolkit.RouteRegistry

	service *service.GPT
	biz     *biz.GPT
}

func New() *GPT {
	return &GPT{}
}

func (s *GPT) ID() string {
	return common.GPTModule
}

func (s *GPT) BindRegistry(routeRegistry toolkit.RouteRegistry) {
	s.routeRegistry = routeRegistry
	s.routeRegistry.SetApiVersion(common.ApiVersion)
}

func (s *GPT) Setup(endpointName string, eventHub event.Hub, backgroundRoutine task.BackgroundRoutine) {
	s.biz = biz.New(endpointName, eventHub, backgroundRoutine)
	if s.biz == nil {
		return
	}

	s.service = service.New(s.biz)
	s.service.BindRegistry(s.routeRegistry)
	s.service.RegisterRoute()
}

func (s *GPT) Teardown() {

}
