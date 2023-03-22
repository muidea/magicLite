package model

import (
	"fmt"
)

const (
	// AccountType account
	AccountType = "account"
	// EndpointType endpoint
	EndpointType = "endpoint"
)

// Entity accessLog entity type
type Entity struct {
	ID        int    `json:"id" orm:"id key auto" tag:"view"`
	EName     string `json:"name" orm:"eName" tag:"view"`
	EID       int    `json:"eID" orm:"eID" tag:"view"`
	EType     string `json:"eType" orm:"eType" tag:"view"`
	Namespace string `json:"namespace" orm:"namespace"`
}

func (s *Entity) String() string {
	switch s.EType {
	case AccountType:
		return fmt.Sprintf("用户:%s", s.EName)
	case EndpointType:
		return fmt.Sprintf("终端:%s", s.EName)
	}

	return s.EName
}

func (s *Entity) IsAdmin() bool {
	superUser := DefaultSuperAccount()

	return s.EType == AccountType && s.EName == superUser.Account
}

func (s *Entity) Same(ptr *Entity) bool {
	return s.ID == ptr.ID
}

type OnlineEntity struct {
	ID          int     `json:"id" orm:"id key auto" tag:"view"`
	SessionID   string  `json:"sessionID" orm:"sessionID" tag:"view"`
	Entity      *Entity `json:"entity" orm:"entity" tag:"view"`
	RefreshTime int64   `json:"refreshTime" orm:"refreshTime" tag:"view"`
	ExpiryTime  int64   `json:"expiryTime" orm:"expiryTime" tag:"view"`
	Namespace   string  `json:"namespace" orm:"namespace"`
}
