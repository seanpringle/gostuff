package main

import (
	"fmt"
)

type NodeMod struct {
	args Nodes
}

func NewNodeMod() *NodeMod {
	return &NodeMod{}
}

func (na *NodeMod) Consume(n Node) {
	na.args = na.args.Push(n)
}

func (na *NodeMod) Consumes() int {
	return 2
}

func (na *NodeMod) Produces() int {
	return 1
}

func (na *NodeMod) Precedence() int {
	return 6
}

func (na *NodeMod) Format() string {
	return fmt.Sprintf("mod(%s, %s)", FormatOne(na.args[1]), FormatOne(na.args[0]))
}

func (na *NodeMod) String() string {
	return fmt.Sprintf("NodeMod{%s, %s}", na.args[0], na.args[1])
}
