package main

import (
	"fmt"
	"strconv"
)

type NodeLitDec struct {
	value float64
}

func NewNodeLitDec(value float64) *NodeLitDec {
	return &NodeLitDec{
		value: value,
	}
}

func (ld *NodeLitDec) Format() string {
	return fmt.Sprintf("Dec(%s)", strconv.FormatFloat(ld.value, 'f', -1, 64))
}

func (ld *NodeLitDec) String() string {
	return fmt.Sprintf("NodeLitDec{%s}", ld.Format())
}

func (ld *NodeLitDec) Produces() int {
	return 1
}
