package main

import (
	"fmt"
	"strings"
)

type NodeAssign struct {
	targets NodeTuple
	sources NodeTuple
}

func NewNodeAssign(block *NodeBlock, targets, sources NodeTuple) *NodeAssign {

	for _, target := range targets {
		if name, is := target.(*NodeName); is {
			block.Define(NewNodeVar(name))
		}
	}

	return &NodeAssign{targets: targets, sources: sources}
}

func (na *NodeAssign) Format() string {

	assigns := []string{}
	for i, src := range na.targets {
		if find, is := src.(*NodeFind); is {
			assigns = append(assigns, find.FormatStore(fmt.Sprintf("get(aa, %d)", i)))
		} else {
			assigns = append(assigns, fmt.Sprintf("%s = get(aa, %d)", src.Format(), i))
		}
	}

	return fmt.Sprintf("func() Tup { aa := join(%s); %s; return aa }()",
		na.sources.Format(),
		strings.Join(assigns, ";"),
	)
}

func (na *NodeAssign) String() string {
	return fmt.Sprintf("NodeAssign{%s,%s}", na.targets, na.sources)
}

func (na *NodeAssign) Produces() int {
	return len(na.sources)
}
