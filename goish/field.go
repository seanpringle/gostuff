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
	return 7
}

func (no *NodeField) Format() string {
	return fmt.Sprintf("field(%s,%s)", FormatOne(no.args[1]), FormatOne(no.args[0]))
}

func (no *NodeField) FormatStore(val string) string {
	return fmt.Sprintf("store(%s,%s,%s)", FormatOne(no.args[1]), FormatOne(no.args[0]), val)
}

func (no *NodeField) String() string {
	return fmt.Sprintf("NodeField{%s}", no.args)
}
