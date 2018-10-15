package main

import (
	"fmt"
)

type NodeGt struct {
	args Nodes
}

func NewNodeGt() *NodeGt {
	return &NodeGt{}
}

func (no *NodeGt) Consume(arg Node) {
	no.args = no.args.Push(arg)
}

func (no *NodeGt) Consumes() int {
	return 2
}

func (no *NodeGt) Produces() int {
	return 1
}

func (no *NodeGt) Precedence() int {
	return 4
}

func (no *NodeGt) Format() string {
	return fmt.Sprintf("Bool{gt(one(%s),one(%s))}", no.args[1].Format(), no.args[0].Format())
}

func (no *NodeGt) String() string {
	return fmt.Sprintf("NodeGt{%s}", no.args)
}
