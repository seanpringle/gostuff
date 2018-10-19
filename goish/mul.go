package main

import (
	"fmt"
)

type NodeMul struct {
	args Nodes
}

func NewNodeMul() *NodeMul {
	return &NodeMul{}
}

func (na *NodeMul) Consume(n Node) {
	na.args = na.args.Push(n)
}

func (na *NodeMul) Consumes() int {
	return 2
}

func (na *NodeMul) Produces() int {
	return 1
}

func (na *NodeMul) Precedence() int {
	return 6
}

func (na *NodeMul) Format() string {
	return fmt.Sprintf("mul(%s, %s)", FormatOne(na.args[1]), FormatOne(na.args[0]))
}

func (na *NodeMul) String() string {
	return fmt.Sprintf("NodeMul{%s, %s}", na.args[0], na.args[1])
}
