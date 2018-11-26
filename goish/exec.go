package main

import (
	"fmt"
)

type NodeExec struct {
	name Node
	args Node
}

func NewNodeExec(args Node) *NodeExec {
	return &NodeExec{args: args}
}

func (ne *NodeExec) Consume(n Node) {
	ne.name = n
}

func (ne *NodeExec) Consumes() int {
	return 1
}

func (ne *NodeExec) Produces() int {
	return 1
}

func (ne *NodeExec) Precedence() int {
	return 7
}

func (ne *NodeExec) Format() string {
	args := ne.args
	if args == nil {
		args = NewNodeLitNil()
	}
	if m, is := ne.name.(*NodeMethod); is {
		return fmt.Sprintf("func() *Args { t, m := %s;\nreturn call(vm, m, join(vm, t, %s)); }()", m.Format(), args.Format())
	}
	return fmt.Sprintf("call(vm, %s, %s)", ne.name.Format(), FormatJoin(args))
}

func (ne *NodeExec) FormatJoin() string {
	return ne.Format()
}

func (ne *NodeExec) String() string {
	return fmt.Sprintf("NodeExec{%s(%s)}", ne.name, ne.args)
}
