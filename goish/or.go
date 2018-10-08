package main

import (
	"fmt"
)

type NodeOr struct {
	args Nodes
}

func NewNodeOr() *NodeOr {
	return &NodeOr{}
}

func (no *NodeOr) Consume(arg Node) {
	no.args = no.args.Push(arg)
}

func (no *NodeOr) Consumes() int {
	return 2
}

func (no *NodeOr) Produces() int {
	return 1
}

func (no *NodeOr) Precedence() int {
	return 2
}

func (no *NodeOr) Format() string {
	return fmt.Sprintf("func() Any { var a Any; a = %s; if !truth(a) { a = %s; }; return a }()", no.args[1].Format(), no.args[0].Format())
}

func (no *NodeOr) String() string {
	return fmt.Sprintf("NodeOr{%s}", no.args)
}
