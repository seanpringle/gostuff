package main

import (
	"fmt"
)

type NodeTuple Nodes

func NewNodeTuple(nodes Nodes) NodeTuple {
	return NodeTuple(nodes)
}

func (nt NodeTuple) Format() string {
	return Nodes(nt).FormatJoin(",")
}

func (nt NodeTuple) FormatJoin() string {
	if len(nt) == 1 {
		return fmt.Sprintf("Tup{%s}", FormatOne(nt[0]))
	}
	return fmt.Sprintf("join(%s)", nt.Format())
}

func (nt NodeTuple) String() string {
	return fmt.Sprintf("NodeTuple{%s}", Nodes(nt))
}

func (nt NodeTuple) Produces() int {
	return len(nt)
}
