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
	return fmt.Sprintf("Bool{lt(%s,%s)}", no.args[1].Format(), no.args[0].Format())
}

func (no *NodeLt) String() string {
	return fmt.Sprintf("NodeLt{%s}", no.args)
}