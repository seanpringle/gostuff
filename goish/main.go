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
	return p.scan() == rune(0) || p.scan() == ';' || p.scan() == ')' || p.scan() == '}' || p.peek("end")
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
		ensure(p.peek("do"), "expected: do")
		p.take()
		p.take()
		body := p.block(block, Scope{}, nil).(*NodeBlock)
		ensure(p.peek("end"), "expected: end")
		p.take()
		p.take()
		p.take()
		return NewNodeFunc(args, body)
	}

	if p.peek("for") {
		p.take()
		p.take()
		p.take()
		iter := p.tuple(block, nil)

		body := func() Node {
			ensure(p.peek("do"), "expected: do")
			p.take()
			p.take()
			body := p.block(block, nil, []string{"break", "continue"})
			ensure(p.peek("end"), "expected: end")
			p.take()
			p.take()
			p.take()
			return body
		}

		if p.scan() == ';' {
			p.take()
			begin := iter
			step := p.tuple(block, nil)

			if p.scan() == ';' {
				p.take()
				check := step
				step = p.tuple(block, nil)
				return NewNodeFor3(begin, check, step, body())
			}

			return NewNodeFor2(begin, step, body())
		}

		return NewNodeFor(iter, body())
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
		ensure(p.peek("end"), "expected: end")
		p.take()
		p.take()
		p.take()
		return NewNodeIf(iter, ontrue, onfalse)
	}

	if p.isnumber(p.scan()) {
		for p.isnumber(p.next()) {
			str = append(str, p.take())
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
		ensure(p.scan() == ']', "expected closing bracket (slice)")
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
			if _, is := last.(Operator); is {
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

		if p.scan() == ',' || p.scan() == '{' || p.terminator() {
			break
		}

		if len(ops) > 0 && ops.Last().(Operator).Consumes() > len(items) {
			continue
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
	args := Nodes{}

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
				return Bool{truth(t[0])}
			}
			return Bool{false}
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

		type Bool struct {
			b bool
		}

		func (b Bool) String() string {
			if b.b {
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
			return Bool{len(s.s) > 0}
		}

		type Int struct {
			i64 int64
		}

		func (i Int) String() string {
			return fmt.Sprintf("%d", i.i64)
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
			return Bool{i.i64 != 0}
		}

		func (i Int) Dec() Dec {
			return Dec{float64(i.i64)}
		}

		type Dec struct {
			f64 float64
		}

		func (d Dec) String() string {
			return strconv.FormatFloat(float64(d.f64), 'f', -1, 64)
		}

		func (i Dec) Bool() Bool {
			return Bool{i.f64 < 0 || i.f64 > 0}
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
			return Bool{rune(r) != rune(0)}
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
				n := pos.(IntIsh).Int().i64
				if int64(len(l.data)) > n {
					l.data[n] = val
				}
			}
		}

		func (l *List) Get(pos Any) Any {
			if l != nil {
				n := pos.(IntIsh).Int().i64
				if int64(len(l.data)) > n {
					return l.data[n]
				}
			}
			return nil
		}

		type Func func(Tup) Tup

		func (f Func) Bool() Bool {
			return Bool{true}
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
			if ai, is := a.(IntIsh); is {
				if bi, is := b.(IntIsh); is {
					return Int{ai.Int().i64 + bi.Int().i64}
				}
			}
			if ad, is := a.(DecIsh); is {
				if bd, is := b.(DecIsh); is {
					return Dec{ad.Dec().f64 + bd.Dec().f64}
				}
			}
			if as, is := a.(StrIsh); is {
				return Str{as.Str().s + tostring(b)}
			}
			panic(fmt.Errorf("invalid addition: %v %v", a, b))
		}

		func sub(a, b Any) Any {
			if ai, is := a.(IntIsh); is {
				if bi, is := b.(IntIsh); is {
					return Int{ai.Int().i64 - bi.Int().i64}
				}
			}
			if ad, is := a.(DecIsh); is {
				if bd, is := b.(DecIsh); is {
					return Dec{ad.Dec().f64 - bd.Dec().f64}
				}
			}
			panic(fmt.Errorf("invalid subtraction: %v %v", a, b))
		}

		func mul(a, b Any) Any {
			if ai, is := a.(IntIsh); is {
				if bi, is := b.(IntIsh); is {
					return Int{ai.Int().i64 * bi.Int().i64}
				}
			}
			if ad, is := a.(DecIsh); is {
				if bd, is := b.(DecIsh); is {
					return Dec{ad.Dec().f64 * bd.Dec().f64}
				}
			}
			panic(fmt.Errorf("invalid multiplication: %v %v", a, b))
		}

		func div(a, b Any) Any {
			if ai, is := a.(IntIsh); is {
				if bi, is := b.(IntIsh); is {
					return Int{ai.Int().i64 / bi.Int().i64}
				}
			}
			if ad, is := a.(DecIsh); is {
				if bd, is := b.(DecIsh); is {
					return Dec{ad.Dec().f64 / bd.Dec().f64}
				}
			}
			panic(fmt.Errorf("invalid division: %v %v", a, b))
		}

		func mod(a, b Any) Any {
			if ai, is := a.(IntIsh); is {
				if bi, is := b.(IntIsh); is {
					return Int{ai.Int().i64 % bi.Int().i64}
				}
			}
			panic(fmt.Errorf("invalid modulus: %v %v", a, b))
		}

		func eq(a, b Any) bool {
			if ai, is := a.(IntIsh); is {
				if bi, is := b.(IntIsh); is {
					return ai.Int().i64 == bi.Int().i64
				}
			}
			if ad, is := a.(DecIsh); is {
				if bd, is := b.(DecIsh); is {
					return math.Abs(ad.Dec().f64-bd.Dec().f64) < 0.000001
				}
			}
			panic(fmt.Errorf("invalid comparison (equal): %v %v", a, b))
		}

		func lt(a, b Any) bool {
			if ai, is := a.(IntIsh); is {
				if bi, is := b.(IntIsh); is {
					return ai.Int().i64 < bi.Int().i64
				}
			}
			if ad, is := a.(DecIsh); is {
				if bd, is := b.(DecIsh); is {
					return ad.Dec().f64 < bd.Dec().f64
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
				if b, is := a.(bool); is {
					return b
				}
				if ab, is := a.(BoolIsh); is {
					return ab.Bool().b
				}
			}
			return false
		}

		func join(aa ...Any) Tup {
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

		func get(aa Tup, i int) Any {
			if i < len(aa) {
				return aa[i]
			}
			return nil
		}

		func call(f Any, aa Tup) Any {
			if fn, is := f.(Func); is {
				return fn(aa)
			}
			panic(fmt.Sprintf("attempt to execute a non-function: %v", f))
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

		func trymethod(t Any, k string, def Any) Any {
			t, m := method(t, Str{k})
			if m != nil {
				return call(m, join(t)).(Tup)[0]
			}
			return def
		}

		func tostring(s Any) string {
			return trymethod(s, "string", Str{"nil"}).(Str).s
		}

		func noop(a Any) Any {
			return a
		}

		var Nprint Any = Func(func(aa Tup) Tup {
			for _, a := range aa {
				fmt.Printf("%s", tostring(a))
			}
			return aa
		})

		var Nchan Any = Func(func(t Tup) Tup {
			n := get(t, 0).(IntIsh).Int().i64
			c := make(chan Any, int(n))
			return Tup{Chan(c)}
		})

		var Ngroup Any = Func(func(t Tup) Tup {
			return Tup{NewGroup()}
		})

		var Nlen Any = Func(func(t Tup) Tup {
			s := get(t, 0).(LenIsh)
			return Tup{Int{s.Len()}}
		})

		var Ntype Any = Func(func(t Tup) Tup {
			return Tup{Str{get(t, 0).Type()}}
		})

		var Nstring Any = Func(func(t Tup) Tup {
			return Tup{Str{fmt.Sprintf("%v", get(t, 0))}}
		})

		var Nlib Any = Func(func(t Tup) Tup {
			return Tup{get(t, 0).Lib()}
		})

		func init() {
			libDef = NewMap(MapData{
				Str{"len"}: Nlen,
				Str{"lib"}: Nlib,
				Str{"type"}: Ntype,
				Str{"string"}: Nstring,
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
					return Tup{get(t, 1), get(t, 2)}
				}),
				Str{"get"}: Func(func(t Tup) Tup {
					return Tup{get(t, 0).(*Map).Get(get(t, 1))}
				}),
				Str{"setmeta"}: Func(func(t Tup) Tup {
					get(t, 0).(*Map).meta = get(t, 1)
					return Tup{get(t, 1)}
				}),
				Str{"getmeta"}: Func(func(t Tup) Tup {
					return Tup{get(t, 0).(*Map).meta}
				}),
				Str{"iterate"}: Func(func(t Tup) Tup {
					m := get(t, 0).(*Map)
					keys := []Any{}
					for k, _ := range m.data {
						keys = append(keys, k)
					}
					n := 0
					return Tup{Func(func(tt Tup) Tup {
						k := get(keys, n)
						n = n+1
						if k != nil {
							return Tup{k, m.Get(k)}
						}
						return Tup{nil, nil}
					})}
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
				Str{"iterate"}: Func(func(t Tup) Tup {
					l := get(t, 0).(*List)
					n := 0
					return Tup{Func(func(tt Tup) Tup {
						v := get(l.data, n)
						n = n+1
						return Tup{v}
					})}
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
					return Tup{Bool{true}}
				}),
				Str{"close"}: Func(func(t Tup) Tup {
					c := get(t, 0).(Chan)
					close(c)
					return nil
				}),
				Str{"iterate"}: Func(func(t Tup) Tup {
					c := get(t, 0).(Chan)
					return Tup{Func(func(tt Tup) Tup {
						return Tup{<-c}
					})}
				}),
			})
			libChan.meta = libDef
			libGroup = NewMap(MapData{
				Str{"run"}: Func(func(t Tup) Tup {
					g := get(t, 0).(*Group)
					f := get(t, 1).(Func)
					g.Run(f, t[2:])
					return Tup{Bool{true}}
				}),
				Str{"wait"}: Func(func(t Tup) Tup {
					g := get(t, 0).(*Group)
					g.Wait()
					return Tup{Bool{true}}
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
				Str{"iterate"}: Func(func(t Tup) Tup {
					s := tostring(get(t, 0))
					chars := []rune{}
					for _, c := range s {
						chars = append(chars, c)
					}
					n := 0
					return Tup{Func(func(tt Tup) Tup {
						if len(chars) > n {
							v := Rune(chars[n])
							n = n+1
							return Tup{v}
						}
						return Tup{nil}
					})}
				}),
			})
			libStr.meta = libDef
		}

	`)

	p.println(`func main() {`)
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
