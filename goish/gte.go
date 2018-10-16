package main

import (
	"fmt"
)

type NodeGte struct {
	args Nodes
}

func NewNodeGte() *NodeGte {
	return &NodeGte{}
}

func (no *NodeGte) Consume(arg Node) {
	no.args = no.args.Push(arg)
}

func (no *NodeGte) Consumes() int {
	return 2
}

func (no *NodeGte) Produces() int {
	return 1
}

func (no *NodeGte) Precedence() int {
	return 4
}

func (no *NodeGte) Format() string {
	return fmt.Sprintf("Bool(gte(one(%s),one(%s)))", no.args[1].Format(), no.args[0].Format())
}

func (no *NodeGte) String() string {
	return fmt.Sprintf("NodeGte{%s}", no.args)
}
