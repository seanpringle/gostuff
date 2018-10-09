package main

import (
	"fmt"
)

type NodeElse struct {
	body Node
}

func NewNodeElse(body Node) *NodeElse {
	return &NodeElse{
		body: body,
	}
}

func (nf *NodeElse) Format() string {
	return fmt.Sprintf("else { %s }", nf.body.Format())
}

func (nf *NodeElse) String() string {
	return fmt.Sprintf("NodeElse{%s}", nf.body)
}
