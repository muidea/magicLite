package client

import (
	"fmt"
	"net/url"
	"strings"

	log "github.com/cihub/seelog"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/net"
	"github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicCommon/session"

	"github.com/muidea/magicCas/pkg/common"
)

// Client interface
type Client interface {
	session.Client

	AttachNameSpace(namespace string)

	VerifySession() (*common.EntityView, session.Token, cd.ErrorCode)
	LoginAccount(account, password string) (*common.EntityView, session.Token, cd.ErrorCode)
	LogoutAccount() cd.ErrorCode
	UpdateAccountPassword(ptr *common.UpdatePasswordParam) (*common.AccountView, cd.ErrorCode)
	VerifyAccount(account, password string) (*common.EntityView, cd.ErrorCode)
	FilterEntity(filter *util.ContentFilter) ([]*common.EntityView, cd.ErrorCode)
	QueryEntity(id int) (*common.EntityView, cd.ErrorCode)
	QueryEntityRole(id int) (*common.RoleView, cd.ErrorCode)
	FilterAccessLog(ptr *common.EntityView, filter *util.Pagination) ([]*common.LogView, int64, cd.ErrorCode)

	FilterAccount(filter *util.ContentFilter) ([]*common.AccountView, int64, cd.ErrorCode)
	FilterAccountLite(filter *util.ContentFilter) ([]*common.AccountLite, int64, cd.ErrorCode)
	QueryAccount(id int) (*common.AccountView, cd.ErrorCode)
	CreateAccount(ptr *common.AccountParam) (*common.AccountView, cd.ErrorCode)
	UpdateAccount(id int, ptr *common.AccountParam) (*common.AccountView, cd.ErrorCode)
	DeleteAccount(id int) (*common.AccountView, cd.ErrorCode)
	CheckAccount(account string) ([]*common.AccountLite, cd.ErrorCode)

	FilterEndpoint(filter *util.ContentFilter) ([]*common.EndpointView, int64, cd.ErrorCode)
	FilterEndpointLite(filter *util.ContentFilter) ([]*common.EndpointLite, int64, cd.ErrorCode)
	QueryEndpoint(id int) (*common.EndpointView, cd.ErrorCode)
	CreateEndpoint(ptr *common.EndpointParam) (*common.EndpointView, cd.ErrorCode)
	UpdateEndpoint(id int, ptr *common.EndpointParam) (*common.EndpointView, cd.ErrorCode)
	DeleteEndpoint(id int) (*common.EndpointView, cd.ErrorCode)

	FilterRole(filter *util.ContentFilter) ([]*common.RoleView, int64, cd.ErrorCode)
	FilterRoleLite(filter *util.ContentFilter) ([]*common.RoleLite, int64, cd.ErrorCode)
	QueryRole(id int) (*common.RoleView, cd.ErrorCode)
	CreateRole(ptr *common.RoleParam) (*common.RoleView, cd.ErrorCode)
	UpdateRole(id int, ptr *common.RoleParam) (*common.RoleView, cd.ErrorCode)
	DeleteRole(id int) (*common.RoleView, cd.ErrorCode)

	FilterNamespace(filter *util.ContentFilter) ([]*common.NamespaceView, int64, cd.ErrorCode)
	FilterNamespaceLite(filter *util.ContentFilter) ([]*common.NamespaceLite, int64, cd.ErrorCode)
	QueryNamespace(id int) (*common.NamespaceView, cd.ErrorCode)
	CreateNamespace(ptr *common.NamespaceParam) (*common.NamespaceView, cd.ErrorCode)
	UpdateNamespace(id int, ptr *common.NamespaceParam) (*common.NamespaceView, cd.ErrorCode)
	DeleteNamespace(id int) (*common.NamespaceView, cd.ErrorCode)
}

// NewClient new client
func NewClient(serverURL string) Client {
	clnt := &client{BaseClient: session.NewBaseClient(serverURL)}

	return clnt
}

type client struct {
	session.BaseClient

	namespace string
}

func (s *client) AttachNameSpace(namespace string) {
	s.namespace = namespace
}

