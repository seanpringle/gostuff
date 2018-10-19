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
	flag := FormatBool(nf.flag)
	if nf.onfalse != nil {
		return fmt.Sprintf("if %s { %s } else { %s }", flag, nf.ontrue.Format(), nf.onfalse.Format())
	}
	return fmt.Sprintf("if %s { %s }", flag, nf.ontrue.Format())
}

func (nf *NodeIf) String() string {
	return fmt.Sprintf("NodeIf{%s:%s:%s}", nf.flag, nf.ontrue, nf.onfalse)
}
