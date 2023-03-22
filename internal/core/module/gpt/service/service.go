package service

import (
	"context"
	"encoding/json"
	"github.com/muidea/magicCas/pkg/toolkit"
	"net/http"

	log "github.com/cihub/seelog"

	commonDef "github.com/muidea/magicCommon/def"
	fn "github.com/muidea/magicCommon/foundation/net"
	fu "github.com/muidea/magicCommon/foundation/util"

	engine "github.com/muidea/magicEngine"

	"github.com/muidea/magicLite/pkg/common"

	"github.com/muidea/magicLite/internal/core/module/gpt/biz"
)

type GPT struct {
	routeRegistry toolkit.RouteRegistry
	validator     fu.Validator

	bizPtr *biz.GPT
}

func New(bizPtr *biz.GPT) *GPT {
	ptr := &GPT{
		bizPtr:    bizPtr,
		validator: fu.NewFormValidator(),
	}

	return ptr
}

func (s *GPT) BindRegistry(
	routeRegistry toolkit.RouteRegistry) {

	s.routeRegistry = routeRegistry
}

func (s *GPT) RegisterRoute() {
	queryRoute := engine.CreateRoute(common.QueryMessage, "POST", s.queryMessage)
	s.routeRegistry.AddRoute(queryRoute)
}

func (s *GPT) queryMessage(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	param := &common.QueryParam{}
	result := &common.QueryResult{}
	for {
		err := fn.ParseJSONBody(req, nil, param)
		if err != nil {
			log.Errorf("ParseJSONBody failed, err:%s", err.Error())
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法请求"
			break
		}

		msgData, msgErr := s.bizPtr.QueryMessage(param.Message)
		if msgErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = msgErr.Error()
			break
		}

		result.Data = msgData
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)

}
