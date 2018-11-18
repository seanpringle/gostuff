package main

import (
	"fmt"
)

type NodeRShift struct {
	args Nodes
}

func NewNodeRShift() *NodeRShift {
	return &NodeRShift{}
}

func (na *NodeRShift) Consume(n Node) {
	na.args = na.args.Push(n)
}

func (na *NodeRShift) Consumes() int {
	return 2
}

func (na *NodeRShift) Produces() int {
	return 1
}

func (na *NodeRShift) Precedence() int {
	return 6
}

func (na *NodeRShift) Format() string {
	return fmt.Sprintf("rshift(%s, %s)", na.args[1].Format(), na.args[0].Format())
}

func (na *NodeRShift) String() string {
	return fmt.Sprintf("NodeRShift{%s, %s}", na.args[0], na.args[1])
}
