package main

import (
	"bufio"
	//"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func ensure(b bool, msg string) {
	if !b {
		panic(errors.New(msg))
	}
}

func parseInt(str string) int64 {
	i64, err := strconv.ParseInt(str, 10, 64)
	assert(err)
	return i64
}

func parseDec(str string) float64 {
	f64, err := strconv.ParseFloat(str, 64)
	assert(err)
	return f64
}

type Node interface {
	Format() string
	String() string
}

type Consumer interface {
	Consume(Node)
	Consumes() int
}

type Producer interface {
	Produces() int
}

type Operator interface {
	Consumer
	Producer
	Precedence() int
}

type Keyer interface {
	Node
	FormatKey() string
}

type CanFormatOne interface {
	FormatOne() string
}

func FormatOne(n Node) string {
	if f, is := n.(CanFormatOne); is {
		return f.FormatOne()
	}
	return fmt.Sprintf("one(vm, %s)", n.Format())
}

type CanFormatJoin interface {
	FormatJoin() string
}

func FormatJoin(n Node) string {
	if f, is := n.(CanFormatJoin); is {
		return f.FormatJoin()
	}
	return fmt.Sprintf("join(vm, %s)", n.Format())
}

type CanFormatBool interface {
	FormatBool() string
}

func FormatBool(n Node) string {
	if f, is := n.(CanFormatBool); is {
		return f.FormatBool()
	}
	return fmt.Sprintf("truth(%s)", FormatOne(n))
}

type Nodes []Node
type Scope map[string]Node

func (nl Nodes) Push(n Node) Nodes {
	nl = append(nl, n)
	return nl
}

func (nl Nodes) Pop() (Node, Nodes) {
	l := len(nl) - 1
	n := nl[l]
	nl = nl[:l]
	return n, nl
}

func (nl Nodes) Last() Node {
	if len(nl) > 0 {
		return nl[len(nl)-1]
	}
	return nil
}

func (nl Nodes) FormatJoin(sep string) string {
	parts := []string{}
	for _, n := range nl {
		parts = append(parts, n.Format())
	}
	return strings.Join(parts, sep)
}

func (nl Nodes) String() string {
	parts := []string{}
	for _, n := range nl {
		parts = append(parts, n.String())
	}
	return strings.Join(parts, ",")
}

type Parser struct {
	src   *bufio.Reader
	dst   *bufio.Writer
	queue []rune
}

func (p *Parser) output() *bufio.Writer {
	return p.dst
}

func (p *Parser) read() bool {

	c, _, err := p.src.ReadRune()

	if err != nil && err != io.EOF {
		panic(err)
	}

	if err != io.EOF && c != unicode.ReplacementChar {
		p.queue = append(p.queue, c)
		return true
	}

	return false
}

func (p *Parser) drop() bool {
	if len(p.queue) > 0 {
		p.queue = p.queue[1:]
		return true
	}
	return false
}

func (p *Parser) char(n int) rune {
	for len(p.queue) < n+1 && p.read() {
	}
	if len(p.queue) > n {
		return p.queue[n]
	}
	return rune(0)
}

func (p *Parser) scan() rune {
	for {
		c := p.char(0)

		if c == rune(0) {
			break
		}
		if p.iswhite(c) {
			p.drop()
			continue
		}
		break
	}
	return p.char(0)
}

func (p *Parser) next() rune {
	return p.char(0)
}

func (p *Parser) take() rune {
	c := p.char(0)
	p.drop()
	return c
}

func (p *Parser) peek(s string) bool {
	p.scan()
	for i, r := range s {
		if p.char(i) != r {
			return false
		}
	}
	r := p.char(len(s))
	return p.iswhite(r) || p.issymbol(r) || r == rune(0)
}

func (p *Parser) terminator() bool {
	return p.scan() == rune(0) || p.scan() == ';' || p.scan() == ')' || p.scan() == '}' || p.scan() == ']' || p.peek("end")
}

func (p *Parser) operator() bool {
	c := p.scan()
	return c == '=' || c == '+' || c == '-' || c == '>' || c == '<' || p.peek("or") || p.peek("and")
}

func (p *Parser) isalpha(c rune) bool {
	return unicode.IsLetter(c)
}

func (p *Parser) isdigit(c rune) bool {
	return unicode.IsNumber(c)
}

func (p *Parser) isnumber(c rune) bool {
	return p.isdigit(c) || c == '.'
}

