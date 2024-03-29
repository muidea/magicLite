package common

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/util"

	"github.com/muidea/magicOrm/provider/remote"
	ou "github.com/muidea/magicOrm/util"
)

const (
	RegisterService   = "/service/register/"
	UnregisterService = "/service/unregister/"
	RegisterModel     = "/model/register/"
	UnregisterModel   = "/model/unregister/"
	QueryModel        = "/model/query/"
	CreateSchema      = "/schema/create/"
	DropSchema        = "/schema/drop/"
	InsertValue       = "/value/insert/"
	DeleteValue       = "/value/delete/"
	UpdateValue       = "/value/update/"
	QueryValue        = "/value/query/"
	QueryValues       = "/values/query/"
	QueryCount        = "/values/count/"
)

const BaseModule = "/kernel/base"

type ServiceParam struct {
	DBServer   string `json:"dbServer"`
	DBName     string `json:"dbName"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	MaxConnNum int    `json:"maxConnNum"`
}

func NewService(
	dbServer string,
	dbName string,
	username string,
	password string,
	maxConnNum int,
) *ServiceParam {

	return &ServiceParam{
		DBServer:   dbServer,
		DBName:     dbName,
		Username:   username,
		Password:   password,
		MaxConnNum: maxConnNum,
	}
}

// ModelParam model param
type ModelParam struct {
	remote.Object
}

// Result result
type Result cd.Result

type QueryModelResult struct {
	Result
	Model []*remote.Object `json:"model"`
}

// ObjectValueResult object value result
type ObjectValueResult struct {
	Result
	Value *remote.ObjectValue `json:"value"`
}

// SliceObjectValueResult slice object value result
type SliceObjectValueResult struct {
	Result
	Total int64                    `json:"total"`
	Value *remote.SliceObjectValue `json:"value"`
}

// ObjectValueCountResult query value count
type ObjectValueCountResult struct {
	Result
	Total int64 `json:"total"`
}

// QueryFilter value filter
type QueryFilter struct {
	EqualFilter    []*remote.ItemValue `json:"equal"`
	NotEqualFilter []*remote.ItemValue `json:"noEqual"`
	BelowFilter    []*remote.ItemValue `json:"below"`
	AboveFilter    []*remote.ItemValue `json:"above"`
	InFilter       []*remote.ItemValue `json:"in"`
	NotInFilter    []*remote.ItemValue `json:"notIn"`
	LikeFilter     []*remote.ItemValue `json:"like"`
	PageFilter     *util.Pagination    `json:"page"`
	MaskValue      *remote.ObjectValue `json:"maskValue"`
	SortFilter     *util.SortFilter    `json:"sort"`
}

// NewFilter new query filter
func NewFilter() *QueryFilter {
	return &QueryFilter{
		EqualFilter:    []*remote.ItemValue{},
		NotEqualFilter: []*remote.ItemValue{},
		BelowFilter:    []*remote.ItemValue{},
		AboveFilter:    []*remote.ItemValue{},
		InFilter:       []*remote.ItemValue{},
		NotInFilter:    []*remote.ItemValue{},
		LikeFilter:     []*remote.ItemValue{},
	}
}

func (s *QueryFilter) FromHttpRequest(req *http.Request) {
	filter := util.NewFilter()
	filter.Decode(req)

	s.FromContentFilter(filter)
}

func (s *QueryFilter) FromContentFilter(filter *util.ContentFilter) {
	if filter == nil {
		return
	}

	if filter.ParamItems != nil {
		for k, _ := range filter.ParamItems.Items {
			val := filter.ParamItems.GetEqual(k)
			if val != nil {
				s.Equal(k, val)
				continue
			}
			val = filter.ParamItems.GetNotEqual(k)
			if val != nil {
				s.NotEqual(k, val)
				continue
			}
			val = filter.ParamItems.GetBelow(k)
			if val != nil {
				s.Below(k, val)
				continue
			}
			val = filter.ParamItems.GetAbove(k)
			if val != nil {
				s.Above(k, val)
				continue
			}
			val = filter.ParamItems.GetIn(k)
			if val != nil {
				s.In(k, val)
				continue
			}
			val = filter.ParamItems.GetNotIn(k)
			if val != nil {
				s.NotIn(k, val)
				continue
			}
			val = filter.ParamItems.GetLike(k)
			if val != nil {
				s.Like(k, val)
				continue
			}
		}
	}

	if filter.Pagination != nil {
		s.PageFilter = filter.Pagination
		return
	}

	s.PageFilter = util.DefaultPagination()
}

func (s *QueryFilter) GetString(key string) (ret string, ok bool) {
	for _, item := range s.EqualFilter {
		if item.Name == key {
			ret, ok = (item.Value).(string)
			return
		}
	}

	return
}

func (s *QueryFilter) GetInt(key string) (ret int, ok bool) {
	for _, item := range s.EqualFilter {
		if item.Name == key {
			val, vOK := (item.Value).(float64)
			if !vOK {
				return
			}

			ret = int(val)
			ok = true
			return
		}
	}

	return
}

func (s *QueryFilter) Pagination(pageFilter *util.Pagination) {
	if pageFilter == nil {
		return
	}

	s.PageFilter = pageFilter
}

// Equal assign equal filter value
func (s *QueryFilter) Equal(key string, val interface{}) (err error) {
	qv := reflect.Indirect(reflect.ValueOf(val))
	qvType, qvErr := ou.GetTypeEnum(qv.Type())
	if qvErr != nil {
		err = qvErr
		return
	}
	if ou.IsSliceType(qvType) {
		err = fmt.Errorf("illegal value type, type:%s", qv.Type().String())
		return
	}

	if qvType == ou.TypeDateTimeField {
		val = qv.Interface().(time.Time).Format("2006-01-02 15:04:05")
	}

	if ou.IsBasicType(qvType) {
		item := &remote.ItemValue{Name: key, Value: val}
		s.EqualFilter = append(s.EqualFilter, item)
		return
	}

	if ou.IsMapType(qvType) {
		mVal, mErr := remote.GetMapValue(val)
		if mErr != nil {
			err = mErr
			return
		}

		item := &remote.ItemValue{Name: key, Value: mVal}
		s.EqualFilter = append(s.EqualFilter, item)
		return
	}

	objVal, objErr := remote.GetObjectValue(val)
	if objErr != nil {
		err = objErr
		return
	}

	item := &remote.ItemValue{Name: key, Value: objVal}
	s.EqualFilter = append(s.EqualFilter, item)

	return
}

// NotEqual  assign no equal filter value
func (s *QueryFilter) NotEqual(key string, val interface{}) (err error) {
	qv := reflect.Indirect(reflect.ValueOf(val))
	qvType, qvErr := ou.GetTypeEnum(qv.Type())
	if qvErr != nil {
		err = qvErr
		return
	}
	if ou.IsSliceType(qvType) {
		err = fmt.Errorf("illegal value type, type:%s", qv.Type().String())
		return
	}

	if qvType == ou.TypeDateTimeField {
		val = qv.Interface().(time.Time).Format("2006-01-02 15:04:05")
	}

	if ou.IsBasicType(qvType) {
		item := &remote.ItemValue{Name: key, Value: val}
		s.NotEqualFilter = append(s.NotEqualFilter, item)
		return
	}

	objVal, objErr := remote.GetObjectValue(val)
	if objErr != nil {
		err = objErr
		return
	}

	item := &remote.ItemValue{Name: key, Value: objVal}
	s.NotEqualFilter = append(s.NotEqualFilter, item)
	return nil
}

// Below assign below filter value
func (s *QueryFilter) Below(key string, val interface{}) (err error) {
	qv := reflect.Indirect(reflect.ValueOf(val))
	qvType, qvErr := ou.GetTypeEnum(qv.Type())
	if qvErr != nil {
		err = qvErr
		return
	}
	if !ou.IsBasicType(qvType) {
		err = fmt.Errorf("illegal value type, type:%s", qv.Type().String())
		return
	}

	if qvType == ou.TypeDateTimeField {
		val = qv.Interface().(time.Time).Format("2006-01-02 15:04:05")
	}

	item := &remote.ItemValue{Name: key, Value: val}
	s.BelowFilter = append(s.BelowFilter, item)

	return nil
}

// Above assign above filter value
func (s *QueryFilter) Above(key string, val interface{}) (err error) {
	qv := reflect.Indirect(reflect.ValueOf(val))
	qvType, qvErr := ou.GetTypeEnum(qv.Type())
	if qvErr != nil {
		err = qvErr
		return
	}
	if !ou.IsBasicType(qvType) {
		err = fmt.Errorf("illegal value type, type:%s", qv.Type().String())
		return
	}

	if qvType == ou.TypeDateTimeField {
		val = qv.Interface().(time.Time).Format("2006-01-02 15:04:05")
	}

	item := &remote.ItemValue{Name: key, Value: val}
	s.AboveFilter = append(s.AboveFilter, item)

	return nil
}

func (s *QueryFilter) getSliceValue(sliceVal interface{}) (ret interface{}, err error) {
	sliceReVal := reflect.Indirect(reflect.ValueOf(sliceVal))
	sliceValType, sliceValErr := ou.GetTypeEnum(sliceReVal.Type())
	if sliceValErr != nil {
		err = sliceValErr
		return
	}

	if !ou.IsSliceType(sliceValType) {
		err = fmt.Errorf("illegal value type, type:%s", sliceReVal.Type().String())
		return
	}

	if sliceReVal.Len() == 0 {
		return
	}

	svType := sliceReVal.Type().Elem()
	if svType.Kind() == reflect.Ptr {
		svType = svType.Elem()
	}

	subType, subErr := ou.GetTypeEnum(svType)
	if subErr != nil {
		err = subErr
		return
	}

	if ou.IsStructType(subType) {
		ret, err = remote.GetSliceObjectValue(sliceVal)
		return
	}

	retVal := []interface{}{}
	for idx := 0; idx < sliceReVal.Len(); idx++ {
		subV := reflect.Indirect(sliceReVal.Index(idx))
		if ou.TypeDateTimeField == subType {
			dtVal := subV.Interface().(time.Time).Format("2006-01-02 15:04:05")
			retVal = append(retVal, dtVal)

			continue
		}

		retVal = append(retVal, subV.Interface())
	}
	ret = retVal

	return
}

// In assign in filter value
func (s *QueryFilter) In(key string, val interface{}) (err error) {
	sliceVal, sliceErr := s.getSliceValue(val)
	if sliceErr != nil {
		err = sliceErr
		return
	}

	item := &remote.ItemValue{Name: key, Value: sliceVal}
	s.InFilter = append(s.InFilter, item)

	return
}

// NotIn assign notIn filter value
func (s *QueryFilter) NotIn(key string, val interface{}) (err error) {
	sliceVal, sliceErr := s.getSliceValue(val)
	if sliceErr != nil {
		err = sliceErr
		return
	}

	item := &remote.ItemValue{Name: key, Value: sliceVal}
	s.NotInFilter = append(s.NotInFilter, item)

	return nil
}

// Like assign like filter value
func (s *QueryFilter) Like(key string, val interface{}) (err error) {
	qv := reflect.Indirect(reflect.ValueOf(val))
	if qv.Kind() != reflect.String {
		err = fmt.Errorf("illegal value type, type:%s", qv.Type().String())
		return
	}

	item := &remote.ItemValue{Name: key, Value: val}
	s.LikeFilter = append(s.LikeFilter, item)

	return nil
}

// ValueMask assign mask value
func (s *QueryFilter) ValueMask(val interface{}) (err error) {
	qv := reflect.Indirect(reflect.ValueOf(val))
	qvType, qvErr := ou.GetTypeEnum(qv.Type())
	if qvErr != nil {
		err = qvErr
		return
	}
	if !ou.IsStructType(qvType) {
		err = fmt.Errorf("illegal mask value, type:%s", qv.Type().String())
		return
	}

	objVal, objErr := remote.GetObjectValue(val)
	if objErr != nil {
		err = objErr
		return
	}

	s.MaskValue = objVal
	return
}

// Page assign page filter value
func (s *QueryFilter) Page(filter *util.Pagination) {
	s.PageFilter = filter
}

// Sort sort result list
func (s *QueryFilter) Sort(sorter *util.SortFilter) {
	s.SortFilter = sorter
}

// ObjectValueFilter object value filter
type ObjectValueFilter struct {
	TypeName    string       `json:"typeName"`
	PkgPath     string       `json:"pkgPath"`
	ValueFilter *QueryFilter `json:"valueFilter"`
}
