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
	return fmt.Sprintf("Bool(lte(one(%s),one(%s)))", no.args[1].Format(), no.args[0].Format())
}

func (no *NodeLte) String() string {
	return fmt.Sprintf("NodeLte{%s}", no.args)
}