func (p *Parser) iswhite(c rune) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}

func (p *Parser) issymbol(c rune) bool {
	return (unicode.IsSymbol(c) || unicode.IsPunct(c)) && !p.iswhite(c)
}

func (p *Parser) isname(c rune) bool {
	return c != '.' && (p.isalpha(c) || p.isnumber(c) || c == '_')
}

func (p *Parser) iskeyword(w string) bool {
	return w == "and" || w == "or"
}

func (p *Parser) node(block *NodeBlock) Node {

	for p.scan() == '-' && p.char(1) == '-' {
		p.take()
		p.take()
		if p.char(0) == '[' && p.char(1) == '[' {
			p.node(nil)
		} else {
			for c := p.take(); c != rune(0) && c != '\n'; c = p.take() {
			}
		}
	}

	if p.terminator() {
		return nil
	}

	if p.scan() == '(' {
		p.take()
		node := NewNodeExec(p.tuple(block, nil))
		ensure(p.scan() == ')', "expected closing paren (exec)")
		p.take()
		return node
	}

	str := []rune{}

	if p.scan() == '.' && p.char(1) == '.' {
		p.take()
		p.take()
		return NewNodeConcat()
	}

	if p.scan() == '.' {
		p.take()
		return NewNodeFind()
	}

	if p.scan() == ':' {
		p.take()
		return NewNodeMethod()
	}

	if p.scan() == '\\' {
		p.take()
		return NewNodeField()
	}

	if p.scan() == '!' || p.peek("not") {
		if p.scan() == '!' {
			p.take()
		} else {
			p.take()
			p.take()
			p.take()
		}
		return NewNodeNot(p.node(block))
	}

	if p.scan() == '=' && p.char(1) == '=' {
		p.take()
		p.take()
		return NewNodeEq()
	}

	if p.scan() == '<' && p.char(1) == '=' {
		p.take()
		return NewNodeLte()
	}

	if p.scan() == '<' {
		p.take()
		return NewNodeLt()
	}

	if p.scan() == '>' && p.char(1) == '=' {
		p.take()
		return NewNodeGte()
	}

	if p.scan() == '>' {
		p.take()
		return NewNodeGt()
	}

	if p.peek("or") {
		p.take()
		p.take()
		return NewNodeOr()
	}

	if p.peek("and") {
		p.take()
		p.take()
		p.take()
		return NewNodeAnd()
	}

	if p.scan() == '+' && p.char(1) == '+' {
		p.take()
		p.take()
		return NewNodeInc()
	}

	if p.scan() == '+' {
		p.take()
		return NewNodeAdd()
	}

	if p.scan() == '-' {
		p.take()
		return NewNodeSub()
	}

	if p.scan() == '*' {
		p.take()
		return NewNodeMul()
	}

	if p.scan() == '/' {
		p.take()
		return NewNodeDiv()
	}

	if p.scan() == '%' {
		p.take()
		return NewNodeMod()
	}

	if p.scan() == '#' {
		p.take()
		return NewNodeLen()
	}

	if p.peek("return") {
		p.take()
		p.take()
		p.take()
		p.take()
		p.take()
		p.take()
		return NewNodeReturn(p.tuple(block, nil))
	}

	if p.peek("do") {
		p.take()
		p.take()
		var args Nodes
		if p.scan() == '(' {
			p.take()
			argTup := p.tuple(block, nil)
			if argTup != nil {
				args = Nodes(argTup.(NodeTuple))
			}
			ensure(p.scan() == ')', "expected closing paren (do)")
			p.take()
		}
		body := p.block(block, nil, nil).(*NodeBlock)
		ensure(p.peek("end"), "expected: end (do)")
		p.take()
		p.take()
		p.take()
		return NewNodeFunc(args, body)
	}

	if p.peek("function") {
		p.take()
		p.take()
		p.take()
		p.take()
		p.take()
		p.take()
		p.take()
		p.take()
		var args Nodes
		if p.scan() == '(' {
			p.take()
			argTup := p.tuple(block, nil)
			if argTup != nil {
				args = Nodes(argTup.(NodeTuple))
			}
			ensure(p.scan() == ')', "expected closing paren (func)")
			p.take()
		}
		body := p.block(block, Scope{}, nil).(*NodeBlock)
		ensure(p.peek("end"), "expected: end (function)")
		p.take()
		p.take()
		p.take()
		return NewNodeFunc(args, body)
	}

	if p.peek("if") {
		p.take()
		p.take()
		iter := p.tuple(block, []string{"then"})
		ensure(p.peek("then"), "expected: then")
		p.take()
		p.take()
		p.take()
		p.take()
		var onfalse Node
		ontrue := p.block(block, nil, []string{"else"})
		if p.peek("else") {
			p.take()
			p.take()
			p.take()
			p.take()
			onfalse = p.block(block, nil, nil)
		}
		ensure(p.peek("end"), "expected: end (if)")
		p.take()
		p.take()
		p.take()
		return NewNodeIf(iter, ontrue, onfalse)
	}

	if p.peek("for") {
		p.take()
		p.take()
		p.take()
		iter := p.tuple(block, nil)
		body := p.tuple(block, nil)
		return NewNodeFor(iter, body)
	}

	if p.peek("while") {
		p.take()
		p.take()
		p.take()
		p.take()
		p.take()
		cond := p.tuple(block, nil)
		body := p.tuple(block, nil)
		return NewNodeWhile(cond, body)
	}

	if p.peek("until") {
		p.take()
		p.take()
		p.take()
		p.take()
		p.take()
		cond := p.tuple(block, nil)
		body := p.tuple(block, nil)
		return NewNodeUntil(cond, body)
	}

	if p.peek("break") {
		p.take()
		p.take()
		p.take()
		p.take()
		p.take()
		return NewNodeBreak()
	}

	if p.peek("continue") {
		p.take()
		p.take()
		p.take()
		p.take()
		p.take()
		p.take()
		p.take()
		p.take()
		return NewNodeContinue()
	}

	if p.isnumber(p.scan()) {
		isDec := false
		for p.isnumber(p.next()) {
			c := p.take()
			str = append(str, c)
			isDec = isDec || c == '.'
		}
		if isDec {
			return NewNodeLitDec(parseDec(string(str)))
		}
		return NewNodeLitInt(parseInt(string(str)))
	}

	if p.scan() == '"' {
		str = append(str, p.take())
		for p.next() != '"' {
			c := p.take()
			str = append(str, c)
			if c == '\\' {
				str = append(str, p.take())
			}
		}
		ensure(p.next() == '"', "expected closing quotes")
		str = append(str, p.take())
		return NewNodeLitStr(string(str))
	}

	if p.scan() == '\'' {
		str = append(str, p.take())
		for p.next() != '\'' {
			c := p.take()
			str = append(str, c)
			if c == '\\' {
				str = append(str, p.take())
			}
		}
		ensure(p.next() == '\'', "expected closing quotes")
		str = append(str, p.take())
		return NewNodeLitRune(string(str))
	}

	if p.isname(p.scan()) {
		for p.isname(p.next()) {
			str = append(str, p.take())
		}
		return NewNodeName(string(str))
	}

	if p.scan() == '{' {
		p.take()
		t := map[Keyer]Node{}

		i := int64(0)
		for p.scan() != rune(0) && p.scan() != '}' {
			var key Keyer
			val := p.expression(block)

			if p.scan() == '=' {
				p.take()
				key = val.(Keyer)
				val = p.expression(block)
			}

			if key == nil {
				key = NewNodeLitInt(i)
				i++
			}

			t[key] = val

			if p.scan() == ',' {
				p.take()
				continue
			}

			break
		}
		ensure(p.scan() == '}', "expected closing brace (map)")
		p.take()

		return NewNodeMap(t)
	}

	if p.scan() == '[' && p.char(1) == '[' {
		p.take()
		p.take()
		str := ""
		for p.next() != rune(0) && !(p.next() == ']' && p.char(1) == ']') {
			c := p.take()
			if c == '`' {
				str = str + "` + \"`\" + `"
			} else {
				str = str + string(c)
			}
		}
		ensure(p.next() == ']' && p.char(1) == ']', "expected closing brackets ([[string]])")
		p.take()
		p.take()
		return NewNodeLitStr("`" + str + "`")
	}

	if p.scan() == '[' {
		p.take()
		var nodes Nodes
		t := p.tuple(block, nil)
		if t != nil {
			nodes = Nodes(t.(NodeTuple))
		}
		ensure(p.scan() == ']', "expected closing bracket (list)")
		p.take()
		return NewNodeList(nodes)
	}

	return nil
}

