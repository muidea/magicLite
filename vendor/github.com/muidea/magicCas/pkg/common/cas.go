package common

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/session"
)

const CasModule = "/kernel/cas"

const (
	// LoginAccount account login
	LoginAccount = "/account/login/"
	// LogoutAccount account logout
	LogoutAccount = "/account/logout/"
	// VerifySession verify session
	VerifySession = "/session/verify/"
	// UpdateAccountPassword update account password
	UpdateAccountPassword = "/account/password/update/"
	// VerifyAccount verify account
	VerifyAccount = "/account/verify/"
)

type VerifyResult struct {
	cd.Result
	Entity       *EntityView   `json:"entity"`
	SessionToken session.Token `json:"sessionToken"`
}

type LoginParam struct {
	Account  string `json:"account" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResult VerifyResult

type LogoutResult struct {
	cd.Result
}

type UpdatePasswordParam struct {
	Account     string `json:"account" validate:"required"`
	CurPassword string `json:"curPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}

type UpdatePasswordResult cd.Result
