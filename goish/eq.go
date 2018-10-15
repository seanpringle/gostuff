package main

import (
	"fmt"
)

type NodeEq struct {
	args Nodes
}

func NewNodeEq() *NodeEq {
	return &NodeEq{}
}

func (no *NodeEq) Consume(arg Node) {
	no.args = no.args.Push(arg)
}

func (no *NodeEq) Consumes() int {
	return 2
}

func (no *NodeEq) Produces() int {
	return 1
}

func (no *NodeEq) Precedence() int {
	return 4
}

func (no *NodeEq) Format() string {
	return fmt.Sprintf("Bool{eq(one(%s),one(%s))}", no.args[1].Format(), no.args[0].Format())
}

func (no *NodeEq) String() string {
	return fmt.Sprintf("NodeEq{%s}", no.args)
}
