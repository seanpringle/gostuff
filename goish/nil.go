package main

type NodeLitNil struct {
}

func NewNodeLitNil() *NodeLitNil {
	return &NodeLitNil{}
}

func (li *NodeLitNil) Format() string {
	return "nil"
}

func (li *NodeLitNil) FormatOne() string {
	return li.Format()
}

func (li *NodeLitNil) String() string {
	return "NodeLitNil{}"
}

func (li *NodeLitNil) Produces() int {
	return 1
}
