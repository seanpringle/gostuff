package main

import (
	"fmt"
)

type NodeBoolOr struct {
	args Nodes
}

func NewNodeBoolOr() *NodeBoolOr {
	return &NodeBoolOr{}
}

func (na *NodeBoolOr) Consume(n Node) {
	na.args = na.args.Push(n)
}

func (na *NodeBoolOr) Consumes() int {
	return 2
}

func (na *NodeBoolOr) Produces() int {
	return 1
}

func (na *NodeBoolOr) Precedence() int {
	return 6
}

func (na *NodeBoolOr) Format() string {
	return fmt.Sprintf("b_or(%s, %s)", FormatOne(na.args[1]), FormatOne(na.args[0]))
}

func (na *NodeBoolOr) String() string {
	return fmt.Sprintf("NodeBoolOr{%s, %s}", na.args[0], na.args[1])
}
