package main

import (
	"fmt"
)

type NodeLShift struct {
	args Nodes
}

func NewNodeLShift() *NodeLShift {
	return &NodeLShift{}
}

func (na *NodeLShift) Consume(n Node) {
	na.args = na.args.Push(n)
}

func (na *NodeLShift) Consumes() int {
	return 2
}

func (na *NodeLShift) Produces() int {
	return 1
}

func (na *NodeLShift) Precedence() int {
	return 6
}

func (na *NodeLShift) Format() string {
	return fmt.Sprintf("lshift(%s, %s)", na.args[1].Format(), na.args[0].Format())
}

func (na *NodeLShift) String() string {
	return fmt.Sprintf("NodeLShift{%s, %s}", na.args[0], na.args[1])
}
