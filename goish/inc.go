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
	return 5
}

func (ni *NodeInc) Format() string {
	loc := ni.item.Format()
	if _, is := ni.item.(*NodeName); is {
		return fmt.Sprintf("func() Any { v := one(%s); %s = add(v, Int{1}); return v; }()", loc, loc)
	}
	panic("post-increment only supported on local variables")
}

func (ni *NodeInc) String() string {
	return fmt.Sprintf("NodeInc{%s}", ni.item)
}
