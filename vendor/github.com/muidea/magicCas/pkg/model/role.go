package model

const SuperID = 999999

// Privilege 权限表
// value: 0 未设权限, 1 读权限, 2 写权限, 3 删除权限, 4 全部权限
// namePath: 权限路径
type Privilege struct {
	ID    int    `json:"id" orm:"id key auto"`
	Value int    `json:"value" orm:"value" validate:"required"`
	Path  string `json:"path" orm:"path" validate:"required"`
}

// Role privilege role
type Role struct {
	ID          int          `json:"id" orm:"id key auto" tag:"view,lite"`
	Name        string       `json:"name" orm:"name" validate:"required" tag:"view,lite,param"`
	Description string       `json:"description" orm:"description" tag:"view,param"`
	Privilege   []*Privilege `json:"privilege" orm:"privilege" tag:"view,param"`
	Namespace   string       `json:"namespace" orm:"namespace"`
}

func (s *Role) IsSuper() bool {
	return s.ID == SuperID
}
