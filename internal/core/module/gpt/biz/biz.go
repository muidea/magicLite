package biz

import (
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"
	"github.com/muidea/magicLite/internal/core/base/biz"
	"github.com/muidea/magicLite/pkg/common"
)

type GPT struct {
	biz.Base
}

func New(
	endpointName string,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *GPT {
	return &GPT{
		Base: biz.New(common.GPTModule, eventHub, backgroundRoutine),
	}
}
