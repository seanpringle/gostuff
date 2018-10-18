package main

import (
	"bufio"
	//"bytes"
	"errors"
	"fmt"
	"io"
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
	return fmt.Sprintf("one(%s)", n.Format())
}

type CanFormatJoin interface {
	FormatJoin() string
}

func FormatJoin(n Node) string {
	if f, is := n.(CanFormatJoin); is {
		return f.FormatJoin()
	}
	return fmt.Sprintf("join(%s)", n.Format())
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
		for c := p.take(); c != rune(0) && c != '\n'; c = p.take() {
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

	if p.scan() == '.' {
		p.take()
		return NewNodeFind()
	}

	if p.scan() == ':' {
		p.take()
		return NewNodeMethod()
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
		ensure(p.scan() == '(', "expected opening paren (func)")
		p.take()
		var args Nodes
		argTup := p.tuple(block, nil)
		if argTup != nil {
			args = Nodes(argTup.(NodeTuple))
		}
		ensure(p.scan() == ')', "expected closing paren (func)")
		p.take()
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

	//var last Node

	for node := p.node(block); node != nil; node = p.node(block) {

		//log.Println(len(items), len(ops), node)
		//log.Println("\t", items)
		//log.Println("\t", ops)

		//if ex, is := node.(*NodeExec); is {
		//	if _, is := last.(Operator); is {
		//		node = ex.args
		//	}
		//}

		//last = node

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

	p.println(`package main`)
	p.println(`import "fmt"`)
	p.println(`import "math"`)
	p.println(`import "strings"`)
	p.println(`import "strconv"`)
	p.println(`import "sync"`)
	p.println(`import "time"`)
	p.println(`import "log"`)
	p.println(`import "os"`)
	p.println(`import "runtime/pprof"`)

	p.println(`type Any interface{
		Type() string
		Lib() *Map
		String() string
	}`)

	p.print(`

		var libDef *Map
		var libMap *Map
		var libList *Map
		var libStr *Map
		var libChan *Map
		var libGroup *Map
		var libTime *Map
		var libTick *Map

		type Stringer = fmt.Stringer

		type Tup []Any

		func (t Tup) Type() string {
			return "tuple"
		}

		func (t Tup) Lib() *Map {
			return libDef
		}

		func (t Tup) String() string {
			items := []string{}
			for _, v := range t {
				if _, is := v.(*Map); is {
					v = Str{"map"}
				}
				if _, is := v.(*List); is {
					v = Str{"slice"}
				}
				if _, is := v.(Tup); is {
					v = Str{"tuple"}
				}
				items = append(items, tostring(v))
			}
			return fmt.Sprintf("[%s]", strings.Join(items, ", "))
		}

		func (t Tup) Bool() Bool {
			if len(t) > 0 {
				return Bool(truth(t[0]))
			}
			return Bool(false)
		}

		type BoolIsh interface {
			Bool() Bool
		}

		type StrIsh interface {
			Str() Str
		}

		type IntIsh interface {
			Int() Int
		}

		type DecIsh interface {
			Dec() Dec
		}

		type LenIsh interface {
			Len() int64
		}

		type Bool bool

		func (b Bool) String() string {
			if bool(b) {
				return "true"
			}
			return "false"
		}

		func (b Bool) Bool() Bool {
			return b
		}

		func (b Bool) Type() string {
			return "bool"
		}

		func (b Bool) Lib() *Map {
			return libDef
		}

		type Str struct {
			s string
		}

		func (s Str) Str() Str {
			return s
		}

		func (s Str) String() string {
			return s.s
		}

		func (s Str) Type() string {
			return "string"
		}

		func (s Str) Len() int64 {
			return int64(len(s.s))
		}

		func (s Str) Lib() *Map {
			return libStr
		}

		func (s Str) Bool() Bool {
			return Bool(len(s.s) > 0)
		}

		type Int int64

		func (i Int) String() string {
			return fmt.Sprintf("%d", int64(i))
		}

		func (i Int) Int() Int {
			return i
		}

		func (i Int) Type() string {
			return "int"
		}

		func (i Int) Lib() *Map {
			return libDef
		}

		func (i Int) Bool() Bool {
			return Bool(int64(i) != 0)
		}

		func (i Int) Dec() Dec {
			return Dec(float64(int64(i)))
		}

		type Dec float64

		func (d Dec) String() string {
			return strconv.FormatFloat(float64(d), 'f', -1, 64)
		}

		func (i Dec) Bool() Bool {
			return Bool(float64(i) < 0 || float64(i) > 0)
		}

		func (d Dec) Dec() Dec {
			return d
		}

		func (d Dec) Type() string {
			return "dec"
		}

		func (d Dec) Lib() *Map {
			return libDef
		}

		type Rune rune

		func (r Rune) Bool() Bool {
			return Bool(rune(r) != rune(0))
		}

		func (r Rune) Type() string {
			return "rune"
		}

		func (r Rune) String() string {
			return string([]rune{rune(r)})
		}

		func (r Rune) Lib() *Map {
			return libDef
		}

		type Time time.Time

		func (t Time) Bool() Bool {
			var z time.Time
			return Bool(time.Time(t) != z)
		}

		func (t Time) Type() string {
			return "time"
		}

		func (t Time) String() string {
			return fmt.Sprintf("%v", time.Time(t))
		}

		func (t Time) Lib() *Map {
			return libTime
		}

		type Ticker struct {
			*time.Ticker
		}

		func NewTicker(d time.Duration) Ticker {
			return Ticker{
				time.NewTicker(d),
			}
		}

		func (t Ticker) Bool() Bool {
			return Bool(t.Ticker != nil)
		}

		func (t Ticker) Type() string {
			return "ticker"
		}

		func (t Ticker) String() string {
			return fmt.Sprintf("%v", t.Ticker)
		}

		func (t Ticker) Lib() *Map {
			return libTick
		}

		func (t Ticker) Read() Any {
			return Time(<-t.Ticker.C)
		}

		type MapData map[Any]Any

		type Map struct {
			data MapData
			meta Any
		}

		func NewMap(data MapData) *Map {
			return &Map{
				data: data,
				meta: libMap,
			}
		}

		func (t *Map) Lib() *Map {
			return t
		}

		func (t *Map) Get(key Any) Any {
			if t != nil {
				if v, ok := t.data[key]; ok {
					return v
				}
			}
			return nil
		}

		func (t *Map) Set(key Any, val Any) {
			if t != nil {
				if val == nil {
					delete(t.data, key)
				} else {
					t.data[key] = val
				}
			}
		}

		func (t *Map) Find(key Any) Any {
			if t != nil {
				if v, ok := t.data[key]; ok {
					return v
				}
				if t.meta != nil {
					if f, is := t.meta.(Func); is {
						return f(Tup{t, key})[0]
					}
					if l, is := t.meta.(Tup); is {
						for _, i := range l {
							r := find(i, key)
							if r != nil {
								return r
							}
						}
					}
					if l, is := t.meta.(*List); is {
						for _, i := range l.data {
							r := find(i, key)
							if r != nil {
								return r
							}
						}
					}
					if t, is := t.meta.(*Map); is {
						return t.Find(key)
					}
					return t.meta
				}
			}
			return nil
		}

		func (t *Map) Type() string {
			return "map"
		}

		func (t *Map) Len() int64 {
			if t != nil {
				return int64(len(t.data))
			}
			return 0
		}

		func (t *Map) String() string {
			if t != nil {
				pairs := []string{}
				for k, v := range t.data {
					if _, is := v.(*Map); is {
						v = Str{"map"}
					}
					if _, is := v.(*List); is {
						v = Str{"slice"}
					}
					pairs = append(pairs, fmt.Sprintf("%v = %v", tostring(k), tostring(v)))
				}
				return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
			}
			return "nil"
		}

		type List struct {
			data []Any
		}

		func NewList(data []Any) *List {
			return &List{data: data}
		}

		func (t *List) Lib() *Map {
			return libList
		}

		func (s *List) Type() string {
			return "list"
		}

		func (s *List) Len() int64 {
			return int64(len(s.data))
		}

		func (s *List) String() string {
			items := []string{}
			for _, v := range s.data {
				if _, is := v.(*Map); is {
					v = Str{"map"}
				}
				if _, is := v.(*List); is {
					v = Str{"slice"}
				}
				if _, is := v.(Tup); is {
					v = Str{"tuple"}
				}
				items = append(items, tostring(v))
			}
			return fmt.Sprintf("[%s]", strings.Join(items, ", "))
		}

		func (l *List) Set(pos Any, val Any) {
			if l != nil {
				n := int64(pos.(IntIsh).Int())
				if int64(len(l.data)) > n {
					l.data[n] = val
				}
			}
		}

		func (l *List) Get(pos Any) Any {
			if l != nil {
				n := int64(pos.(IntIsh).Int())
				if int64(len(l.data)) > n {
					return l.data[n]
				}
			}
			return nil
		}

		type Func func(Tup) Tup

		func (f Func) Bool() Bool {
			return Bool(true)
		}

		func (f Func) String() string {
			return "func"
		}

		func (f Func) Type() string {
			return "func"
		}

		func (f Func) Lib() *Map {
			return libDef
		}

		type Chan chan Any

		func (c Chan) Type() string {
			return "channel"
		}

		func (c Chan) String() string {
			return "channel"
		}

		func (c Chan) Lib() *Map {
			return libChan
		}

		type Group struct {
			g sync.WaitGroup
		}

		func NewGroup() *Group {
			return &Group{}
		}

		func (g *Group) Type() string {
			return "group"
		}

		func (g *Group) String() string {
			return "group"
		}

		func (g *Group) Lib() *Map {
			return libGroup
		}

		func (g *Group) Run(f Func, t Tup) {
			g.g.Add(1)
			go func() {
				defer func() {
					recover()
					g.g.Done()
				}()
				f(append([]Any{g}, t...))
			}()
		}

		func (g *Group) Done() {
			g.g.Done()
		}

		func (g *Group) Wait() {
			g.g.Wait()
			g.g = sync.WaitGroup{}
		}

		func add(a, b Any) Any {
			if ai, is := a.(Int); is {
				if bi, is := b.(Int); is {
					return Int(ai + bi)
				}
			}
			if ai, is := a.(IntIsh); is {
				if bi, is := b.(IntIsh); is {
					return Int(ai.Int() + bi.Int())
				}
			}
			if ad, is := a.(DecIsh); is {
				if bd, is := b.(DecIsh); is {
					return Dec(float64(ad.Dec()) + float64(bd.Dec()))
				}
			}
			if as, is := a.(StrIsh); is {
				return Str{as.Str().s + tostring(b)}
			}
			panic(fmt.Errorf("invalid addition: %v %v", a, b))
		}

		func sub(a, b Any) Any {
			if ai, is := a.(Int); is {
				if bi, is := b.(Int); is {
					return Int(ai - bi)
				}
			}
			if ai, is := a.(IntIsh); is {
				if bi, is := b.(IntIsh); is {
					return Int(ai.Int() - bi.Int())
				}
			}
			if ad, is := a.(DecIsh); is {
				if bd, is := b.(DecIsh); is {
					return Dec(float64(ad.Dec()) - float64(bd.Dec()))
				}
			}
			panic(fmt.Errorf("invalid subtraction: %v %v", a, b))
		}

		func mul(a, b Any) Any {
			if ai, is := a.(Int); is {
				if bi, is := b.(Int); is {
					return Int(ai * bi)
				}
			}
			if ai, is := a.(IntIsh); is {
				if bi, is := b.(IntIsh); is {
					return Int(int64(ai.Int()) * int64(bi.Int()))
				}
			}
			if ad, is := a.(DecIsh); is {
				if bd, is := b.(DecIsh); is {
					return Dec(float64(ad.Dec()) * float64(bd.Dec()))
				}
			}
			panic(fmt.Errorf("invalid multiplication: %v %v", a, b))
		}

		func div(a, b Any) Any {
			if ai, is := a.(Int); is {
				if bi, is := b.(Int); is {
					return Int(ai / bi)
				}
			}
			if ai, is := a.(IntIsh); is {
				if bi, is := b.(IntIsh); is {
					return Int(int64(ai.Int()) / int64(bi.Int()))
				}
			}
			if ad, is := a.(DecIsh); is {
				if bd, is := b.(DecIsh); is {
					return Dec(float64(ad.Dec()) / float64(bd.Dec()))
				}
			}
			panic(fmt.Errorf("invalid division: %v %v", a, b))
		}

		func mod(a, b Any) Any {
			if ai, is := a.(Int); is {
				if bi, is := b.(Int); is {
					return Int(ai % bi)
				}
			}
			if ai, is := a.(IntIsh); is {
				if bi, is := b.(IntIsh); is {
					return Int(int64(ai.Int()) % int64(bi.Int()))
				}
			}
			panic(fmt.Errorf("invalid modulus: %v %v", a, b))
		}

		func eq(a, b Any) bool {
			if ai, is := a.(Int); is {
				if bi, is := b.(Int); is {
					return ai == bi
				}
			}
			if ai, is := a.(IntIsh); is {
				if bi, is := b.(IntIsh); is {
					return int64(ai.Int()) == int64(bi.Int())
				}
			}
			if ad, is := a.(DecIsh); is {
				if bd, is := b.(DecIsh); is {
					return math.Abs(float64(ad.Dec()) - float64(bd.Dec())) < 0.000001
				}
			}
			return func() (rs bool) {
				defer func() {
					recover()
				}()
				rs = a == b
				return
			}()
		}

		func lt(a, b Any) bool {
			if ai, is := a.(Int); is {
				if bi, is := b.(Int); is {
					return ai < bi
				}
			}
			if ai, is := a.(IntIsh); is {
				if bi, is := b.(IntIsh); is {
					return int64(ai.Int()) < int64(bi.Int())
				}
			}
			if ad, is := a.(DecIsh); is {
				if bd, is := b.(DecIsh); is {
					return float64(ad.Dec()) < float64(bd.Dec())
				}
			}
			if as, is := a.(Str); is {
				if bs, is := b.(Str); is {
					return strings.Compare(as.s, bs.s) < 0
				}
			}
			panic(fmt.Errorf("invalid comparison (less): %v %v", a, b))
		}

		func lte(a, b Any) bool {
			return lt(a, b) || eq(a, b)
		}

		func gt(a, b Any) bool {
			return !lte(a, b)
		}

		func gte(a, b Any) bool {
			return !lt(a, b)
		}

		func truth(a interface{}) bool {
			if a != nil {
				if b, is := a.(Bool); is {
					return bool(b)
				}
				if b, is := a.(bool); is {
					return b
				}
				if ab, is := a.(BoolIsh); is {
					return bool(ab.Bool())
				}
			}
			return false
		}

		func join(aa ...Any) Tup {
			if len(aa) == 1 {
				if t, is := aa[0].(Tup); is {
					return t
				} else {
					return aa
				}
			}
			var rr Tup
			for _, a := range aa {
				if t, is := a.(Tup); is {
					rr = append(rr, t...)
				} else {
					rr = append(rr, a)
				}
			}
			return rr
		}

		func one(a Any) Any {
			if t, is := a.(Tup); is {
				if len(t) > 0 {
					return t[0]
				}
				return nil
			}
			return a
		}

		func get(aa Tup, i int) Any {
			if aa != nil && i < len(aa) {
				return aa[i]
			}
			return nil
		}

		func call(f Any, aa Tup) Tup {
			return f.(Func)(aa)
		}

		func find(t Any, key Any) Any {
			if t != nil {
				return t.Lib().Find(key)
			}
			return nil
		}

		func method(t Any, key Any) (Any, Any) {
			return t, find(t, key)
		}

		func store(t Any, key Any, val Any) Any {
			t.(*Map).Set(key, val)
			return val
		}

		func iterate(o Any) Func {
			if oi := trymethod(o, "iterate", nil); oi != nil {
				return oi.(Func)
			}
			panic(fmt.Sprintf("not iterable: %v", o))
		}

		type loopBroke int

		func loopbreak() {
			panic(loopBroke(0))
		}

		func loop(fn func()) {
			defer func() {
				if r := recover(); r != nil {
					if _, is := r.(loopBroke); is {
						return
					}
					panic(r)
				}
			}()

			fn()
		}

		func trymethod(t Any, k string, def Any) Any {
			t, m := method(t, Str{k})
			if m != nil {
				return get(call(m, join(t)), 0)
			}
			return def
		}

		func tostring(s Any) string {
			return trymethod(s, "string", Str{"nil"}).(Str).s
		}

		func noop(a Any) {
		}

		var Nprint Any = Func(func(aa Tup) Tup {
			parts := []string{}
			for _, a := range aa {
				parts = append(parts, tostring(a))
			}
			fmt.Printf("%s\n", strings.Join(parts, " "))
			return aa
		})

		var Nchan Any = Func(func(t Tup) Tup {
			n := int64(get(t, 0).(IntIsh).Int())
			c := make(chan Any, int(n))
			return Tup{Chan(c)}
		})

		var Ngroup Any = Func(func(t Tup) Tup {
			return Tup{NewGroup()}
		})

		var Ntime *Map

		var Ntype Any = Func(func(t Tup) Tup {
			return Tup{Str{get(t, 0).Type()}}
		})

		var Nsetprototype Any = Func(func(t Tup) Tup {
			if m, is := get(t, 0).(*Map); is {
				m.meta = get(t, 1)
				return nil
			}
			panic(fmt.Sprintf("cannot set prototype"))
		})

		var Ngetprototype Any = Func(func(t Tup) Tup {
			if m, is := get(t, 0).(*Map); is {
				return Tup{m.meta}
			}
			return Tup{get(t, 0).Lib()}
		})

		func init() {
			libDef = NewMap(MapData{
				Str{"len"}: Func(func(t Tup) Tup {
					s := get(t, 0).(LenIsh)
					return Tup{Int(s.Len())}
				}),
				Str{"type"}: Ntype,
				Str{"string"}: Func(func(t Tup) Tup {
					return Tup{Str{fmt.Sprintf("%v", get(t, 0))}}
				}),
			})
			libMap = NewMap(MapData{
				Str{"keys"}: Func(func(t Tup) Tup {
					keys := []Any{}
					for k, _ := range get(t, 0).(*Map).data {
						keys = append(keys, k)
					}
					return Tup{NewList(keys)}
				}),
				Str{"set"}: Func(func(t Tup) Tup {
					get(t, 0).(*Map).Set(get(t, 1), get(t, 2))
					return Tup{get(t, 0)}
				}),
				Str{"get"}: Func(func(t Tup) Tup {
					return Tup{get(t, 0).(*Map).Get(get(t, 1))}
				}),
			})
			libMap.meta = libDef
			libList = NewMap(MapData{
				Str{"push"}: Func(func(t Tup) Tup {
					l := get(t, 0).(*List)
					v := get(t, 1)
					l.data = append(l.data, v)
					return Tup{l}
				}),
				Str{"pop"}: Func(func(t Tup) Tup {
					l := get(t, 0).(*List)
					n := len(l.data)-1
					var v Any
					if len(l.data) < n {
						v = l.data[n]
						l.data = l.data[0:n]
					}
					return Tup{v}
				}),
				Str{"join"}: Func(func(t Tup) Tup {
					l := get(t, 0).(*List)
					j := get(t, 1)
					var ls []string
					for _, s := range l.data {
						ls = append(ls, tostring(s))
					}
					return Tup{Str{strings.Join(ls, tostring(j))}}
				}),
				Str{"set"}: Func(func(t Tup) Tup {
					get(t, 0).(*List).Set(get(t, 1), get(t, 2))
					return Tup{get(t, 1), get(t, 2)}
				}),
				Str{"get"}: Func(func(t Tup) Tup {
					return Tup{get(t, 0).(*List).Get(get(t, 1))}
				}),
			})
			libList.meta = libDef
			libChan = NewMap(MapData{
				Str{"read"}: Func(func(t Tup) Tup {
					c := get(t, 0).(Chan)
					return Tup{<-c}
				}),
				Str{"write"}: Func(func(t Tup) Tup {
					c := get(t, 0).(Chan)
					a := get(t, 1)
					c <-a
					return Tup{Bool(true)}
				}),
				Str{"close"}: Func(func(t Tup) Tup {
					c := get(t, 0).(Chan)
					close(c)
					return nil
				}),
			})
			libChan.meta = libDef
			libGroup = NewMap(MapData{
				Str{"run"}: Func(func(t Tup) Tup {
					g := get(t, 0).(*Group)
					f := get(t, 1).(Func)
					g.Run(f, t[2:])
					return Tup{Bool(true)}
				}),
				Str{"wait"}: Func(func(t Tup) Tup {
					g := get(t, 0).(*Group)
					g.Wait()
					return Tup{Bool(true)}
				}),
			})
			libGroup.meta = libDef
			libStr = NewMap(MapData{
				Str{"split"}: Func(func(t Tup) Tup {
					s := tostring(get(t, 0))
					j := tostring(get(t, 1))
					l := []Any{}
					for _, p := range strings.Split(s, j) {
						l = append(l, Str{p})
					}
					return Tup{NewList(l)}
				}),
			})
			libStr.meta = libDef
			libTick = NewMap(MapData{
				Str{"read"}: Func(func(t Tup) Tup {
					ti := get(t, 0).(Ticker)
					return Tup{ti.Read()}
				}),
			})
			libTick.meta = libDef
			libTime = NewMap(MapData{
				Str{"ms"}: Int(int64(time.Millisecond)),
				Str{"ticker"}: Func(func(t Tup) Tup {
					d := int64(get(t, 0).(Int))
					return Tup{NewTicker(time.Duration(d))}
				}),
			})
			libTime.meta = libDef
			Ntime = libTime
		}

	`)

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

	`)

	block := p.block(nil, Scope{}, nil)
	p.println(block.Format())
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
