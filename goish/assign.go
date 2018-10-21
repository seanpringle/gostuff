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

	if len(na.targets) == 1 && len(na.sources) == 1 {
		if find, is := na.targets[0].(*NodeFind); is {
			return fmt.Sprintf("func() Any { a := %s; %s; return a }()",
				FormatOne(na.sources),
				find.FormatStore("a"),
			)
		} else if find, is := na.targets[0].(*NodeField); is {
			return fmt.Sprintf("func() Any { a := %s; %s; return a }()",
				FormatOne(na.sources),
				find.FormatStore("a"),
			)
		} else {
			return fmt.Sprintf("func() Any { a := %s; %s = a; return a }()",
				FormatOne(na.sources),
				na.targets.Format(),
			)
		}
	}

	assigns := []string{}
	for i, src := range na.targets {
		if find, is := src.(*NodeFind); is {
			assigns = append(assigns, find.FormatStore(fmt.Sprintf("aa.get(%d)", i)))
		} else if find, is := src.(*NodeField); is {
			assigns = append(assigns, find.FormatStore(fmt.Sprintf("aa.get(%d)", i)))
		} else {
			assigns = append(assigns, fmt.Sprintf("%s = aa.get(%d)", src.Format(), i))
		}
	}

	return fmt.Sprintf("func() *Args { aa := join(vm, %s); %s; return aa }()",
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