func (s *client) GetContextValues() url.Values {
	vals := s.BaseClient.GetContextValues()
	if s.namespace != "" {
		vals.Set(common.AuthNamespace, s.namespace)
	}

	return vals
}

func (s *client) VerifySession() (*common.EntityView, session.Token, cd.ErrorCode) {
	result := &common.VerifyResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.VerifySession}, "")
	url.RawQuery = vals.Encode()
	_, err := net.HTTPGet(s.BaseClient.GetHTTPClient(), url.String(), result, s.GetContextValues())
	if err != nil {
		log.Errorf("verify session failed, err:%s", err.Error())
		return nil, "", cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("刷新会话状态失败，%s", result.Reason)
		log.Errorf("verify session failed, err:%s", err.Error())
		return nil, "", result.ErrorCode
	}

	return result.Entity, result.SessionToken, result.ErrorCode
}

func (s *client) LoginAccount(account, password string) (*common.EntityView, session.Token, cd.ErrorCode) {
	result := &common.LoginResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.LoginAccount}, "")
	url.RawQuery = vals.Encode()
	param := &common.LoginParam{Account: account, Password: password}
	_, err := net.HTTPPost(s.BaseClient.GetHTTPClient(), url.String(), param, result, s.GetContextValues())
	if err != nil {
		log.Errorf("login account failed, err:%s", err.Error())
		return nil, "", cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("登录失败，%s", result.Reason)
		log.Errorf("login account failed, err:%s", err.Error())
		return nil, "", result.ErrorCode
	}

	return result.Entity, result.SessionToken, result.ErrorCode
}

func (s *client) LogoutAccount() cd.ErrorCode {
	result := &common.LogoutResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.LogoutAccount}, "")
	url.RawQuery = vals.Encode()
	_, err := net.HTTPDelete(s.BaseClient.GetHTTPClient(), url.String(), result, s.GetContextValues())
	if err != nil {
		log.Errorf("logout account failed, err:%s", err.Error())
		return cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("登出失败，%s", result.Reason)
		log.Errorf("logout account failed, err:%s", err.Error())
		return result.ErrorCode
	}

	return result.ErrorCode
}

func (s *client) UpdateAccountPassword(ptr *common.UpdatePasswordParam) (*common.AccountView, cd.ErrorCode) {
	result := &common.AccountResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.UpdateAccountPassword}, "")
	url.RawQuery = vals.Encode()
	_, err := net.HTTPPut(s.BaseClient.GetHTTPClient(), url.String(), ptr, result, s.GetContextValues())
	if err != nil {
		log.Errorf("update account password failed, err:%s", err.Error())
		return nil, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("修改账号密码失败，%s", result.Reason)
		log.Errorf("update account password failed, err:%s", err.Error())
		return nil, result.ErrorCode
	}

	return result.Account, result.ErrorCode
}

func (s *client) VerifyAccount(account, password string) (*common.EntityView, cd.ErrorCode) {
	result := &common.EntityResult{}

	vals := url.Values{}
	filter := util.NewFilter()
	filter.Set("account", account)
	filter.Set("password", password)
	vals = filter.Encode(vals)
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.VerifyAccount}, "")
	url.RawQuery = vals.Encode()
	_, err := net.HTTPGet(s.BaseClient.GetHTTPClient(), url.String(), result, s.GetContextValues())
	if err != nil {
		log.Errorf("verify account failed, err:%s", err.Error())
		return nil, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("验证Account失败，%s", result.Reason)
		log.Errorf("verify account failed, err:%s", err.Error())
		return nil, result.ErrorCode
	}

	return result.Entity, result.ErrorCode
}

func (s *client) FilterEntity(filter *util.ContentFilter) ([]*common.EntityView, cd.ErrorCode) {
	result := &common.EntityListResult{}

	vals := url.Values{}
	if filter != nil {
		vals = filter.Encode(vals)
	}
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.FilterEntity}, "")
	url.RawQuery = vals.Encode()

	_, err := net.HTTPGet(s.BaseClient.GetHTTPClient(), url.String(), result, s.GetContextValues())
	if err != nil {
		log.Errorf("filter entity failed, err:%s", err.Error())
		return nil, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("查询Entity列表失败，%s", result.Reason)
		log.Errorf("filter entity failed, err:%s", err.Error())
		return nil, result.ErrorCode
	}

	return result.Entity, result.ErrorCode
}

