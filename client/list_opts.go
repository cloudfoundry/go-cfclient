package client

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 50
)

const (
	PageField          = "page"
	PerPageField       = "per_page"
	OrderByField       = "order_by"
	LabelSelectorField = "label_selector"
	CreatedAtsField    = "created_ats"
	UpdatedAtsField    = "updated_ats"
	IncludeField       = "include"
	GUIDsField         = "guids"
	NamesField         = "names"
)

// ListOptions is the shared common type for all other list option types
type ListOptions struct {
	Page    int    `filter:"page,omitempty"`
	PerPage int    `filter:"per_page,omitempty"`
	OrderBy string `filter:"order_by,omitempty"`

	LabelSelector Filter
	CreateAts     TimestampFilter `filter:"created_ats,omitempty"`
	UpdatedAts    TimestampFilter `filter:"updated_ats,omitempty"`
}

func NewListOptions() *ListOptions {
	return &ListOptions{
		Page:    DefaultPage,
		PerPage: DefaultPageSize,
	}
}

func (lo ListOptions) ToQueryString() url.Values {
	s := ListOptionsSerializer{}
	val, err := s.Serialize(&lo)
	if err != nil {
		panic("go-cfclient filter bug: " + err.Error())
	}
	return val
	// v := url.Values{}
	// if lo.Page > 0 {
	// 	v.Set(PageField, strconv.Itoa(lo.Page))
	// }
	// if lo.PerPage > 0 {
	// 	v.Set(PerPageField, strconv.Itoa(lo.PerPage))
	// }
	// if len(lo.OrderBy) > 0 {
	// 	v.Set(OrderByField, lo.OrderBy)
	// }
	// v = appendQueryStrings(v, lo.LabelSelector.ToQueryString(LabelSelectorField))
	// v = appendQueryStrings(v, lo.CreateAts.ToQuerystring(CreatedAtsField))
	// v = appendQueryStrings(v, lo.UpdatedAts.ToQuerystring(UpdatedAtsField))
	// return v
}

var filterType = reflect.TypeOf(Filter{})
var timeFilterType = reflect.TypeOf(TimestampFilter{})
var timeType = reflect.TypeOf(time.Time{})

type ListOptionsSerializer struct{}

func (los ListOptionsSerializer) Serialize(o any) (url.Values, error) {
	if o == nil {
		return url.Values{}, nil
	}

	val := reflect.ValueOf(o)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return url.Values{}, nil
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return url.Values{}, nil
	}

	return los.reflectValues(val)
}

func (los ListOptionsSerializer) reflectValues(val reflect.Value) (url.Values, error) {
	vals := url.Values{}
	for i := 0; i < val.Type().NumField(); i++ {
		sv := val.Field(i)
		rawTag := val.Type().Field(i).Tag.Get("filter")
		if rawTag == "" {
			continue
		}
		tag, err := parseTag(rawTag)
		if err != nil {
			return vals, err
		}
		if tag.name == "-" {
			continue
		}
		if tag.omitEmpty && isEmptyValue(sv) {
			continue
		}

		sv = getNonPointerValue(sv)
		switch sv.Type() {
		case filterType:
			// TODO: parse filter
			break
		case timeFilterType:
			tfv := reflect.ValueOf(sv)
			tfv = getNonPointerValue(tfv)
			var timestamps []string
			var relationalOperator RelationalOperator
			for tfi := 0; i < tfv.Type().NumField(); tfi++ {
				tfsv := tfv.Field(tfi)
				if tfsv.Kind() == reflect.Slice || tfsv.Kind() == reflect.Array {
					for tfsvi := 0; tfsvi < tfsv.Len(); tfsvi++ {
						tsv := tfsv.Index(tfsvi)
						if tsv.Type() == timeType {
							timestamp := tsv.Interface().(time.Time)
							timestamps = append(timestamps, timestamp.Format(time.RFC3339))
						}
					}
				} else if tfsv.Kind() == reflect.Int {
					relationalOperator = tfsv.Interface().(RelationalOperator)
				}
			}
			if len(timestamps) > 0 {
				key := tag.name
				if relationalOperator != RelationalOperatorNone {
					key = key + "[" + relationalOperator.String() + "]"
				}
				vals.Add(key, strings.Join(timestamps, ","))
			}
			break
		default:
			vals.Add(tag.name, fmt.Sprint(sv.Interface()))
		}
	}

	return vals, nil
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func getNonPointerValue(sv reflect.Value) reflect.Value {
	for sv.Kind() == reflect.Ptr {
		if sv.IsNil() {
			break
		}
		sv = sv.Elem()
	}
	return sv
}

func parseTag(tagValue string) (filterTag, error) {
	s := strings.Split(tagValue, ",")
	if len(s) == 2 {
		return filterTag{
			name:      s[0],
			omitEmpty: s[1] == "omitempty",
		}, nil
	} else if len(s) == 1 {
		return filterTag{
			name:      s[0],
			omitEmpty: false,
		}, nil
	}
	return filterTag{}, errors.New("missing required filter tag name")
}

type filterTag struct {
	name      string
	omitEmpty bool
}
