package query

import (
	"sync"
)

type Column map[TupleId]Value

type Table struct {
	Fields map[Field]Column
	Ids    map[TupleId]struct{}
	lock   sync.Mutex
	log    []func()
}

func NewTable() *Table {
	return &Table{
		Fields: map[Field]Column{},
		Ids:    map[TupleId]struct{}{},
	}
}

func (t *Table) Commit() {
	t.lock.Lock()
	for _, fn := range t.log {
		fn()
	}
	t.log = []func(){}
	t.lock.Unlock()
}

func (t *Table) Upsert(tuple Tuple) TupleId {
	id := tuple.Id()
	tuple = tuple.Copy()
	t.lock.Lock()
	t.log = append(t.log, func() {
		t.Ids[id] = struct{}{}
		for field, value := range tuple {
			if field == Id {
				continue
			}
			if _, exists := t.Fields[field]; !exists {
				t.Fields[field] = Column{}
			}
			column := t.Fields[field]
			column[id] = value
		}
	})
	t.lock.Unlock()
	return id
}

func (t *Table) Delete(id TupleId) {
	t.lock.Lock()
	t.log = append(t.log, func() {
		delete(t.Ids, id)
		for field, column := range t.Fields {
			delete(column, id)
			if len(column) == 0 {
				delete(t.Fields, field)
			}
		}
	})
	t.lock.Unlock()
}

func (t *Table) Insert(tuple Tuple) TupleId {
	id := tuple.Id()
	t.Delete(id)
	t.Upsert(tuple)
	return id
}

func (t *Table) Select(id TupleId, fields ...Field) Tuple {
	var tuple Tuple

	init := func() {
		if tuple == nil {
			tuple = Tuple{
				(Id): id,
			}
		}
	}

	get := func(field Field, column Column) {
		if column != nil && column[id] != nil {
			init()
			tuple[field] = VCopy(column[id])
		}
	}

	if len(fields) == 0 {
		for field, column := range t.Fields {
			get(field, column)
		}
	} else {
		for _, field := range fields {
			if field == Id {
				init()
			} else {
				get(field, t.Fields[field])
			}
		}
	}
	return tuple
}