func (s *client) QueryEntity(id int) (*common.EntityView, cd.ErrorCode) {
	result := &common.EntityResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.QueryEntity}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err := net.HTTPGet(s.BaseClient.GetHTTPClient(), url.String(), result, s.GetContextValues())
	if err != nil {
		log.Errorf("query entity failed, id:%d, err:%s", id, err.Error())
		return nil, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("查询指定Entity失败，%s", result.Reason)
		log.Errorf("query entity failed, id:%d, err:%s", id, err.Error())
		return nil, result.ErrorCode
	}

	return result.Entity, result.ErrorCode
}

func (s *client) QueryEntityRole(id int) (*common.RoleView, cd.ErrorCode) {
	result := &common.EntityRoleResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.QueryEntityRole}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()
	_, err := net.HTTPGet(s.BaseClient.GetHTTPClient(), url.String(), result, s.GetContextValues())
	if err != nil {
		log.Errorf("query entity role failed, id:%d, err:%s", id, err.Error())
		return nil, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("查询指定Entity权限组失败，%s", result.Reason)
		log.Errorf("query entity role failed, id:%d, err:%s", id, err.Error())
		return nil, result.ErrorCode
	}

	return result.Role, result.ErrorCode
}

func (s *client) FilterAccessLog(entityPtr *common.EntityView, paginationPtr *util.Pagination) ([]*common.LogView, int64, cd.ErrorCode) {
	result := &common.AccessLogResult{}

	filter := util.NewFilter()
	filter.Equal("creater", entityPtr.ID)
	filter.Pagination = paginationPtr
	vals := url.Values{}
	vals = filter.Encode(vals)

	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.QueryAccessLog}, "")
	url.RawQuery = vals.Encode()

	_, err := net.HTTPGet(s.BaseClient.GetHTTPClient(), url.String(), result, s.GetContextValues())
	if err != nil {
		log.Errorf("filter access log failed, err:%s", err.Error())
		return nil, 0, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("查询日志失败，%s", result.Reason)
		log.Errorf("filter access log failed, err:%s", err.Error())
		return nil, 0, result.ErrorCode
	}

	return result.AccessLog, result.Total, result.ErrorCode
}

func (s *client) FilterAccount(filter *util.ContentFilter) ([]*common.AccountView, int64, cd.ErrorCode) {
	result := &common.AccountListResult{}

	vals := url.Values{}
	if filter == nil {
		filter = util.NewFilter()
	}
	filter.Set("mode", cd.ViewMode)
	vals = filter.Encode(vals)

	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.FilterAccount}, "")
	url.RawQuery = vals.Encode()

	_, err := net.HTTPGet(s.BaseClient.GetHTTPClient(), url.String(), result, s.GetContextValues())
	if err != nil {
		log.Errorf("filter account failed, err:%s", err.Error())
		return nil, 0, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("查询账号列表失败，%s", result.Reason)
		log.Errorf("filter account failed, err:%s", err.Error())
		return nil, 0, result.ErrorCode
	}

	return result.Account, result.Total, result.ErrorCode
}

func (s *client) FilterAccountLite(filter *util.ContentFilter) ([]*common.AccountLite, int64, cd.ErrorCode) {
	result := &common.AccountLiteListResult{}

	vals := url.Values{}
	if filter == nil {
		filter = util.NewFilter()
	}
	filter.Set("mode", cd.LiteMode)
	vals = filter.Encode(vals)

	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.FilterAccount}, "")
	url.RawQuery = vals.Encode()

	_, err := net.HTTPGet(s.BaseClient.GetHTTPClient(), url.String(), result, s.GetContextValues())
	if err != nil {
		log.Errorf("filter account lite failed, err:%s", err.Error())
		return nil, 0, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("查询账号列表失败，%s", result.Reason)
		log.Errorf("filter account lite failed, err:%s", err.Error())
		return nil, 0, result.ErrorCode
	}

	return result.Account, result.Total, result.ErrorCode
}

