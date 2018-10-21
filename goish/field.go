package main

import (
	"fmt"
)

type NodeField struct {
	args Nodes
}

func NewNodeField() *NodeField {
	return &NodeField{}
}

func (no *NodeField) Consume(arg Node) {
	no.args = no.args.Push(arg)
}

func (no *NodeField) Consumes() int {
	return 2
}

func (no *NodeField) Produces() int {
	return 1
}

func (no *NodeField) Precedence() int {
	return 8
}

func (no *NodeField) Format() string {
	return fmt.Sprintf("field(%s,%s)", no.args[1].Format(), no.args[0].Format())
}

func (no *NodeField) FormatStore(val string) string {
	return fmt.Sprintf("store(%s,%s,%s)", no.args[1].Format(), no.args[0].Format(), val)
}

func (no *NodeField) String() string {
	return fmt.Sprintf("NodeField{%s}", no.args)
}
