package main

type NodeBreak struct {
}

func NewNodeBreak() *NodeBreak {
	return &NodeBreak{}
}

func (nr *NodeBreak) Format() string {
	return "panic(loopBreak(0))"
}

func (nr *NodeBreak) String() string {
	return "NodeBreak{}"
}
