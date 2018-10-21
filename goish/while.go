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
	return fmt.Sprintf("loop(func() { for truth(%s) { vm.da(call(vm, %s, nil)) } })", nf.flag.Format(), nf.body.Format())
}

func (nf *NodeWhile) String() string {
	return fmt.Sprintf("NodeWhile{%s:%s}", nf.flag, nf.body)
}
