package main

type NodeBreak struct {
}

func NewNodeBreak() *NodeBreak {
	return &NodeBreak{}
}

func (nr *NodeBreak) Format() string {
	return "break"
}

func (nr *NodeBreak) String() string {
	return "NodeBreak{}"
}
