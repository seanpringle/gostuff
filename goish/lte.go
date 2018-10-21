package main

import (
	"fmt"
)

type NodeLte struct {
	args Nodes
}

func NewNodeLte() *NodeLte {
	return &NodeLte{}
}

func (no *NodeLte) Consume(arg Node) {
	no.args = no.args.Push(arg)
}

func (no *NodeLte) Consumes() int {
	return 2
}

func (no *NodeLte) Produces() int {
	return 1
}

func (no *NodeLte) Precedence() int {
	return 4
}

func (no *NodeLte) Format() string {
	return fmt.Sprintf("Bool(%s)", no.FormatBool())
}

func (no *NodeLte) FormatBool() string {
	return fmt.Sprintf("lte(%s, %s)", FormatOne(no.args[1]), FormatOne(no.args[0]))
}

func (no *NodeLte) FormatOne() string {
	return no.Format()
}

func (no *NodeLte) String() string {
	return fmt.Sprintf("NodeLte{%s}", no.args)
}
