package query

type Column map[TupleId]Value

type Table struct {
	Fields map[Field]Column
	log    []func()
}

func NewTable() *Table {
	return &Table{
		Fields: map[Field]Column{},
	}
}

func (t *Table) Commit() {
	for _, fn := range t.log {
		fn()
	}
	t.log = []func(){}
}

func (t *Table) Insert(tuple Tuple) TupleId {
	id := tuple.Id()
	tuple = tuple.Copy()
	t.log = append(t.log, func() {
		for _, column := range t.Fields {
			delete(column, id)
		}
		for field, value := range tuple {
			if _, exists := t.Fields[field]; !exists {
				t.Fields[field] = Column{}
			}
			column := t.Fields[field]
			column[id] = value
		}
		for field, column := range t.Fields {
			if len(column) == 0 {
				delete(t.Fields, field)
			}
		}
	})
	return id
}

func (t *Table) Delete(id TupleId) {
	t.log = append(t.log, func() {
		for field, column := range t.Fields {
			delete(column, id)
			if len(column) == 0 {
				delete(t.Fields, field)
			}
		}
	})
}

func (t *Table) Select(ids []TupleId) []Tuple {
	var r []Tuple
	for _, id := range ids {
		var tuple Tuple
		for field, column := range t.Fields {
			if column[id] != nil {
				if tuple == nil {
					tuple = Tuple{
						(Id): id,
					}
				}
				tuple[field] = VCopy(column[id])
			}
		}
		if tuple != nil {
			r = append(r, tuple)
		}
	}
	return r
}

func (t *Table) Equal(field Field, value Value) []TupleId {
	var ids []TupleId
	for id, fv := range t.Fields[field] {
		if VEqual(value, fv) {
			ids = append(ids, id)
		}
	}
	return ids
}
