package main

import (
	"fmt"
)

type NodeNot struct {
	node Node
}

func NewNodeNot(arg Node) *NodeNot {
	return &NodeNot{arg}
}

func (no *NodeNot) Produces() int {
	return 1
}

func (no *NodeNot) Format() string {
	return fmt.Sprintf("Bool(!truth(%s))", FormatOne(no.node))
}

func (no *NodeNot) FormatBool() string {
	return fmt.Sprintf("!truth(%s)", FormatOne(no.node))
}

func (no *NodeNot) FormatOne() string {
	return no.Format()
}

func (no *NodeNot) String() string {
	return fmt.Sprintf("NodeNot{%s}", no.node)
}
