package main

import (
	"fmt"
)

type NodeCatch struct {
	fn Node
}

func NewNodeCatch(node Node) *NodeCatch {
	return &NodeCatch{node}
}

func (nd *NodeCatch) Format() string {
	return fmt.Sprintf("defer catch(vm, %s)", nd.fn.Format())
}

func (nd *NodeCatch) String() string {
	return fmt.Sprintf("NodeCatch{%s}", nd.fn)
}
