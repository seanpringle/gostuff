package main

import (
	"fmt"
)

type NodeBoolXor struct {
	args Nodes
}

func NewNodeBoolXor() *NodeBoolXor {
	return &NodeBoolXor{}
}

func (na *NodeBoolXor) Consume(n Node) {
	na.args = na.args.Push(n)
}

func (na *NodeBoolXor) Consumes() int {
	return 2
}

func (na *NodeBoolXor) Produces() int {
	return 1
}

func (na *NodeBoolXor) Precedence() int {
	return 6
}

func (na *NodeBoolXor) Format() string {
	return fmt.Sprintf("b_xor(%s, %s)", FormatOne(na.args[1]), FormatOne(na.args[0]))
}

func (na *NodeBoolXor) String() string {
	return fmt.Sprintf("NodeBoolXor{%s, %s}", na.args[0], na.args[1])
}
