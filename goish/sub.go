package main

import (
	"fmt"
)

type NodeSub struct {
	args Nodes
}

func NewNodeSub() *NodeSub {
	return &NodeSub{}
}

func (na *NodeSub) Consume(n Node) {
	na.args = na.args.Push(n)
}

func (na *NodeSub) Consumes() int {
	return 2
}

func (na *NodeSub) Produces() int {
	return 1
}

func (na *NodeSub) Precedence() int {
	return 5
}

func (na *NodeSub) Format() string {
	return fmt.Sprintf("sub(%s, %s)", FormatOne(na.args[1]), FormatOne(na.args[0]))
}

func (na *NodeSub) FormatOne() string {
	return na.Format()
}

func (na *NodeSub) String() string {
	return fmt.Sprintf("NodeSub{%s, %s}", na.args[0], na.args[1])
}
