package main

import (
	"fmt"
)

type NodeExtract struct {
	arg Node
}

func NewNodeExtract() *NodeExtract {
	return &NodeExtract{}
}

func (no *NodeExtract) Consume(arg Node) {
	no.arg = arg
}

func (no *NodeExtract) Consumes() int {
	return 1
}

func (no *NodeExtract) Produces() int {
	return 1
}

func (no *NodeExtract) Precedence() int {
	return 3
}

func (no *NodeExtract) Format() string {
	return fmt.Sprintf("extract(vm, %s)", FormatOne(no.arg))
}

func (no *NodeExtract) FormatOne() string {
	return FormatOne(no.arg)
}

func (no *NodeExtract) String() string {
	return fmt.Sprintf("NodeExtract{%s}", no.arg)
}