func (s *client) QueryAccount(id int) (*common.AccountView, cd.ErrorCode) {
	result := &common.AccountResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.QueryAccount}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err := net.HTTPGet(s.BaseClient.GetHTTPClient(), url.String(), result, s.GetContextValues())
	if err != nil {
		log.Errorf("query account failed, err:%s", err.Error())
		return nil, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("查询指定账号失败，%s", result.Reason)
		log.Errorf("query account failed, err:%s", err.Error())
		return nil, result.ErrorCode
	}

	return result.Account, result.ErrorCode
}

func (s *client) CreateAccount(ptr *common.AccountParam) (*common.AccountView, cd.ErrorCode) {
	result := &common.AccountResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.CreateAccount}, "")
	url.RawQuery = vals.Encode()
	_, err := net.HTTPPost(s.BaseClient.GetHTTPClient(), url.String(), ptr, result, s.GetContextValues())
	if err != nil {
		log.Errorf("create account failed, err:%s", err.Error())
		return nil, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("新建账号失败，%s", result.Reason)
		log.Errorf("create account failed, err:%s", err.Error())
		return nil, result.ErrorCode
	}

	return result.Account, result.ErrorCode
}

func (s *client) UpdateAccount(id int, ptr *common.AccountParam) (*common.AccountView, cd.ErrorCode) {
	result := &common.AccountResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.UpdateAccount}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err := net.HTTPPut(s.BaseClient.GetHTTPClient(), url.String(), ptr, result, s.GetContextValues())
	if err != nil {
		log.Errorf("update account failed, err:%s", err.Error())
		return nil, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("修改账号失败，%s", result.Reason)
		log.Errorf("update account failed, err:%s", err.Error())
		return nil, result.ErrorCode
	}

	return result.Account, result.ErrorCode
}

func (s *client) DeleteAccount(id int) (*common.AccountView, cd.ErrorCode) {
	result := &common.AccountResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.DeleteAccount}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err := net.HTTPDelete(s.BaseClient.GetHTTPClient(), url.String(), result, s.GetContextValues())
	if err != nil {
		log.Errorf("delete account failed, err:%s", err.Error())
		return nil, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("删除账号失败，%s", result.Reason)
		log.Errorf("delete account failed, err:%s", err.Error())
		return nil, result.ErrorCode
	}

	return result.Account, result.ErrorCode
}

func (s *client) CheckAccount(account string) ([]*common.AccountLite, cd.ErrorCode) {
	result := &common.AccountLiteListResult{}

	vals := url.Values{}
	filter := util.NewFilter()
	filter.Set("account", account)
	filter.Set("mode", cd.LiteMode)
	vals = filter.Encode(vals)
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.FilterAccount}, "")
	url.RawQuery = vals.Encode()

	_, err := net.HTTPGet(s.BaseClient.GetHTTPClient(), url.String(), result, s.GetContextValues())
	if err != nil {
		log.Errorf("check account failed, err:%s", err.Error())
		return nil, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("查询终端失败，%s", result.Reason)
		log.Errorf("check account failed, err:%s", err.Error())
		return nil, result.ErrorCode
	}

	return result.Account, result.ErrorCode
}

func (s *client) FilterEndpoint(filter *util.ContentFilter) ([]*common.EndpointView, int64, cd.ErrorCode) {
	result := &common.EndpointListResult{}

	vals := url.Values{}
	if filter == nil {
		filter = util.NewFilter()
	}
	filter.Set("mode", cd.ViewMode)
	vals = filter.Encode(vals)

	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.FilterEndpoint}, "")
	url.RawQuery = vals.Encode()

	_, err := net.HTTPGet(s.BaseClient.GetHTTPClient(), url.String(), result, s.GetContextValues())
	if err != nil {
		log.Errorf("filter endpoint failed, err:%s", err.Error())
		return nil, 0, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("查询终端失败，%s", result.Reason)
		log.Errorf("filter endpoint failed, err:%s", err.Error())
		return nil, 0, result.ErrorCode
	}

	return result.Endpoint, result.Total, result.ErrorCode
}

