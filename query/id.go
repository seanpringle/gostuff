package query

import (
	"sync/atomic"
)

type Field = int
type Hook = int

const Id Field = 0

var sequence int64 = 0

func NewId() TupleId {
	return TupleId(atomic.AddInt64(&sequence, 1))
}
