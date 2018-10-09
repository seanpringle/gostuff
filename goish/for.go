package main

import (
	"fmt"
)

type NodeFor struct {
	flag Node
	body Node
}

func NewNodeFor(flag, body Node) *NodeFor {
	return &NodeFor{
		flag: flag,
		body: body,
	}
}

func (nf *NodeFor) Format() string {
	return fmt.Sprintf("for truth(%s) { %s }", nf.flag.Format(), nf.body.Format())
}

func (nf *NodeFor) String() string {
	return fmt.Sprintf("NodeFor{%s:%s}", nf.flag, nf.body)
}

type NodeFor2 struct {
	begin Node
	step  Node
	body  Node
}

func NewNodeFor2(begin, step, body Node) *NodeFor2 {
	return &NodeFor2{begin, step, body}
}

func (nf *NodeFor2) Format() string {
	return fmt.Sprintf("for %s; truth(%s); { %s }",
		nf.begin.Format(),
		nf.step.Format(),
		nf.body.Format(),
	)
}

func (nf *NodeFor2) String() string {
	return fmt.Sprintf("NodeFor2{%s;%s;%s}", nf.begin, nf.step, nf.body)
}

type NodeFor3 struct {
	begin Node
	check Node
	step  Node
	body  Node
}

func NewNodeFor3(begin, check, step, body Node) *NodeFor3 {
	return &NodeFor3{begin, check, step, body}
}

func (nf *NodeFor3) Format() string {
	return fmt.Sprintf("for %s; truth(%s); %s { %s }",
		nf.begin.Format(),
		nf.check.Format(),
		nf.step.Format(),
		nf.body.Format(),
	)
}

func (nf *NodeFor3) String() string {
	return fmt.Sprintf("NodeFor3{%s;%s;%s;%s}", nf.begin, nf.check, nf.step, nf.body)
}