func (p *Parser) expression(block *NodeBlock) Node {

	if p.terminator() {
		return nil
	}

	items := Nodes{}
	ops := Nodes{}

	shunt := func(prec int) {
		var op Node
		for ops.Last() != nil && ops.Last().(Operator).Precedence() >= prec {
			op, ops = ops.Pop()
			var item Node
			for i := 0; i < op.(Consumer).Consumes() && len(items) > 0; i++ {
				item, items = items.Pop()
				op.(Consumer).Consume(item)
			}
			items = items.Push(op)
		}
	}

	var last Node

	for node := p.node(block); node != nil; node = p.node(block) {

		if ex, is := node.(*NodeExec); is {
			_, op := last.(Operator)
			if op || len(items) == 0 {
				node = ex.args
			}
		}

		last = node

		if op, is := node.(Operator); is {
			shunt(op.Precedence())
			ops = ops.Push(node)
		} else {
			items = items.Push(node)
		}

		consuming := 0
		for _, op := range ops {
			consuming += op.(Operator).Consumes()
		}

		if consuming > len(items)+len(ops)-1 {
			continue
		}

		//if p.scan() == ',' || p.scan() == '{' || p.scan() == '[' || p.terminator() {
		if p.scan() == ',' || p.terminator() {
			break
		}

		if p.scan() == '-' && p.char(1) == '-' {
			break
		}

		if !p.issymbol(p.scan()) && !p.peek("or") && !p.peek("and") {
			break
		}
	}

	for len(ops) > 0 {
		shunt(0)
	}

	if len(items) == 0 {
		return nil
	}

	if len(items) > 1 {
		panic(fmt.Sprintf("unbalanced expression %s", items))
	}

	return items[0]
}