func (s *client) FilterEndpointLite(filter *util.ContentFilter) ([]*common.EndpointLite, int64, cd.ErrorCode) {
	result := &common.EndpointLiteListResult{}

	vals := url.Values{}
	if filter == nil {
		filter = util.NewFilter()
	}
	filter.Set("mode", cd.LiteMode)
	vals = filter.Encode(vals)

	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.FilterEndpoint}, "")
	url.RawQuery = vals.Encode()

	_, err := net.HTTPGet(s.BaseClient.GetHTTPClient(), url.String(), result, s.GetContextValues())
	if err != nil {
		log.Errorf("filter endpoint lite failed, err:%s", err.Error())
		return nil, 0, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("查询终端失败，%s", result.Reason)
		log.Errorf("filter endpoint lite failed, err:%s", err.Error())
		return nil, 0, result.ErrorCode
	}

	return result.Endpoint, result.Total, result.ErrorCode
}

func (s *client) QueryEndpoint(id int) (*common.EndpointView, cd.ErrorCode) {
	result := &common.EndpointResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.QueryEndpoint}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err := net.HTTPGet(s.BaseClient.GetHTTPClient(), url.String(), result, s.GetContextValues())
	if err != nil {
		log.Errorf("query endpoint failed, err:%s", err.Error())
		return nil, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("查询指定终端失败，%s", result.Reason)
		log.Errorf("query endpoint failed, err:%s", err.Error())
		return nil, result.ErrorCode
	}

	return result.Endpoint, result.ErrorCode
}

func (s *client) CreateEndpoint(ptr *common.EndpointParam) (*common.EndpointView, cd.ErrorCode) {
	result := &common.EndpointResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.CreateEndpoint}, "")
	url.RawQuery = vals.Encode()
	_, err := net.HTTPPost(s.BaseClient.GetHTTPClient(), url.String(), ptr, result, s.GetContextValues())
	if err != nil {
		log.Errorf("create endpoint failed, err:%s", err.Error())
		return nil, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("新建终端失败，%s", result.Reason)
		log.Errorf("create endpoint failed, err:%s", err.Error())
		return nil, result.ErrorCode
	}

	return result.Endpoint, result.ErrorCode
}

func (s *client) UpdateEndpoint(id int, ptr *common.EndpointParam) (*common.EndpointView, cd.ErrorCode) {
	result := &common.EndpointResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.UpdateEndpoint}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err := net.HTTPPut(s.BaseClient.GetHTTPClient(), url.String(), ptr, result, s.GetContextValues())
	if err != nil {
		log.Errorf("update endpoint failed, err:%s", err.Error())
		return nil, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("更新终端失败，%s", result.Reason)
		log.Errorf("update endpoint failed, err:%s", err.Error())
		return nil, result.ErrorCode
	}

	return result.Endpoint, result.ErrorCode
}

func (s *client) DeleteEndpoint(id int) (*common.EndpointView, cd.ErrorCode) {
	result := &common.EndpointResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.DeleteEndpoint}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err := net.HTTPDelete(s.BaseClient.GetHTTPClient(), url.String(), result, s.GetContextValues())
	if err != nil {
		log.Errorf("delete endpoint failed, err:%s", err.Error())
		return nil, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("删除终端失败，%s", result.Reason)
		log.Errorf("delete endpoint failed, err:%s", err.Error())
		return nil, result.ErrorCode
	}

	return result.Endpoint, result.ErrorCode
}

func (s *client) FilterRole(filter *util.ContentFilter) ([]*common.RoleView, int64, cd.ErrorCode) {
	result := &common.RoleListResult{}

	vals := url.Values{}
	if filter == nil {
		filter = util.NewFilter()
	}
	filter.Set("mode", cd.ViewMode)
	vals = filter.Encode(vals)

	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.FilterRole}, "")
	url.RawQuery = vals.Encode()

	_, err := net.HTTPGet(s.BaseClient.GetHTTPClient(), url.String(), result, s.GetContextValues())
	if err != nil {
		log.Errorf("filter role failed, err:%s", err.Error())
		return nil, 0, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("查询权限组失败，%s", result.Reason)
		log.Errorf("filter role failed, err:%s", err.Error())
		return nil, 0, result.ErrorCode
	}

	return result.Role, result.Total, result.ErrorCode
}

