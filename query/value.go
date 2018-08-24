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

type ValueEqual interface {
	Equal(Value) bool
}

func VEqual(v1, v2 Value) bool {
	if _, ok := v1.(ValueEqual); ok {
		return v1.(ValueEqual).Equal(v2)
	}
	if _, ok := v2.(ValueEqual); ok {
		return v2.(ValueEqual).Equal(v1)
	}
	return false
}
