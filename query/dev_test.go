package query

import (
	"log"
	"testing"
)

const (
	Name Field = Id + 1 + iota
)

func Test1(test *testing.T) {

	t1 := NewTable()

	id := NewId()

	t1.Insert(Tuple{
		(Id):   id,
		(Name): "Sean",
	})

	t1.Insert(Tuple{
		(Id):   NewId(),
		(Name): "Eliza",
	})

	t1.Commit()
	Save("save", t1)

	for _, tuple := range t1.Select(t1.Equal(Name, id)) {
		log.Println(tuple)
	}
}
