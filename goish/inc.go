package main

import (
	"fmt"
)

type NodeInc struct {
	item Node
}

func NewNodeInc() *NodeInc {
	return &NodeInc{}
}

func (ni *NodeInc) Consume(n Node) {
	ni.item = n
}

func (ni *NodeInc) Consumes() int {
	return 1
}

func (ni *NodeInc) Produces() int {
	return 1
}

func (ni *NodeInc) Precedence() int {
	return 9
}

func (ni *NodeInc) Format() string {
	loc := ni.item.Format()
	if _, is := ni.item.(*NodeName); is {
		return fmt.Sprintf("func() Any { v := one(vm, %s); %s = add(v, Int(1)); return v; }()", loc, loc)
	}
	panic(fmt.Sprintf("post-increment only supported on local variables: %v", ni.item.Format()))
}

func (ni *NodeInc) String() string {
	return fmt.Sprintf("NodeInc{%s}", ni.item)
}
