package query

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
