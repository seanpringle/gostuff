package query

import (
	"sync"
)

type Column map[TupleId]Value

type Table struct {
	Fields map[Field]Column
	Ids    map[TupleId]struct{}
	hooks  map[Hook]func(TupleId)
	lock   sync.Mutex
	log    []func() TupleId
}

func NewTable() *Table {
	return &Table{
		Fields: map[Field]Column{},
		Ids:    map[TupleId]struct{}{},
		hooks:  map[Hook]func(TupleId){},
	}
}

func (t *Table) Hook(hook Hook, fn func(TupleId)) {
	t.hooks[hook] = fn
}

func (t *Table) Unhook(hook Hook) {
	delete(t.hooks, hook)
}

func (t *Table) Commit() {
	t.lock.Lock()
	ids := map[TupleId]struct{}{}
	for _, fn := range t.log {
		ids[fn()] = struct{}{}
	}
	for id, _ := range ids {
		for _, fn := range t.hooks {
			fn(id)
		}
	}
	t.log = nil
	t.lock.Unlock()
}

func (t *Table) Upsert(tuple Tuple) TupleId {
	id := tuple.Id()
	tuple = tuple.Copy()
	t.lock.Lock()
	t.log = append(t.log, func() TupleId {
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
		return id
	})
	t.lock.Unlock()
	return id
}

func (t *Table) Delete(id TupleId) {
	t.lock.Lock()
	t.log = append(t.log, func() TupleId {
		delete(t.Ids, id)
		for field, column := range t.Fields {
			delete(column, id)
			if len(column) == 0 {
				delete(t.Fields, field)
			}
		}
		return id
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
