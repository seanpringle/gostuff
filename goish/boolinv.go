package main

import (
	"fmt"
)

type NodeBoolInv struct {
	arg Node
}

func NewNodeBoolInv(n Node) *NodeBoolInv {
	return &NodeBoolInv{n}
}

func (na *NodeBoolInv) Format() string {
	return fmt.Sprintf("b_inv(%s)", FormatOne(na.arg))
}

func (na *NodeBoolInv) String() string {
	return fmt.Sprintf("NodeBoolInv{%s}", na.arg)
}
