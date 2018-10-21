package main

import (
	"fmt"
	"strings"
)

type NodeMap struct {
	table map[Keyer]Node
}

func NewNodeMap(t map[Keyer]Node) *NodeMap {
	return &NodeMap{
		table: t,
	}
}

func (nt *NodeMap) Format() string {
	pairs := []string{}
	for k, v := range nt.table {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k.FormatKey(), FormatOne(v)))
	}
	return fmt.Sprintf("NewMap(MapData{%s})", strings.Join(pairs, ","))
}

func (nt *NodeMap) String() string {
	return fmt.Sprintf("NodeMap{%s}", nt.table)
}
