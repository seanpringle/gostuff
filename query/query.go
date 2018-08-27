package query

import (
	"sort"
	"strings"
	"sync"
)

type Direction = int
type Aggregate = int

const (
	ASC Direction = iota + 1
	DESC
)

const (
	MAX Aggregate = iota + 1
	MIN
	ANY
)

type order struct {
	Field
	Direction
}

type Query struct {
	table   *Table
	fields  []Field
	limit   int
	done    chan struct{}
	group   sync.WaitGroup
	view    []TupleId
	filters []func(chan TupleId) chan TupleId
	groups  []Field
	orders  []order
	aggs    map[Field]Aggregate
}

func Select(t *Table, fields ...Field) *Query {
	return &Query{
		table:  t,
		fields: fields,
		aggs:   map[Field]Aggregate{},
		done:   make(chan struct{}, 1000),
	}
}

func (q *Query) agg(agg Aggregate, fields ...Field) {
	for _, field := range fields {
		q.aggs[field] = agg
	}
}

func (q *Query) Min(fields ...Field) *Query {
	q.agg(MIN, fields...)
	return q
}

func (q *Query) Max(fields ...Field) *Query {
	q.agg(MAX, fields...)
	return q
}

func (q *Query) Any(fields ...Field) *Query {
	q.agg(ANY, fields...)
	return q
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
			column := q.table.Fields[field]
			for id := range input {
				if (field == Id && fn(id)) || fn(column[id]) {
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

func (q *Query) Group(fields ...Field) *Query {
	for _, field := range fields {
		q.groups = append(q.groups, field)
	}
	return q
}

func (q *Query) Order(field Field, direction Direction) *Query {
	q.orders = append(q.orders, order{field, direction})
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

	stages := len(q.filters)

	fetch := make(chan Tuple, 8)
	stages++
	q.group.Add(1)
	go func() {
		for id := range ids {
			select {
			case <-q.done:
				break
			case fetch <- q.table.Select(id, q.fields...):
			}
		}
		close(fetch)
		q.group.Done()
	}()

	results := fetch

	if len(q.groups) > 0 {
		stages++
		q.group.Add(1)
		group := make(chan Tuple, 8)
		go func(input chan Tuple) {
			aggregate := map[string]Tuple{}
			for tuple := range input {
				var keys []string
				for _, field := range q.groups {
					keys = append(keys, VString(tuple[field]))
				}
				key := strings.Join(keys, ",")
				if _, ok := aggregate[key]; !ok {
					aggregate[key] = Tuple{}
					for _, field := range q.groups {
						aggregate[key][field] = tuple[field]
					}
					for field, _ := range q.aggs {
						aggregate[key][field] = tuple[field]
					}
				} else {
					for field, method := range q.aggs {
						switch method {
						case MIN:
							if VLt(tuple[field], aggregate[key][field]) {
								aggregate[key][field] = tuple[field]
							}
						case MAX:
							if VGt(tuple[field], aggregate[key][field]) {
								aggregate[key][field] = tuple[field]
							}
						}
					}
				}
			}
			for _, tuple := range aggregate {
				select {
				case <-q.done:
					break
				case group <- tuple:
				}
			}
			close(group)
			q.group.Done()
		}(results)
		results = group
	}

	if len(q.orders) > 0 {
		stages++
		q.group.Add(1)
		order := make(chan Tuple, 8)
		go func(input chan Tuple) {
			var tuples []Tuple
			for tuple := range input {
				tuples = append(tuples, tuple)
			}
			for _, level := range q.orders {
				if level.Direction == ASC {
					sort.SliceStable(tuples, func(i, j int) bool {
						return VLt(tuples[i][level.Field], tuples[j][level.Field])
					})
				} else {
					sort.SliceStable(tuples, func(i, j int) bool {
						return VGt(tuples[i][level.Field], tuples[j][level.Field])
					})
				}
			}
			for _, tuple := range tuples {
				select {
				case <-q.done:
					break
				case order <- tuple:
				}
			}
			close(order)
			q.group.Done()
		}(results)
		results = order
	}

	limit := make(chan Tuple, 8)
	//q.group.Add(1)
	go func(input chan Tuple) {
		count := 0
		for tuple := range input {
			limit <- tuple
			count++
			if q.limit > 0 && count >= q.limit {
				break
			}
		}
		for i := 0; i < stages; i++ {
			q.done <- struct{}{}
		}
		close(limit)
		q.group.Wait()
	}(results)
	results = limit

	return results
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
	//q.fields = []Field{field}
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
