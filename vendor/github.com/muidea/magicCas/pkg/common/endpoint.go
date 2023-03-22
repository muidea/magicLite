package common

import (
	cd "github.com/muidea/magicCommon/def"

	"github.com/muidea/magicCas/pkg/model"
)

const EndpointModule = "/module/endpoint"

const (
	FilterEndpoint = "/endpoint/query/"
	QueryEndpoint  = "/endpoint/query/:id"
	CreateEndpoint = "/endpoint/create/"
	UpdateEndpoint = "/endpoint/update/:id"
	DeleteEndpoint = "/endpoint/delete/:id"
	NotifyEndpoint = "/endpoint/notify/:id"
)

// EndpointView endpoint view
type EndpointView struct {
	ID          int          `json:"id"`
	Endpoint    string       `json:"endpoint"`
	AuthToken   string       `json:"authToken"`
	Description string       `json:"description"`
	Account     *AccountLite `json:"account"`
	Status      *Status      `json:"status"`
	CreateTime  int64        `json:"createTime"`
}

// FromEndpoint from endpoint
func (s *EndpointView) FromEndpoint(ptr *model.Endpoint) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Endpoint = ptr.Endpoint
	s.AuthToken = ptr.AuthToken
	s.Description = ptr.Description
	if ptr.Account != nil {
		s.Account = &AccountLite{}
		s.Account.FromAccount(ptr.Account)
	}
	s.Status = GetStatus(ptr.Status)

	s.CreateTime = ptr.CreateTime
}

func (s *EndpointView) ToEndpoint() (ret *model.Endpoint) {
	ret = &model.Endpoint{
		ID:          s.ID,
		Endpoint:    s.Endpoint,
		AuthToken:   s.AuthToken,
		Description: s.Description,
		Status:      s.Status.ID,
		CreateTime:  s.CreateTime,
	}

	if s.Account != nil {
		ret.Account = s.Account.ToAccount()
	}

	return
}

// Entity get endpoint entity
func (s *EndpointView) Entity() *model.Entity {
	return &model.Entity{EID: s.ID, EType: model.EndpointType, EName: s.Endpoint}
}

type EndpointLite struct {
	ID       int     `json:"id"`
	Endpoint string  `json:"endpoint"`
	Status   *Status `json:"status"`
}

func (s *EndpointLite) FromEndpoint(ptr *model.Endpoint) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Endpoint = ptr.Endpoint
	s.Status = GetStatus(ptr.Status)
}

func (s *EndpointLite) ToEndpoint() (ret *model.Endpoint) {
	ret = &model.Endpoint{
		ID:       s.ID,
		Endpoint: s.Endpoint,
		Status:   s.Status.ID,
	}

	return
}

type EndpointParam struct {
	Endpoint    string       `json:"endpoint"`
	Description string       `json:"description"`
	Account     *AccountLite `json:"account"`
	Status      *Status      `json:"status"`
}

func (s *EndpointParam) ToEndpoint() (ret *model.Endpoint) {
	ret = &model.Endpoint{
		Endpoint:    s.Endpoint,
		Description: s.Description,
	}
	if s.Account != nil {
		ret.Account = s.Account.ToAccount()
	}
	if s.Status != nil {
		ret.Status = s.Status.ID
	}
	return
}

type EndpointResult struct {
	cd.Result
	Endpoint *EndpointView `json:"endpoint"`
}

type EndpointLiteListResult struct {
	cd.Result
	Total    int64           `json:"total"`
	Endpoint []*EndpointLite `json:"endpoint"`
}

type EndpointListResult struct {
	cd.Result
	Total    int64           `json:"total"`
	Endpoint []*EndpointView `json:"endpoint"`
}

type EndpointStatisticResult struct {
	EndpointListResult
}
