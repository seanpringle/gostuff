package main

import (
	"fmt"
)

type NodeBoolAnd struct {
	args Nodes
}

func NewNodeBoolAnd() *NodeBoolAnd {
	return &NodeBoolAnd{}
}

func (na *NodeBoolAnd) Consume(n Node) {
	na.args = na.args.Push(n)
}

func (na *NodeBoolAnd) Consumes() int {
	return 2
}

func (na *NodeBoolAnd) Produces() int {
	return 1
}

func (na *NodeBoolAnd) Precedence() int {
	return 6
}

func (na *NodeBoolAnd) Format() string {
	return fmt.Sprintf("b_and(%s, %s)", FormatOne(na.args[1]), FormatOne(na.args[0]))
}

func (na *NodeBoolAnd) String() string {
	return fmt.Sprintf("NodeBoolAnd{%s, %s}", na.args[0], na.args[1])
}