func (s *client) FilterRoleLite(filter *util.ContentFilter) ([]*common.RoleLite, int64, cd.ErrorCode) {
	result := &common.RoleLiteListResult{}

	vals := url.Values{}
	if filter == nil {
		filter = util.NewFilter()
	}
	filter.Set("mode", cd.LiteMode)
	vals = filter.Encode(vals)

	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.FilterRole}, "")
	url.RawQuery = vals.Encode()

	_, err := net.HTTPGet(s.BaseClient.GetHTTPClient(), url.String(), result, s.GetContextValues())
	if err != nil {
		log.Errorf("filter role lite failed, err:%s", err.Error())
		return nil, 0, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("查询权限组失败，%s", result.Reason)
		log.Errorf("filter role lite failed, err:%s", err.Error())
		return nil, 0, result.ErrorCode
	}

	return result.Role, result.Total, result.ErrorCode
}

func (s *client) QueryRole(id int) (*common.RoleView, cd.ErrorCode) {
	result := &common.RoleResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.QueryRole}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err := net.HTTPGet(s.BaseClient.GetHTTPClient(), url.String(), result, s.GetContextValues())
	if err != nil {
		log.Errorf("query role failed, err:%s", err.Error())
		return nil, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("查询指定权限组失败，%s", result.Reason)
		log.Errorf("query role failed, err:%s", err.Error())
		return nil, result.ErrorCode
	}

	return result.Role, result.ErrorCode
}

func (s *client) CreateRole(ptr *common.RoleParam) (*common.RoleView, cd.ErrorCode) {
	result := &common.RoleResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.CreateRole}, "")
	url.RawQuery = vals.Encode()
	_, err := net.HTTPPost(s.BaseClient.GetHTTPClient(), url.String(), ptr, result, s.GetContextValues())
	if err != nil {
		log.Errorf("create role failed, err:%s", err.Error())
		return nil, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("新建权限组失败，%s", result.Reason)
		log.Errorf("create role failed, err:%s", err.Error())
		return nil, result.ErrorCode
	}

	return result.Role, result.ErrorCode
}

func (s *client) UpdateRole(id int, ptr *common.RoleParam) (*common.RoleView, cd.ErrorCode) {
	result := &common.RoleResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.UpdateRole}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err := net.HTTPPut(s.BaseClient.GetHTTPClient(), url.String(), ptr, result, s.GetContextValues())
	if err != nil {
		log.Errorf("update role failed, err:%s", err.Error())
		return nil, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("更新权限组失败，%s", result.Reason)
		log.Errorf("update role failed, err:%s", err.Error())
		return nil, result.ErrorCode
	}

	return result.Role, result.ErrorCode
}

func (s *client) DeleteRole(id int) (*common.RoleView, cd.ErrorCode) {
	result := &common.RoleResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.DeleteRole}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err := net.HTTPDelete(s.BaseClient.GetHTTPClient(), url.String(), result, s.GetContextValues())
	if err != nil {
		log.Errorf("delete role failed, err:%s", err.Error())
		return nil, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("删除权限组失败，%s", result.Reason)
		log.Errorf("delete role failed, err:%s", err.Error())
		return nil, result.ErrorCode
	}

	return result.Role, result.ErrorCode
}

func (s *client) FilterNamespace(filter *util.ContentFilter) ([]*common.NamespaceView, int64, cd.ErrorCode) {
	result := &common.NamespaceStatisticResult{}

	vals := url.Values{}
	if filter == nil {
		filter = util.NewFilter()
	}
	filter.Set("mode", cd.ViewMode)
	vals = filter.Encode(vals)
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.FilterNamespace}, "")
	url.RawQuery = vals.Encode()

	_, err := net.HTTPGet(s.BaseClient.GetHTTPClient(), url.String(), result, s.GetContextValues())
	if err != nil {
		log.Errorf("filter namespace failed, err:%s", err.Error())
		return nil, 0, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("查询租户失败，%s", result.Reason)
		log.Errorf("filter namespace failed, err:%s", err.Error())
		return nil, 0, result.ErrorCode
	}

	return result.Namespace, result.Total, result.ErrorCode
}

