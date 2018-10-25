package main

import (
	"fmt"
)

type NodeLitRune struct {
	value string
}

func NewNodeLitRune(value string) *NodeLitRune {
	return &NodeLitRune{
		value: value,
	}
}

func (ls *NodeLitRune) Format() string {
	return fmt.Sprintf(`Rune(%s)`, ls.value)
}

func (ls *NodeLitRune) FormatOne() string {
	return ls.Format()
}

func (ls *NodeLitRune) String() string {
	return fmt.Sprintf("NodeLitRune{%s}", ls.Format())
}

func (ls *NodeLitRune) Produces() int {
	return 1
}

func (ls *NodeLitRune) FormatKey() string {
	return fmt.Sprintf("Rune(%s)", ls.value)
}
