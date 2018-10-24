package main

import (
	"fmt"
)

type NodeUntil struct {
	flag Node
	body Node
}

func NewNodeUntil(flag, body Node) *NodeUntil {
	return &NodeUntil{
		flag: flag,
		body: body,
	}
}

func (nf *NodeUntil) Format() string {
	return fmt.Sprintf("loop(func() { for !(%s) { vm.da(call(vm, %s, nil)) } })", FormatBool(nf.flag), nf.body.Format())
}

func (nf *NodeUntil) String() string {
	return fmt.Sprintf("NodeUntil{%s:%s}", nf.flag, nf.body)
}
