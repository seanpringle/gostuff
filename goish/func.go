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
		if agg, is := arg.(*NodeNameAgg); is {
			init = append(init, fmt.Sprintf("%s := aa.agg(%d); noop(%s)", agg.Format(), i, agg.Format()))
			break
		}
		init = append(init, fmt.Sprintf("%s := aa.get(%d); noop(%s)", arg.(*NodeName).Format(), i, arg.(*NodeName).Format()))
	}

	if b, is := nf.body.(*NodeBlock); is {
		for _, arg := range nf.args {
			if b.scope != nil {
				delete(b.scope, arg.Format())
			}
		}
	}

	//return fmt.Sprintf("Func(func(aa Tup) Tup { %s; return func() Tup { %s; return Tup{nil}; }() })",
	return fmt.Sprintf("Func(func(vm *VM, aa *Args) *Args { %s; vm.da(aa); %s; return nil })",
		strings.Join(init, ";"),
		nf.body.Format(),
	)
}

func (nf *NodeFunc) String() string {
	return fmt.Sprintf("NodeFunc{%s:%s}", Nodes(nf.args), nf.body)
}
