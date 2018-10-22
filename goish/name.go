package main

import (
	"fmt"
)

var KId int = 1
var Keys map[string]int = map[string]int{}

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

func (nn *NodeName) FormatOne() string {
	return fmt.Sprintf(`N%s`, nn.name)
}

func (nn *NodeName) String() string {
	return fmt.Sprintf("NodeName{%s}", nn.name)
}

func (nn *NodeName) Produces() int {
	return 1
}

func (nn *NodeName) FormatKey() string {
	//return fmt.Sprintf(`Str{%q}`, nn.name)
	if _, ok := Keys[nn.name]; !ok {
		Keys[nn.name] = KId
		KId++
	}
	return fmt.Sprintf(`S%d /* %s */`, Keys[nn.name], nn.name)
}
