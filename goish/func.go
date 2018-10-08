package main

import (
	"fmt"
	"strings"
)

type NodeFunc struct {
	args Nodes
	body Node
}

func NewNodeFunc(args Nodes, body Node) *NodeFunc {
	return &NodeFunc{
		args: args,
		body: body,
	}
}

func (nf *NodeFunc) Format() string {

	init := []string{}
	for i, arg := range nf.args {
		init = append(init, fmt.Sprintf("%s := get(aa, %d)", arg.(*NodeName).Format(), i))
	}

	return fmt.Sprintf("Func(func(aa Tup) Tup { %s; return %s })",
		strings.Join(init, ";"),
		nf.body.Format(),
	)
}

func (nf *NodeFunc) String() string {
	return fmt.Sprintf("NodeFunc{%s:%s}", Nodes(nf.args), nf.body)
}
