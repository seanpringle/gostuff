package main

import (
	"fmt"
)

type NodeIf struct {
	flag    Node
	ontrue  Node
	onfalse Node
}

func NewNodeIf(flag, ontrue, onfalse Node) *NodeIf {
	return &NodeIf{
		flag:    flag,
		ontrue:  ontrue,
		onfalse: onfalse,
	}
}

func (nf *NodeIf) Format() string {
	if nf.onfalse != nil {
		return fmt.Sprintf("if truth(%s) { %s } else { %s }", nf.flag.Format(), nf.ontrue.Format(), nf.onfalse.Format())

	}
	return fmt.Sprintf("if truth(%s) { %s }", nf.flag.Format(), nf.ontrue.Format())
}

func (nf *NodeIf) String() string {
	return fmt.Sprintf("NodeIf{%s:%s:%s}", nf.flag, nf.ontrue, nf.onfalse)
}
