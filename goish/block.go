package main

import (
	"fmt"
	"strings"
)

type NodeBlock struct {
	expr   Nodes
	scope  Scope
	parent *NodeBlock
}

func NewNodeBlock(parent *NodeBlock, scope Scope) *NodeBlock {
	return &NodeBlock{
		scope:  scope,
		parent: parent,
	}
}

func (nb *NodeBlock) Consume(n Node) {
	nb.expr = nb.expr.Push(n)
}

func (nb *NodeBlock) Produces() int {
	return 1
}

func (nb *NodeBlock) Define(v *NodeVar) {
	block := nb
	scope := nb.scope
	for scope == nil && block.parent != nil {
		block = block.parent
		scope = block.scope
	}
	if scope == nil {
		panic("missing scope")
	}
	scope[v.Format()] = v
}

func (nb *NodeBlock) Format() string {
	//if nb.scope == nil {
	//	return nb.expr.FormatJoin("\n")
	//}
	scope := []string{}
	for _, n := range nb.scope {
		scope = append(scope, n.Format())
	}
	return fmt.Sprintf("block(func() { %s; %s })",
		strings.Join(scope, "\n"),
		nb.expr.FormatJoin("\n"),
	)
}

func (nb *NodeBlock) String() string {
	scope := []string{}
	for _, n := range nb.scope {
		scope = append(scope, n.String())
	}
	return fmt.Sprintf("NodeBlock{%s:%s}", scope, nb.expr)
}
