package main

import (
	"fmt"
)

type NodeWhile struct {
	flag Node
	body Node
}

func NewNodeWhile(flag, body Node) *NodeWhile {
	return &NodeWhile{
		flag: flag,
		body: body,
	}
}

func (nf *NodeWhile) Format() string {
	return fmt.Sprintf("loop(func() { for (%s) { vm.da(call(vm, %s, nil)) } })", FormatBool(nf.flag), nf.body.Format())
}

func (nf *NodeWhile) String() string {
	return fmt.Sprintf("NodeWhile{%s:%s}", nf.flag, nf.body)
}
