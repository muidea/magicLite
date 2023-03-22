package common

import (
	"time"

	cd "github.com/muidea/magicCommon/def"

	"github.com/muidea/magicCas/pkg/model"
)

const BaseModule = "/kernel/base"

const (
	// SaveEntity save entity
	SaveEntity = "/entity/save/"
	// FilterEntity filter entity
	FilterEntity = "/entity/query/"
	// QueryEntity query entity
	QueryEntity = "/entity/query/:id"
	// DeleteEntity delete entity
	DeleteEntity = "/entity/delete/"
	// ClearEntity clear entity
	ClearEntity = "/entity/clear/"
	// QueryEntityRole query entity role
	QueryEntityRole = "/entity/role/:id"

	// QueryAccessLog query access log
	QueryAccessLog = "/access/log/query/"
	// WriteAccessLog write access log
	WriteAccessLog = "/access/log/write/"

	// NotifyTimer notify timer
	NotifyTimer = "/base/timer/notify/"
)

type EntityView struct {
	ID    int    `json:"id"`
	EName string `json:"name"`
	EID   int    `json:"eID"`
	EType string `json:"eType"`
}

func (s *EntityView) FromEntity(ptr *model.Entity) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.EID = ptr.EID
	s.EName = ptr.EName
	s.EType = ptr.EType
}

func (s *EntityView) ToEntity() (ret *model.Entity) {
	ret = &model.Entity{
		ID:    s.ID,
		EName: s.EName,
		EID:   s.EID,
		EType: s.EType,
	}
	return
}

// IsValid check valid
func (s *EntityView) IsValid() bool {
	if s.EID <= 0 {
		return false
	}

	switch s.EType {
	case model.AccountType, model.EndpointType:
		return true
	}

	return false
}

func (s *EntityView) IsAdmin() bool {
	superUser := model.DefaultSuperAccount()

	return s.EType == model.AccountType && s.EName == superUser.Account
}

func (s *EntityView) IsSupper() bool {
	if s.EType == model.EndpointType {
		return true
	}

	return s.IsAdmin()
}

func (s *EntityView) IsSame(ptr *model.Entity) bool {
	if ptr == nil {
		return false
	}

	return s.EID == ptr.EID && s.EType == ptr.EType && s.EName == ptr.EName
}

// UnknownEntity 未知账号
func UnknownEntity() *EntityView {
	return &EntityView{ID: -1, EID: -1, EName: "未知账号"}
}

type EntityResult struct {
	cd.Result
	Entity *EntityView `json:"entity"`
}

type EntityListResult struct {
	cd.Result
	Entity []*EntityView `json:"entity"`
}

type EntityRoleResult struct {
	cd.Result
	Role *RoleView `json:"role"`
}

type LogView struct {
	ID         int         `json:"id"`
	Address    string      `json:"address"`
	Memo       string      `json:"memo"`
	Creater    *EntityView `json:"creater"`
	CreateTime int64       `json:"createTime"`
}

func (s *LogView) FromLog(ptr *model.Log, createrPtr *model.Entity) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Address = ptr.Address
	s.Memo = ptr.Memo
	s.Creater = &EntityView{}
	s.Creater.FromEntity(createrPtr)
	s.CreateTime = ptr.CreateTime
}

type AccessLogResult struct {
	cd.Result
	Total     int64      `json:"total"`
	AccessLog []*LogView `json:"accessLog"`
}

const (
	// InvalidPermission 无效权限
	InvalidPermission = iota
	// UnSetPermission 未设权限
	UnSetPermission
	// ReadPermission 只读权限
	ReadPermission
	// WritePermission 可写权限
	WritePermission
	// DeletePermission 删除权限
	DeletePermission
	// AllPermission 全部权限
	AllPermission
)

// Permission info
type Permission struct {
	Value int    `json:"value"`
	Name  string `json:"name"`
}

var permissionList = []*Permission{
	{Value: InvalidPermission, Name: "无效权限"},
	{Value: UnSetPermission, Name: "未设权限"},
	{Value: ReadPermission, Name: "只读权限"},
	{Value: WritePermission, Name: "可写权限"},
	{Value: DeletePermission, Name: "删除权限"},
	{Value: AllPermission, Name: "全部权限"},
}

// GetPermissionList get permission list
func GetPermissionList() []*Permission {
	return permissionList
}

// GetPermission get permission
func GetPermission(value int) (ret *Permission) {
	switch value {
	case UnSetPermission,
		ReadPermission,
		WritePermission,
		DeletePermission,
		AllPermission:
		ret = permissionList[value]
	default:
		ret = permissionList[0]
	}

	return
}

// EnumPrivilegeItemResult enum privilege item result
type EnumPrivilegeItemResult struct {
	cd.Result
	Privilege []*Privilege `json:"privilege"`
}

type TimerNotify struct {
	PreTime time.Time
	CurTime time.Time
}
