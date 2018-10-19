package main

import (
	"fmt"
)

type NodeAdd struct {
	args Nodes
}

func NewNodeAdd() *NodeAdd {
	return &NodeAdd{}
}

func (na *NodeAdd) Consume(n Node) {
	na.args = na.args.Push(n)
}

func (na *NodeAdd) Consumes() int {
	return 2
}

func (na *NodeAdd) Produces() int {
	return 1
}

func (na *NodeAdd) Precedence() int {
	return 5
}

func (na *NodeAdd) Format() string {
	return fmt.Sprintf("add(%s, %s)", FormatOne(na.args[1]), FormatOne(na.args[0]))
}

func (na *NodeAdd) FormatOne() string {
	return na.Format()
}

func (na *NodeAdd) String() string {
	return fmt.Sprintf("NodeAdd{%s, %s}", na.args[0], na.args[1])
}
