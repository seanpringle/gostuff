package main

import (
	"fmt"
)

type NodeLitInt struct {
	value int64
}

func NewNodeLitInt(value int64) *NodeLitInt {
	return &NodeLitInt{
		value: value,
	}
}

func (li *NodeLitInt) Format() string {
	return fmt.Sprintf("Int(%d)", li.value)
}

func (li *NodeLitInt) FormatOne() string {
	return li.Format()
}

func (li *NodeLitInt) String() string {
	return fmt.Sprintf("NodeLitInt{%d}", li.value)
}

func (li *NodeLitInt) Produces() int {
	return 1
}

func (li *NodeLitInt) FormatKey() string {
	return fmt.Sprintf(`Int(%d)`, li.value)
}