func (s *client) FilterNamespaceLite(filter *util.ContentFilter) ([]*common.NamespaceLite, int64, cd.ErrorCode) {
	result := &common.NamespaceLiteListResult{}

	vals := url.Values{}
	if filter == nil {
		filter = util.NewFilter()
	}
	filter.Set("mode", cd.LiteMode)
	vals = filter.Encode(vals)
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.FilterNamespace}, "")
	url.RawQuery = vals.Encode()

	_, err := net.HTTPGet(s.BaseClient.GetHTTPClient(), url.String(), result, s.GetContextValues())
	if err != nil {
		log.Errorf("filter namespace lite failed, err:%s", err.Error())
		return nil, 0, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("查询租户失败，%s", result.Reason)
		log.Errorf("filter namespace lite failed, err:%s", err.Error())
		return nil, 0, result.ErrorCode
	}

	return result.Namespace, result.Total, result.ErrorCode
}

func (s *client) QueryNamespace(id int) (*common.NamespaceView, cd.ErrorCode) {
	result := &common.NamespaceResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.QueryNamespace}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err := net.HTTPGet(s.BaseClient.GetHTTPClient(), url.String(), result, s.GetContextValues())
	if err != nil {
		log.Errorf("query namespace failed, err:%s", err.Error())
		return nil, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("查询指定租户失败，%s", result.Reason)
		log.Errorf("query namespace failed, err:%s", err.Error())
		return nil, result.ErrorCode
	}

	return result.Namespace, result.ErrorCode
}

func (s *client) CreateNamespace(ptr *common.NamespaceParam) (*common.NamespaceView, cd.ErrorCode) {
	result := &common.NamespaceResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.CreateNamespace}, "")
	url.RawQuery = vals.Encode()
	_, err := net.HTTPPost(s.BaseClient.GetHTTPClient(), url.String(), ptr, result, s.GetContextValues())
	if err != nil {
		log.Errorf("create namespace failed, err:%s", err.Error())
		return nil, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("新建租户失败，%s", result.Reason)
		log.Errorf("create namespace failed, err:%s", err.Error())
		return nil, result.ErrorCode
	}

	return result.Namespace, result.ErrorCode
}

func (s *client) UpdateNamespace(id int, ptr *common.NamespaceParam) (*common.NamespaceView, cd.ErrorCode) {
	result := &common.NamespaceResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.UpdateNamespace}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err := net.HTTPPut(s.BaseClient.GetHTTPClient(), url.String(), ptr, result, s.GetContextValues())
	if err != nil {
		log.Errorf("update namespace failed, err:%s", err.Error())
		return nil, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("更新租户失败，%s", result.Reason)
		log.Errorf("update namespace failed, err:%s", err.Error())
		return nil, result.ErrorCode
	}

	return result.Namespace, result.ErrorCode
}

func (s *client) DeleteNamespace(id int) (*common.NamespaceView, cd.ErrorCode) {
	result := &common.NamespaceResult{}

	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.BaseClient.GetServerURL())
	url.Path = strings.Join([]string{url.Path, common.ApiVersion, common.DeleteNamespace}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()

	_, err := net.HTTPDelete(s.BaseClient.GetHTTPClient(), url.String(), result, s.GetContextValues())
	if err != nil {
		log.Errorf("delete namespace failed, err:%s", err.Error())
		return nil, cd.UnExpected
	}

	if result.ErrorCode != cd.Success {
		err = fmt.Errorf("删除租户失败，%s", result.Reason)
		log.Errorf("delete namespace failed, err:%s", err.Error())
		return nil, result.ErrorCode
	}

	return result.Namespace, result.ErrorCode
}
