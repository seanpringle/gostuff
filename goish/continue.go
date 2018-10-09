package main

type NodeContinue struct {
}

func NewNodeContinue() *NodeContinue {
	return &NodeContinue{}
}

func (nr *NodeContinue) Format() string {
	return "continue"
}

func (nr *NodeContinue) String() string {
	return "NodeContinue{}"
}
