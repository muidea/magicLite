package biz

import (
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicLite/pkg/common"

	"github.com/muidea/magicLite/internal/core/base/biz"
)

type Base struct {
	biz.Base
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Base {
	ptr := &Base{
		Base: biz.New(common.BaseModule, eventHub, backgroundRoutine),
	}

	return ptr
}

func (s *Base) Notify(event event.Event, result event.Result) {
}
