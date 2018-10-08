package main

import (
	"fmt"
)

type NodeAnd struct {
	args Nodes
}

func NewNodeAnd() *NodeAnd {
	return &NodeAnd{}
}

func (no *NodeAnd) Consume(arg Node) {
	no.args = no.args.Push(arg)
}

func (no *NodeAnd) Consumes() int {
	return 2
}

func (no *NodeAnd) Produces() int {
	return 1
}

func (no *NodeAnd) Precedence() int {
	return 3
}

func (no *NodeAnd) Format() string {
	return fmt.Sprintf("func() Any { var a Any; a = %s; if truth(a) { var b Any; b = %s; if truth(b) { return b; }; }; return nil }()", no.args[1].Format(), no.args[0].Format())
}

func (no *NodeAnd) String() string {
	return fmt.Sprintf("NodeAnd{%s}", no.args)
}
