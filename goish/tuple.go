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

func (nt NodeTuple) FormatOne() string {
	if len(nt) == 1 {
		return FormatOne(nt[0])
	}
	return fmt.Sprintf("one(%s)", nt.Format())
}

func (nt NodeTuple) FormatBool() string {
	if len(nt) == 1 {
		return FormatBool(nt[0])
	}
	return nt.Format()
}

func (nt NodeTuple) String() string {
	return fmt.Sprintf("NodeTuple{%s}", Nodes(nt))
}

func (nt NodeTuple) Produces() int {
	return len(nt)
}
