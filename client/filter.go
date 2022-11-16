package client

import (
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

type Filter struct {
	Values []string
	Not    bool
}
