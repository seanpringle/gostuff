package main

import (
	"fmt"
)

type NodeLitStr struct {
	value string
}

func NewNodeLitStr(value string) *NodeLitStr {
	return &NodeLitStr{
		value: value,
	}
}

func (ls *NodeLitStr) Format() string {
	return fmt.Sprintf(`Str{%s}`, ls.value)
}

func (ls *NodeLitStr) String() string {
	return fmt.Sprintf("NodeLitStr{%s}", ls.Format())
}

func (ls *NodeLitStr) Produces() int {
	return 1
}

func (ls *NodeLitStr) FormatKey() string {
	return fmt.Sprintf("Str{%s}", ls.value)
}
