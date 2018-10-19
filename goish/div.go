package main

import (
	"fmt"
)

type NodeDiv struct {
	args Nodes
}

func NewNodeDiv() *NodeDiv {
	return &NodeDiv{}
}

func (na *NodeDiv) Consume(n Node) {
	na.args = na.args.Push(n)
}

func (na *NodeDiv) Consumes() int {
	return 2
}

func (na *NodeDiv) Produces() int {
	return 1
}

func (na *NodeDiv) Precedence() int {
	return 6
}

func (na *NodeDiv) Format() string {
	return fmt.Sprintf("div(%s, %s)", FormatOne(na.args[1]), FormatOne(na.args[0]))
}

func (na *NodeDiv) String() string {
	return fmt.Sprintf("NodeDiv{%s, %s}", na.args[0], na.args[1])
}
