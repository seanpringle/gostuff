package main

import (
	"fmt"
)

type NodeTry struct {
	fn Node
}

func NewNodeTry(node Node) *NodeTry {
	return &NodeTry{node}
}

func (nd *NodeTry) Format() string {
	return fmt.Sprintf("try(vm, %s)", nd.fn.Format())
}

func (nd *NodeTry) String() string {
	return fmt.Sprintf("NodeTry{%s}", nd.fn)
}
