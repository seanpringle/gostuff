package main

import (
	"fmt"
)

type NodeName struct {
	name string
}

func NewNodeName(name string) *NodeName {
	return &NodeName{
		name: name,
	}
}

func (nn *NodeName) Format() string {
	return fmt.Sprintf(`N%s`, nn.name)
}

func (nn *NodeName) String() string {
	return fmt.Sprintf("NodeName{%s}", nn.name)
}

func (nn *NodeName) Produces() int {
	return 1
}

func (nn *NodeName) FormatKey() string {
	return fmt.Sprintf(`Str{%q}`, nn.name)
}
