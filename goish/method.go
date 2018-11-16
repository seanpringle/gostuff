package main

import (
	"fmt"
)

type NodeMethod struct {
	args Nodes
}

func NewNodeMethod() *NodeMethod {
	return &NodeMethod{}
}

func (no *NodeMethod) Consume(arg Node) {
	no.args = no.args.Push(arg)
}

func (no *NodeMethod) Consumes() int {
	return 2
}

func (no *NodeMethod) Produces() int {
	return 1
}

func (no *NodeMethod) Precedence() int {
	return 7
}

func (no *NodeMethod) Format() string {
	//	if ex, is := no.args[1].(*NodeExec); is {
	return fmt.Sprintf("method(%s,%s)", FormatOne(no.args[1]), no.args[0].(Keyer).FormatKey())
	//	}
	//	return fmt.Sprintf("method(%s,%s)", no.args[1].Format(), no.args[0].(Keyer).FormatKey())
}

func (no *NodeMethod) String() string {
	return fmt.Sprintf("NodeMethod{%s}", no.args)
}
