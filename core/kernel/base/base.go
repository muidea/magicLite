package base

import (
	"github.com/muidea/magicCommon/application"
	"github.com/muidea/magicCommon/module"
	"github.com/muidea/magicLite/common"
	"github.com/muidea/magicLite/core/kernel/base/biz"
)

func init() {
	module.Register(New())
}

type Base struct {
	biz *biz.Base
}

func New() module.Module {
	return &Base{}
}

func (s *Base) ID() string {
	return common.BaseModule
}

func (s *Base) Setup(endpointName string) {
	app := application.GetApp()
	s.biz = biz.New(
		app.EventHub(),
		app.BackgroundRoutine(),
	)
}

func (s *Base) Teardown() {

}
