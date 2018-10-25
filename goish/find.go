package main

import (
	"fmt"
)

type NodeFind struct {
	args Nodes
}

func NewNodeFind() *NodeFind {
	return &NodeFind{}
}

func (no *NodeFind) Consume(arg Node) {
	no.args = no.args.Push(arg)
}

func (no *NodeFind) Consumes() int {
	return 2
}

func (no *NodeFind) Produces() int {
	return 1
}

func (no *NodeFind) Precedence() int {
	return 7
}

func (no *NodeFind) Format() string {
	return fmt.Sprintf("find(%s,%s)", FormatOne(no.args[1]), no.args[0].(Keyer).FormatKey())
}

func (no *NodeFind) FormatStore(val string) string {
	return fmt.Sprintf("store(%s,%s,%s)", FormatOne(no.args[1]), no.args[0].(Keyer).FormatKey(), val)
}

func (no *NodeFind) String() string {
	return fmt.Sprintf("NodeFind{%s}", no.args)
}
