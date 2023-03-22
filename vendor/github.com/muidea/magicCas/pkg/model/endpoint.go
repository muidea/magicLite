package model

// Endpoint endpoint
type Endpoint struct {
	ID          int      `json:"id" orm:"id key auto" tag:"view,lite"`
	Endpoint    string   `json:"endpoint" orm:"endpoint" tag:"view,lite,param"`
	AuthToken   string   `json:"authToken" orm:"authToken" tag:"view,param"`
	Description string   `json:"description" orm:"description" tag:"view,param"`
	Status      int      `json:"status" orm:"status" tag:"view,lite,param"`
	Account     *Account `json:"account" orm:"account" tag:"view,param"`
	CreateTime  int64    `json:"createTime" orm:"createTime" tag:"view"`
	Namespace   string   `json:"namespace" orm:"namespace"`
}

// Enable enable status
func (s *Endpoint) Enable() bool {
	return s.Status == EnableStatus
}

// Disable disable status
func (s *Endpoint) Disable() bool {
	return s.Status == DisableStatus
}

// Entity get endpoint entity
func (s *Endpoint) Entity() *Entity {
	return &Entity{EName: s.Endpoint, EID: s.ID, EType: EndpointType}
}

//IsSame is same endpoint
func (s *Endpoint) IsSame(endpoint *Endpoint) bool {
	if s.ID != endpoint.ID {
		return false
	}

	if s.Endpoint != endpoint.Endpoint {
		return false
	}

	if s.AuthToken != endpoint.AuthToken {
		return false
	}

	if s.Description != endpoint.Description {
		return false
	}

	if s.Status != endpoint.Status {
		return false
	}

	if s.Account == nil || endpoint.Account != nil {
		return false
	}
	if s.Account != nil || endpoint.Account == nil {
		return false
	}
	if s.Account.ID != endpoint.Account.ID {
		return false
	}

	return true
}

func DefaultSuperEndpoint() *Endpoint {
	return &Endpoint{
		Endpoint:    "defaultEndpoint",
		Description: "default description",
	}
}
