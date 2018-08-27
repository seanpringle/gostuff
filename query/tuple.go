package query

import (
	"encoding/gob"
	"fmt"
)

func init() {
	gob.Register(TupleId(0))
}

type TupleId int64

func (t TupleId) Eq(v Value) bool {
	if _, ok := v.(TupleId); ok {
		return t == v.(TupleId)
	}
	return false
}

func (t TupleId) Lt(v Value) bool {
	if _, ok := v.(TupleId); ok {
		return t < v.(TupleId)
	}
	return false
}

func (t TupleId) String() string {
	return fmt.Sprintf("%d", t)
}

type Tuple map[Field]Value

func (t Tuple) Copy() Tuple {
	tuple := Tuple{}
	for field, value := range t {
		tuple[field] = VCopy(value)
	}
	return tuple
}

func (t Tuple) Id() TupleId {
	if _, ok := t[Id]; !ok {
		t[Id] = NewId()
	}
	return t[Id].(TupleId)
}
