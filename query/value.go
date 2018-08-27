package query

import (
	"fmt"
	"reflect"
	"strings"
)

type Value interface{}

type ValueCopy interface {
	Copy() Value
}

func VCopy(v Value) Value {
	if _, ok := v.(ValueCopy); ok {
		return v.(ValueCopy).Copy()
	}
	return v
}

type ValueString interface {
	String() string
}

func VString(v Value) string {
	if _, ok := v.(ValueString); ok {
		return v.(ValueString).String()
	}
	if _, ok := v.(string); ok {
		return v.(string)
	}
	return fmt.Sprintf("%v", v)
}

type ValueEq interface {
	Eq(Value) bool
}

func VEq(v1, v2 Value) bool {
	if _, ok := v1.(ValueEq); ok {
		return v1.(ValueEq).Eq(v2)
	}
	if _, ok := v2.(ValueEq); ok {
		return v2.(ValueEq).Eq(v1)
	}
	if reflect.TypeOf(v1) == reflect.TypeOf(v2) {
		return v1 == v2
	}
	return false
}

type ValueLt interface {
	Lt(Value) bool
}

func VLt(v1, v2 Value) bool {
	if _, ok := v1.(ValueLt); ok {
		return v1.(ValueLt).Lt(v2)
	}
	if _, ok := v2.(ValueLt); ok {
		return v2.(ValueLt).Lt(v1)
	}
	if s1, ok := v1.(string); ok {
		if s2, ok := v2.(string); ok {
			return strings.Compare(s1, s2) == -1
		}
	}
	if s1, ok := v1.(int64); ok {
		if s2, ok := v2.(int64); ok {
			return s1 < s2
		}
	}
	if s1, ok := v1.(float64); ok {
		if s2, ok := v2.(float64); ok {
			return s1 < s2
		}
	}
	return false
}

func VGt(v1, v2 Value) bool {
	return !VEq(v1, v2) && !VLt(v1, v2)
}
