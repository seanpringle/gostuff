package query

import (
	"encoding/gob"
	"fmt"
	"os"
	"sync/atomic"
)

type Field = int
type TupleId = int64

const Id Field = 0

var (
	sequence int64 = 0
)

func NewId() TupleId {
	return TupleId(atomic.AddInt64(&sequence, 1))
}

func Save(path string, tables ...*Table) (wtf error) {

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(error); !ok {
				r = fmt.Errorf("unknown error")
			}
			wtf = fmt.Errorf("query.Save: %s", r.(error).Error())
		}
	}()

	assert := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	file, err := os.Create(path)
	assert(err)
	defer file.Close()

	save := gob.NewEncoder(file)
	assert(save.Encode(tables))

	return
}
