package main

import (
	"fmt"
)

type NodeDefer struct {
	fn Node
}

func NewNodeDefer(node Node) *NodeDefer {
	return &NodeDefer{node}
}

func (nd *NodeDefer) Format() string {
	return fmt.Sprintf("defer func() { %s }()", nd.fn.Format())
}

func (nd *NodeDefer) String() string {
	return fmt.Sprintf("NodeDefer{%s}", nd.fn)
}
