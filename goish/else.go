package main

import (
	"fmt"
)

type NodeElse struct {
	flag Node
	body Node
}

func NewNodeElse(flag, body Node) *NodeElse {
	return &NodeElse{
		flag: flag,
		body: body,
	}
}

func (nf *NodeElse) Format() string {
	return fmt.Sprintf("else { %s }", nf.flag.Format(), nf.body.Format())
}

func (nf *NodeElse) String() string {
	return fmt.Sprintf("NodeElse{%s:%s}", nf.flag, nf.body)
}