func (p *Parser) tuple(block *NodeBlock, terms []string) Node {

	terminate := func() bool {
		if p.terminator() {
			return true
		}
		for _, term := range terms {
			if p.peek(term) {
				return true
			}
		}
		return false
	}

	if terminate() {
		return nil
	}

	args := Nodes{}

	for expr := p.expression(block); expr != nil; expr = p.expression(block) {
		args = append(args, expr)

		if p.scan() == ',' {
			p.take()
			if !terminate() {
				continue
			}
		}

		break
	}

	if p.scan() == '=' && p.char(1) != '=' {
		p.take()

		vars := args
		args = Nodes{}

		for expr := p.expression(block); expr != nil; expr = p.expression(block) {
			args = append(args, expr)

			if p.scan() == ',' {
				p.take()
				if !terminate() {
					continue
				}
			}

			break
		}

		return NewNodeAssign(block, NewNodeTuple(vars), NewNodeTuple(args))
	}

	if len(args) > 0 {
		return NewNodeTuple(args)
	}

	return nil
}

func (p *Parser) block(parent *NodeBlock, scope Scope, terms []string) Node {
	block := NewNodeBlock(parent, scope)
	for expr := p.tuple(block, terms); expr != nil; expr = p.tuple(block, terms) {
		block.Consume(expr)
	}
	return block
}

func (p *Parser) print(ss ...string) {
	for _, s := range ss {
		if _, err := p.output().WriteString(s); err != nil {
			panic(err)
		}
	}
}

func (p *Parser) println(ss ...string) {
	p.print(append(ss[:len(ss)], "\n")...)
}

func (p *Parser) run() (wtf error) {

	//	defer func() {
	//		if r := recover(); r != nil {
	//			switch r.(type) {
	//			case error:
	//				wtf = r.(error)
	//			default:
	//				wtf = fmt.Errorf("Parser.run(): %v", r)
	//			}
	//		}
	//	}()

	{
		data, err := ioutil.ReadFile("base.go.tpl")
		if err != nil {
			log.Fatalf("base.go.tpl %v", err)
		}
		p.println(string(data))
	}

	block := p.block(nil, Scope{}, nil).Format()

	for k, n := range Keys {
		p.println(fmt.Sprintf(`const S%d Text = Text(%q)`, n, k))
	}

	p.println(`func main() {`)

	p.println(`

		f, err := os.Create("cpuprofile")
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()

		for _, f := range onInit {
			f()
		}

		vm := &VM{}

	`)

	p.println(block)
	p.println(`}`)
	return
}

func main() {

	p := &Parser{
		src: bufio.NewReader(os.Stdin),
		dst: bufio.NewWriter(os.Stdout),
	}

	if err := p.run(); err != nil {
		log.Fatal(err)
	}

	p.dst.Flush()
}
