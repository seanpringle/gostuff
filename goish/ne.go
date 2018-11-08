package main

import (
	"fmt"
)

type NodeNe struct {
	args Nodes
}

func NewNodeNe() *NodeNe {
	return &NodeNe{}
}

func (no *NodeNe) Consume(arg Node) {
	no.args = no.args.Push(arg)
}

func (no *NodeNe) Consumes() int {
	return 2
}

func (no *NodeNe) Produces() int {
	return 1
}

func (no *NodeNe) Precedence() int {
	return 4
}

func (no *NodeNe) Format() string {
	return fmt.Sprintf("Bool(%s)", no.FormatBool())
}

func (no *NodeNe) FormatBool() string {
	return fmt.Sprintf("!eq(%s, %s)", FormatOne(no.args[1]), FormatOne(no.args[0]))
}

func (no *NodeNe) FormatOne() string {
	return no.Format()
}

func (no *NodeNe) String() string {
	return fmt.Sprintf("NodeNe{%s}", no.args)
}
