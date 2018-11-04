package main

import (
	"fmt"
	"strings"
)

type NodeRec struct {
	table map[*NodeName]Node
}

func NewNodeRec(t map[*NodeName]Node) *NodeRec {
	return &NodeRec{
		table: t,
	}
}

func (nt *NodeRec) Format() string {
	keys := []string{}
	pairs := []string{}
	for k, v := range nt.table {
		keys = append(keys, fmt.Sprintf("%s Any", k.Format()))
		pairs = append(pairs, fmt.Sprintf("%s: %s", k.Format(), FormatOne(v)))
	}
	return fmt.Sprintf(`
		func() Any {
			type R struct {
				Record
				%s
			}
			return &R{
				Record: Record{},
				%s,
			}
		}()`,
		strings.Join(keys, "\n"),
		strings.Join(pairs, ",\n"),
	)
}

func (nt *NodeRec) String() string {
	return fmt.Sprintf("NodeRec{%s}", nt.table)
}

type NodeRecField struct {
	name *NodeName
}

func (nt *NodeRecField) Format() string {
	return fmt.Sprintf(".%s", nt.name.Format())
}

func (nt *NodeRecField) String() string {
	return fmt.Sprintf("NodeRecField{%s}", nt.name)
}
