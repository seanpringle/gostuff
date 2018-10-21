package main

import (
	"fmt"
)

type NodeLt struct {
	args Nodes
}

func NewNodeLt() *NodeLt {
	return &NodeLt{}
}

func (no *NodeLt) Consume(arg Node) {
	no.args = no.args.Push(arg)
}

func (no *NodeLt) Consumes() int {
	return 2
}

func (no *NodeLt) Produces() int {
	return 1
}

func (no *NodeLt) Precedence() int {
	return 4
}

func (no *NodeLt) Format() string {
	return fmt.Sprintf("Bool(%s)", no.FormatBool())
}

func (no *NodeLt) FormatBool() string {
	return fmt.Sprintf("lt(%s, %s)", FormatOne(no.args[1]), FormatOne(no.args[0]))
}

func (no *NodeLt) FormatOne() string {
	return no.Format()
}

func (no *NodeLt) String() string {
	return fmt.Sprintf("NodeLt{%s}", no.args)
}
