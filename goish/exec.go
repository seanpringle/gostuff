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
		args = NewNodeLitInt(0)
	}
	if m, is := ne.name.(*NodeMethod); is {
		return fmt.Sprintf("func() Tup { t, m := %s; return call(m, join(t, %s)); }()", m.Format(), args.Format())
	}
	//if p, is := ne.args.(Producer); is {
	//	if p.Produces() == 1 {
	//		return fmt.Sprintf("call(%s, %s)", ne.name.Format(), args)
	//	}
	//}
	if f, is := ne.name.(*NodeName); is {
		return fmt.Sprintf("%s.(Func)(%s)", f.Format(), FormatJoin(args))
	}
	return fmt.Sprintf("call(%s, %s)", ne.name.Format(), FormatJoin(args))
}

func (ne *NodeExec) FormatOne() string {
	return fmt.Sprintf("%s[0]", ne.Format())
}

func (ne *NodeExec) String() string {
	return fmt.Sprintf("NodeExec{%s(%s)}", ne.name, ne.args)
}
