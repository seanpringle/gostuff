package main

import (
	"fmt"
)

type NodeReturn struct {
	node Node
}

func NewNodeReturn(node Node) *NodeReturn {
	return &NodeReturn{node: node}
}

func (nr *NodeReturn) Format() string {
	return fmt.Sprintf("return join(%s)", nr.node.Format())
}

func (nr *NodeReturn) String() string {
	return fmt.Sprintf("NodeReturn{%s}", nr.node)
}
