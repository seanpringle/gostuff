package history

import (
	"github.com/seanpringle/gostuff/imath"
)

type History struct {
	Next func() string
	Prev func() string
	Add  func(string)
	Last func() string
	Get  func() string
}

func New() *History {

	self := &History{}

	list := []string{}
	pos := 0

	to := func(p int) {
		pos = imath.Max(0, imath.Min(len(list)-1, p))
	}

	forward := func() {
		to(pos + 1)
	}

	backward := func() {
		to(pos - 1)
	}

	current := func() string {
		if len(list) > pos {
			return list[pos]
		}
		return ""
	}

	drop := func(str string) {
		filter := []string{}
		for _, item := range list {
			if item != str {
				filter = append(filter, item)
			}
		}
		list = filter
	}

	add := func(str string) {
		list = append(list, str)
	}

	self.Next = func() string {
		forward()
		return current()
	}

	self.Prev = func() string {
		backward()
		return current()
	}

	self.Last = func() string {
		to(len(list) - 1)
		return current()
	}

	self.Add = func(str string) {
		drop(str)
		add(str)
		to(len(list) - 1)
	}

	self.Get = func() string {
		return current()
	}

	return self
}
