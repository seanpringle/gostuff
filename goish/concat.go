package main

import (
	"fmt"
)

type NodeConcat struct {
	args Nodes
}

func NewNodeConcat() *NodeConcat {
	return &NodeConcat{}
}

func (na *NodeConcat) Consume(n Node) {
	na.args = na.args.Push(n)
}

func (na *NodeConcat) Consumes() int {
	return 2
}

func (na *NodeConcat) Produces() int {
	return 1
}

func (na *NodeConcat) Precedence() int {
	return 5
}

func (na *NodeConcat) Format() string {
	return fmt.Sprintf("concat(%s, %s)", FormatOne(na.args[1]), FormatOne(na.args[0]))
}

func (na *NodeConcat) FormatOne() string {
	return na.Format()
}

func (na *NodeConcat) String() string {
	return fmt.Sprintf("NodeConcat{%s, %s}", na.args[0], na.args[1])
}
