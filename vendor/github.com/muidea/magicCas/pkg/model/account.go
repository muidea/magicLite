package model

// Account account
type Account struct {
	ID          int    `json:"id" orm:"id key auto" tag:"view,lite"`
	Account     string `json:"account" orm:"account" tag:"view,lite,param"`
	Password    string `json:"password" orm:"password" tag:"param"`
	EMail       string `json:"email" orm:"email" tag:"view,param"`
	Description string `json:"description" orm:"description" tag:"view,param"`
	Status      int    `json:"status" orm:"status" tag:"view,lite,param"`
	Role        *Role  `json:"role" orm:"role" tag:"view,param"`
	CreateTime  int64  `json:"createTime" orm:"createTime" tag:"view"`
	Namespace   string `json:"namespace" orm:"namespace"`
}

// Enable enable status
func (s *Account) Enable() bool {
	return s.Status == EnableStatus
}

// Disable disable status
func (s *Account) Disable() bool {
	return s.Status == DisableStatus
}

// Entity get account entity
func (s *Account) Entity() *Entity {
	return &Entity{EName: s.Account, EID: s.ID, EType: AccountType, Namespace: s.Namespace}
}

//IsSame is same account
func (s *Account) IsSame(account *Account) bool {
	if s.ID != account.ID {
		return false
	}

	if s.Account != account.Account {
		return false
	}

	if s.Password != account.Password {
		return false
	}

	if s.EMail != account.EMail {
		return false
	}

	if s.Description != account.Description {
		return false
	}

	if s.Status != account.Status {
		return false
	}

	if s.Role != nil && account.Role != nil {
		if s.Role.ID != account.Role.ID {
			return false
		}
	}
	if s.Role != nil && account.Role == nil {
		return false
	}
	if s.Role == nil && account.Role != nil {
		return false
	}

	return true
}

func (s *Account) SupperAccount() bool {
	return s.Account == "administrator" && s.Status == InitStatus && s.Role == nil
}

func DefaultSuperAccount() *Account {
	return &Account{
		Account:     "administrator",
		Password:    "administrator",
		EMail:       "rangh@foxmail.com",
		Description: "default administrator",
	}
}
