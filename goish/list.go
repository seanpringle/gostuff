package main

import (
	"fmt"
)

type NodeList struct {
	nodes Nodes
}

func NewNodeList(nodes Nodes) *NodeList {
	return &NodeList{nodes: nodes}
}

func (ns *NodeList) Format() string {
	return fmt.Sprintf("NewList([]Any{%s})", ns.nodes.FormatJoin(","))
}

func (ns *NodeList) String() string {
	return fmt.Sprintf("NodeList{%s}", ns.nodes)
}
