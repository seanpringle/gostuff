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
	if nr.node == nil {
		return fmt.Sprintf("return join(vm, nil)")
	}
	return fmt.Sprintf("return %s", FormatJoin(nr.node))
}

func (nr *NodeReturn) String() string {
	return fmt.Sprintf("NodeReturn{%s}", nr.node)
}
