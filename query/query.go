package query

import (
	"sync"
)

type Query struct {
	table  *Table
	fields []Field
	ids    chan TupleId
	done   chan struct{}
	limit  int
	group  sync.WaitGroup
	stages int
}

func Select(t *Table, fields ...Field) *Query {
	q := &Query{
		table:  t,
		fields: fields,
		done:   make(chan struct{}, 1000),
	}

	output := make(chan TupleId, 8)
	q.group.Add(1)
	q.stages++
	go func() {
		for id, _ := range t.Ids {
			select {
			case <-q.done:
				break
			case output <- id:
			}
		}
		close(output)
		q.group.Done()
	}()

	q.ids = output
	return q
}

func (q *Query) Limit(n int) *Query {
	q.limit = n
	return q
}

func (q *Query) Where(field Field, fn func(Value) bool) *Query {
	output := make(chan TupleId, 8)
	input := q.ids
	q.ids = output
	q.stages++
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
	tuples := make(chan Tuple, 8)
	go func() {
		count := 0
		for id := range q.ids {
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
		for i := 1; i < q.stages; i++ {
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
