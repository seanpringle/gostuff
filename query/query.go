package query

import (
	"sync"
)

type Query struct {
	table   *Table
	fields  []Field
	limit   int
	done    chan struct{}
	group   sync.WaitGroup
	view    []TupleId
	filters []func(chan TupleId) chan TupleId
}

func Select(t *Table, fields ...Field) *Query {
	return &Query{
		table:  t,
		fields: fields,
		done:   make(chan struct{}, 1000),
	}
}

func (q *Query) View(ids ...TupleId) *Query {
	q.view = ids
	return q
}

func (q *Query) Limit(n int) *Query {
	q.limit = n
	return q
}

func (q *Query) Where(field Field, fn func(Value) bool) *Query {
	q.filters = append(q.filters, func(input chan TupleId) chan TupleId {
		output := make(chan TupleId, 8)
		q.group.Add(1)
		go func() {
			for id := range input {
				if fn(q.table.Fields[field][id]) {
					select {
					case <-q.done:
						break
					case output <- id:
					}
				}
			}
			close(output)
			q.group.Done()
		}()
		return output
	})
	return q
}

func (q *Query) In(field Field, values ...Value) *Query {
	fn := func(v Value) bool {
		for _, value := range values {
			if VEq(value, v) {
				return true
			}
		}
		return false
	}
	q.Where(field, fn)
	return q
}

func (q *Query) Eq(field Field, value Value) *Query {
	fn := func(v Value) bool {
		return VEq(value, v)
	}
	q.Where(field, fn)
	return q
}

func (q *Query) Ne(field Field, value Value) *Query {
	fn := func(v Value) bool {
		return !VEq(value, v)
	}
	q.Where(field, fn)
	return q
}

func (q *Query) Lt(field Field, value Value) *Query {
	fn := func(v Value) bool {
		return VLt(value, v)
	}
	q.Where(field, fn)
	return q
}

func (q *Query) Lte(field Field, value Value) *Query {
	fn := func(v Value) bool {
		return VEq(value, v) || VLt(value, v)
	}
	q.Where(field, fn)
	return q
}

func (q *Query) Gt(field Field, value Value) *Query {
	fn := func(v Value) bool {
		return !VEq(value, v) && !VLt(value, v)
	}
	q.Where(field, fn)
	return q
}

func (q *Query) Gte(field Field, value Value) *Query {
	fn := func(v Value) bool {
		return !VLt(value, v)
	}
	q.Where(field, fn)
	return q
}

func (q *Query) Run() chan Tuple {

	ids := make(chan TupleId, 8)
	q.group.Add(1)
	go func(ids chan TupleId) {
		if q.view != nil {
			// subset scan
			for _, id := range q.view {
				select {
				case <-q.done:
					break
				case ids <- id:
				}
			}
		} else {
			// table scan
			for id, _ := range q.table.Ids {
				select {
				case <-q.done:
					break
				case ids <- id:
				}
			}
		}
		close(ids)
		q.group.Done()
	}(ids)

	for _, fn := range q.filters {
		ids = fn(ids)
	}

	tuples := make(chan Tuple, 8)
	go func() {
		count := 0
		for id := range ids {
			select {
			case <-q.done:
				break
			case tuples <- q.table.Select(id, q.fields...):
			}
			count++
			if q.limit > 0 && count >= q.limit {
				break
			}
		}
		for i := 0; i < len(q.filters)+1; i++ {
			q.done <- struct{}{}
		}
		q.group.Wait()
		close(tuples)
	}()

	return tuples
}

func (q *Query) One() Tuple {
	return <-q.Limit(1).Run()
}

func (q *Query) All() []Tuple {
	var tuples []Tuple
	for tuple := range q.Run() {
		tuples = append(tuples, tuple)
	}
	return tuples
}

func (q *Query) List(field Field) []Value {
	var values []Value
	q.fields = []Field{field}
	for tuple := range q.Run() {
		values = append(values, tuple[field])
	}
	return values
}

func (q *Query) Table() *Table {
	t := NewTable()
	for tuple := range q.Run() {
		t.Insert(tuple)
	}
	t.Commit()
	return t
}
