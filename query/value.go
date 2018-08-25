package query

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
	return v1 == v2
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
	return false
}
