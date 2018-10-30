package main

import (
	"fmt"
	"strings"
)

type NodeList struct {
	nodes Nodes
}

func NewNodeList(nodes Nodes) *NodeList {
	return &NodeList{nodes: nodes}
}

func (ns *NodeList) Format() string {
	s := []string{}
	for _, n := range ns.nodes {
		s = append(s, FormatOne(n))
	}
	return fmt.Sprintf("NewList([]Any{%s})", strings.Join(s, ","))
}

func (ns *NodeList) String() string {
	return fmt.Sprintf("NodeList{%s}", ns.nodes)
}
