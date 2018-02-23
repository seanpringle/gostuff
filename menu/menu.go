package menu

import (
	"fmt"
	"github.com/seanpringle/gostuff/imath"
	"sort"
	"strings"
)

type Menu struct {
	Right    func()
	Left     func()
	Set      func(string)
	Ins      func(string)
	Back     func()
	Input    func() string
	Position func() int
	Matches  func() []string
	Choose   func()
}

func New(prefix string, items []string) *Menu {

	self := &Menu{}
	sort.Strings(items)

	input := ""
	matches := []string{}
	position := int(0)

	filter := func() {
		matches = []string{}
		for _, item := range items {
			if len(input) == 0 || strings.Contains(strings.ToLower(item), strings.ToLower(input)) {
				matches = append(matches, item)
			}
		}
		position = imath.Max(0, imath.Min(len(matches)-1, position))
	}

	filter()

	self.Right = func() {
		position = imath.Max(0, imath.Min(len(matches)-1, position+1))
	}

	self.Left = func() {
		position = imath.Max(0, imath.Min(len(matches)-1, position-1))
	}

	self.Set = func(str string) {
		input = str
		filter()
	}

	self.Ins = func(str string) {
		input = strings.Join([]string{input, str}, "")
		filter()
	}

	self.Back = func() {
		if len(input) > 0 {
			input = input[:len(input)-1]
			filter()
		}
	}

	self.Input = func() string {
		return fmt.Sprintf("%s%s", prefix, input)
	}

	self.Position = func() int {
		return position
	}

	self.Matches = func() []string {
		return matches
	}

	self.Choose = func() {
		if len(matches) > 0 {
			input = matches[position]
			filter()
		}
	}

	return self
}
