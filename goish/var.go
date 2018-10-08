package main

import (
	"fmt"
)

type NodeVar struct {
	name *NodeName
}

func NewNodeVar(name *NodeName) *NodeVar {
	return &NodeVar{
		name: name,
	}
}

func (nd *NodeVar) Format() string {
	return fmt.Sprintf("var %s Any", nd.name.Format())
}

func (nd *NodeVar) String() string {
	return fmt.Sprintf("NodeVar{%s}", nd.name)
}

func (nd *NodeVar) Produces() int {
	return 1
}
