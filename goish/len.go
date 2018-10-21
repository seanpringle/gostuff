package main

import (
	"fmt"
)

type NodeLen struct {
	item Node
}

func NewNodeLen() *NodeLen {
	return &NodeLen{}
}

func (ni *NodeLen) Consume(n Node) {
	ni.item = n
}

func (ni *NodeLen) Consumes() int {
	return 1
}

func (ni *NodeLen) Produces() int {
	return 1
}

func (ni *NodeLen) Precedence() int {
	return 5
}

func (ni *NodeLen) Format() string {
	return fmt.Sprintf("length(%s)", ni.item.Format())
}

func (ni *NodeLen) String() string {
	return fmt.Sprintf("NodeLen{%s}", ni.item)
}
