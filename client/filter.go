package client

import (
	"net/url"
	"strings"
	"time"
)

type RelationalOperator int

const (
	RelationalOperatorNone RelationalOperator = iota
	GreaterThan
	LessThan
	GreaterThanOrEqual
	LessThanOrEqual
)

func (r RelationalOperator) String() string {
	switch r {
	case GreaterThan:
		return "gt"
	case GreaterThanOrEqual:
		return "gte"
	case LessThan:
		return "lt"
	case LessThanOrEqual:
		return "lte"
	}
	return ""
}

type TimestampFilter struct {
	Timestamp []time.Time
	Operator  RelationalOperator
}

func (tf TimestampFilter) ToQuerystring(fieldName string) url.Values {
	v := url.Values{}
	if len(tf.Timestamp) > 0 {
		var t []string
		for _, ts := range tf.Timestamp {
			t = append(t, ts.Format(time.RFC3339))
		}
		key := fieldName
		if tf.Operator != RelationalOperatorNone {
			key = key + "[" + tf.Operator.String() + "]"
		}
		v.Set(key, strings.Join(t, ","))
	}
	return v
}

type Filter struct {
	Values []string
	Not    bool
}

func (f Filter) ToQueryString(fieldName string) url.Values {
	v := url.Values{}
	if len(f.Values) > 0 {
		key := fieldName
		if f.Not {
			key = key + "[not]"
		}
		v.Set(key, strings.Join(f.Values, ","))
	}
	return v
}

func appendQueryStrings(a, b url.Values) url.Values {
	for k, v := range a {
		b.Set(k, v[0]) // url.Values get only returns 1st item anyway
	}
	return b
}
