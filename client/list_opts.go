package client

import (
	"net/url"
	"strconv"
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
	Page    int
	PerPage int
	OrderBy string

	LabelSelector Filter
	CreateAts     TimestampFilter
	UpdatedAts    TimestampFilter
}

func NewListOptions() *ListOptions {
	return &ListOptions{
		Page:    DefaultPage,
		PerPage: DefaultPageSize,
	}
}

func (lo ListOptions) ToQueryString() url.Values {
	v := url.Values{}
	if lo.Page > 0 {
		v.Set(PageField, strconv.Itoa(lo.Page))
	}
	if lo.PerPage > 0 {
		v.Set(PerPageField, strconv.Itoa(lo.PerPage))
	}
	if len(lo.OrderBy) > 0 {
		v.Set(OrderByField, lo.OrderBy)
	}
	v = appendQueryStrings(v, lo.LabelSelector.ToQueryString(LabelSelectorField))
	v = appendQueryStrings(v, lo.CreateAts.ToQuerystring(CreatedAtsField))
	v = appendQueryStrings(v, lo.UpdatedAts.ToQuerystring(UpdatedAtsField))
	return v
}
