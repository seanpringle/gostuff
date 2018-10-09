package main

import (
	"fmt"
)

type NodeIf struct {
	flag Node
	body Node
}

func NewNodeIf(flag, body Node) *NodeIf {
	return &NodeIf{
		flag: flag,
		body: body,
	}
}

func (nf *NodeIf) Format() string {
	return fmt.Sprintf("if truth(%s) { %s }", nf.flag.Format(), nf.body.Format())
}

func (nf *NodeIf) String() string {
	return fmt.Sprintf("NodeIf{%s:%s}", nf.flag, nf.body)
}
