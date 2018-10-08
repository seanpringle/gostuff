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
	args := ""
	if ne.args != nil {
		args = ne.args.Format()
	}
	if m, is := ne.name.(*NodeMethod); is {
		return fmt.Sprintf("func() Any { t, m := %s; return call(m, join(t, %s)); }()", m.Format(), args)
	}
	return fmt.Sprintf("call(%s, join(%s))", ne.name.Format(), args)
}

func (ne *NodeExec) String() string {
	return fmt.Sprintf("NodeExec{%s(%s)}", ne.name, ne.args)
}
