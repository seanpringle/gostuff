package main

type NodeNameAgg struct {
	*NodeName
}

func NewNodeNameAgg(name string) *NodeNameAgg {
	return &NodeNameAgg{NewNodeName(name)}
}
