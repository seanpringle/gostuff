package main

import (
	"fmt"
)

type NodeLoop struct {
	body Node
}

func (nl *NodeLoop) Format() string {
	return fmt.Sprintf("loop(func() Tup { %s; return Tup{nil} })", nl.body.Format())
}

func (nl *NodeLoop) String() string {
	return fmt.Sprintf("NodeLoop{%s}", nl.body)
}

type NodeFor struct {
	item Node
	body Node
}

func NewNodeFor(item, body Node) *NodeFor {
	return &NodeFor{
		item: item,
		body: body,
	}
}

func (nf *NodeFor) Format() string {
	return fmt.Sprintf("loop(func() { it := iterate(%s); for { aa := it(vm, nil); if aa.get(0) == nil { vm.da(aa); break }; vm.da(call(vm, %s, aa)); } })",
		FormatOne(nf.item), nf.body.Format(),
	)
}

func (nf *NodeFor) String() string {
	return fmt.Sprintf("NodeFor{%s:%s}", nf.item, nf.body)
}
