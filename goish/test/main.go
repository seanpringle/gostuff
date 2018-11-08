package main

import "fmt"
import "math"
import "strings"

//import "strconv"
import "sync"
import "time"
import "log"
import "os"
import "sort"
import "io"
import "io/ioutil"
import "runtime/pprof"
import "encoding/hex"
import "bufio"
import "errors"
import "encoding/json"
import "regexp"

type Any interface {
	Type() string
	Lib() Searchable
	String() string
}

type Searchable interface {
	Search(Any) Any
}

type Linkable interface {
	Link(Searchable)
}

type VM struct {
	argsC   [8]*Args
	argsN   int
	reCache map[string]*regexp.Regexp
}

func (vm *VM) ga(n int) *Args {
	if n > 32 {
		panic("maximum of 32 arguments per call")
	}
	if vm.argsN > 0 {
		vm.argsN--
		aa := vm.argsC[vm.argsN]
		aa.used = 0
		return aa
	}
	return &Args{}
}

func (vm *VM) da(a *Args) {
	if a != nil && vm.argsN < 8 {
		a.used = 0
		vm.argsC[vm.argsN] = a
		vm.argsN++
	}
}

type Args struct {
	used  uint
	cells [32]Any
}

func (aa *Args) Type() string {
	return "args"
}

func (aa *Args) String() string {
	return "args"
}

func (aa *Args) Lib() Searchable {
	return protoDef
}

func (aa *Args) get(i int) Any {
	if i >= 32 {
		panic("maximum of 32 arguments per call")
	}
	if aa != nil && (1<<uint(i))&aa.used != 0 {
		return aa.cells[i]
	}
	return nil
}

func (aa *Args) agg(i int) Any {
	if i >= 32 {
		panic("maximum of 32 arguments per call")
	}
	l := []Any{}
	for aa != nil && i < 32 && (1<<uint(i))&aa.used != 0 {
		l = append(l, aa.cells[i])
		i++
	}
	return NewList(l)
}

func (aa *Args) set(i int, v Any) {
	if i >= 32 {
		panic("maximum of 32 arguments per call")
	}
	if aa != nil {
		aa.cells[i] = v
		aa.used = aa.used | (1 << uint(i))
	}
}

func (aa *Args) len() int {
	l := 0
	if aa != nil {
		u := aa.used
		for u != 0 {
			l++
			u = u >> 1
		}
	}
	return l
}

var protoDef *Map
var protoInt *Map
var protoDec *Map
var protoMap *Map
var protoList *Map
var protoText *Map
var protoChan *Map
var protoGroup *Map
var protoInst *Map
var protoTick *Map
var protoBlob *Map
var protoStream *Map

var libIO *Map
var libTime *Map
var libSync *Map

var onInit []func()

type Stringer = fmt.Stringer

type BoolIsh interface {
	Bool() Bool
}

type TextIsh interface {
	Text() Text
}

type BlobIsh interface {
	Blob() Blob
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

func (b Bool) Lib() Searchable {
	return protoDef
}

type Text string

func (s Text) Text() Text {
	return s
}

func (s Text) String() string {
	return string(s)
}

func (s Text) Type() string {
	return "text"
}

func (s Text) Len() int64 { // characters
	l := int64(0)
	for range string(s) {
		l++
	}
	return l
}

func (s Text) Lib() Searchable {
	return protoText
}

func (s Text) Bool() Bool {
	return Bool(len(s) > 0)
}

func (s Text) Blob() Blob {
	return []byte(string(s))
}

type Blob []byte

func (b Blob) Type() string {
	return "blob"
}

func (b Blob) Lib() Searchable {
	return protoBlob
}

func (b Blob) String() string {
	return hex.Dump([]byte(b))
}

func (b Blob) Blob() Blob {
	return b
}

func (b Blob) Text() Text {
	return Text(string(b))
}

func (b Blob) Len() int64 {
	return int64(len(b))
}

func (b Blob) Bool() Bool {
	return Bool(len(b) > 0)
}

type Int int64

func (i Int) String() string {
	return fmt.Sprintf("%v", int64(i))
}

func (i Int) Int() Int {
	return i
}

func (i Int) Type() string {
	return "int"
}

func (i Int) Lib() Searchable {
	return protoInt
}

func (i Int) Bool() Bool {
	return Bool(int64(i) != 0)
}

func (i Int) Dec() Dec {
	return Dec(float64(int64(i)))
}

func (i Int) Text() Text {
	return Text(i.String())
}

type SInt struct {
	i Int
}

func (i *SInt) String() string {
	return i.i.String()
}

func (i *SInt) Int() Int {
	return i.i
}

func (i *SInt) Type() string {
	return "int"
}

func (i *SInt) Lib() Searchable {
	return protoInt
}

func (i *SInt) Bool() Bool {
	return i.i.Bool()
}

func (i *SInt) Dec() Dec {
	return Dec(float64(int64(i.i)))
}

func (i *SInt) Text() Text {
	return Text(i.String())
}

const IntCache = 4096

var SInts [IntCache]SInt

func init() {
	for i := 0; i < IntCache; i++ {
		SInts[i] = SInt{Int(i)}
	}
}

type Dec float64

func (d Dec) String() string {
	//return strconv.FormatFloat(float64(d), 'f', -1, 64)
	return fmt.Sprintf("%v", float64(d))
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

func (d Dec) Lib() Searchable {
	return protoDec
}

func (d Dec) Text() Text {
	return Text(d.String())
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

func (r Rune) Text() Text {
	return Text(r.String())
}

func (r Rune) Lib() Searchable {
	return protoDef
}

type Status struct {
	e error
}

func NewStatus(e error) Status {
	return Status{e: e}
}

func (e Status) Bool() Bool {
	return Bool(e.e == nil)
}

func (e Status) Type() string {
	return "status"
}

func (e Status) String() string {
	if e.e != nil {
		return error(e.e).Error()
	}
	return "ok"
}

func (e Status) Lib() Searchable {
	return protoDef
}

type Instant time.Time

func (t Instant) Bool() Bool {
	var z time.Time
	return Bool(time.Time(t) != z)
}

func (t Instant) Type() string {
	return "instant"
}

func (t Instant) String() string {
	return fmt.Sprintf("%v", time.Time(t))
}

func (t Instant) Lib() Searchable {
	return protoInst
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
	return "ticker"
}

func (t Ticker) Lib() Searchable {
	return protoTick
}

func (t Ticker) Read() Any {
	return Instant(<-t.Ticker.C)
}

func (t Ticker) Stop() Any {
	t.Ticker.Stop()
	return nil
}

type MapData map[Any]Any

type Map struct {
	data MapData
	meta Searchable
}

func NewMap(data MapData) *Map {
	return &Map{
		data: data,
		meta: protoMap,
	}
}

func (t *Map) Lib() Searchable {
	return t
}

func (t *Map) Link(lib Searchable) {
	t.meta = lib
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
		t.data[key] = val
	}
}

func (t *Map) Has(key Any) Any {
	if t != nil {
		_, has := t.data[key]
		return Bool(has)
	}
	return Bool(false)
}

func (t *Map) Drop(key Any) {
	if t != nil {
		delete(t.data, key)
	}
}

func (t *Map) Search(key Any) Any {
	if t != nil {
		if v, ok := t.data[key]; ok {
			return v
		}
		if t.meta != nil {
			return t.meta.Search(key)
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

func (t *Map) Bool() Bool {
	return Bool(t.Len() > 0)
}

func (t *Map) String() string {
	if t != nil {
		pairs := []string{}
		for k, v := range t.data {
			if _, is := v.(*Map); is {
				v = Text(v.Type())
			}
			if _, is := v.(*List); is {
				v = Text(v.Type())
			}
			//pairs = append(pairs, fmt.Sprintf("%s = %s", tostring(k), tostring(v)))
			pairs = append(pairs, fmt.Sprintf("%s = %s", tostring(k), tostring(v)))
		}
		return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
	}
	return "nil"
}

type List struct {
	data []Any
	meta Searchable
}

func NewList(data []Any) *List {
	return &List{data: data, meta: protoList}
}

func (t *List) Lib() Searchable {
	return t.meta
}

func (l *List) Link(lib Searchable) {
	l.meta = lib
}

func (l *List) Type() string {
	return "list"
}

func (l *List) Len() int64 {
	return int64(len(l.data))
}

func (l *List) Bool() Bool {
	return Bool(len(l.data) > 0)
}

func (s *List) String() string {
	items := []string{}
	for _, v := range s.data {
		if _, is := v.(*Map); is {
			v = Text(v.Type())
		}
		if _, is := v.(*List); is {
			v = Text(v.Type())
		}
		items = append(items, tostring(v))
	}
	return fmt.Sprintf("[%s]", strings.Join(items, ", "))
}

func (l *List) Set(pos Any, val Any) {
	if l != nil {
		n := int64(pos.(IntIsh).Int())
		if n >= 0 && int64(len(l.data)) > n {
			l.data[n] = val
		}
	}
}

func (l *List) Get(pos Any) Any {
	if l != nil {
		n := int64(pos.(IntIsh).Int())
		if n >= 0 && int64(len(l.data)) > n {
			return l.data[n]
		}
	}
	return nil
}

func (l *List) Search(key Any) Any {
	for _, item := range l.data {
		if lib, is := item.(Searchable); is {
			if val := lib.Search(key); val != nil {
				return val
			}
		}
	}
	return nil
}

type Func func(*VM, *Args) *Args

func (f Func) Bool() Bool {
	return Bool(true)
}

func (f Func) String() string {
	return "func"
}

func (f Func) Type() string {
	return "func"
}

func (f Func) Lib() Searchable {
	return protoDef
}

func (f Func) Search(key Any) Any {
	vm := &VM{}
	return one(vm, call(vm, f, join(vm, key)))
}

type Chan struct {
	c chan Any
}

func NewChan(n int) *Chan {
	return &Chan{
		c: make(chan Any, n),
	}
}

func (c *Chan) Type() string {
	return "channel"
}

func (c *Chan) String() string {
	return "channel"
}

func (c *Chan) Lib() Searchable {
	return protoChan
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

func (g *Group) Lib() Searchable {
	return protoGroup
}

func (g *Group) Run(f Func, t *Args) {
	g.g.Add(1)
	vm := &VM{}
	aa := vm.ga(t.len() + 1)
	aa.set(0, g)
	for i := 0; i < t.len(); i++ {
		aa.set(i+1, t.get(i))
	}
	go func() {
		defer func() {
			recover()
			g.g.Done()
		}()
		f(vm, aa)
	}()
}

func (g *Group) Done() {
	g.g.Done()
}

func (g *Group) Wait() {
	g.g.Wait()
	g.g = sync.WaitGroup{}
}

type Stream struct {
	s interface{}
}

func NewStream(s interface{}) Stream {
	r, reader := s.(io.Reader)
	w, writer := s.(io.Writer)
	if reader && writer {
		return Stream{bufio.NewReadWriter(bufio.NewReader(r), bufio.NewWriter(w))}
	}
	if writer {
		return Stream{bufio.NewWriter(w)}
	}
	return Stream{bufio.NewReader(r)}
}

func (s Stream) Type() string {
	return "stream"
}

func (s Stream) Lib() Searchable {
	return protoStream
}

func (s Stream) String() string {
	return "stream"
}

func (s Stream) Flush() error {
	if f, is := s.s.(Flusher); is {
		return f.Flush()
	}
	return errors.New("not a flusher")
}

func (s Stream) Close() error {
	if c, is := s.s.(Closer); is {
		return c.Close()
	}
	return errors.New("not a closer")
}

func isInt(a Any) (n Int, b bool) {
	if ai, is := a.(Int); is {
		return ai, true
	}
	if ai, is := a.(*SInt); is {
		return ai.Int(), true
	}
	return 0, false
}

func toInt(n Int) Any {
	if n >= 0 && n < IntCache {
		return &SInts[n]
	}
	return n
}

func concat(a, b Any) Text {
	return Text(totext(a) + totext(b))
}

func add(a, b Any) Any {
	if ai, is := isInt(a); is {
		if bi, is := isInt(b); is {
			return toInt(ai + bi)
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
	panic(fmt.Errorf("invalid addition: %v %v", a, b))
}

func sub(a, b Any) Any {
	if ai, is := isInt(a); is {
		if bi, is := isInt(b); is {
			return toInt(ai - bi)
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
	if ai, is := isInt(a); is {
		if bi, is := isInt(b); is {
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
	if ai, is := isInt(a); is {
		if bi, is := isInt(b); is {
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
	if ai, is := isInt(a); is {
		if bi, is := isInt(b); is {
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
	if ai, is := isInt(a); is {
		if bi, is := isInt(b); is {
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
			return math.Abs(float64(ad.Dec())-float64(bd.Dec())) < 0.000001
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
	if ai, is := isInt(a); is {
		if bi, is := isInt(b); is {
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
	if as, is := a.(Text); is {
		if bs, is := b.(Text); is {
			return strings.Compare(string(as), string(bs)) < 0
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
		if aa, is := a.(Any); is {
			return tobool(aa)
		}
	}
	return false
}

func join(vm *VM, t ...Any) *Args {
	if len(t) == 1 {
		if aa, is := t[0].(*Args); is {
			return aa
		}
		r := vm.ga(1)
		r.set(0, t[0])
		return r
	}
	l := 0
	for _, v := range t {
		if aa, is := v.(*Args); is {
			l += aa.len()
		} else {
			l++
		}
	}
	r := vm.ga(l)
	i := 0
	for _, v := range t {
		if aa, is := v.(*Args); is {
			l := aa.len()
			for j := 0; j < l; j++ {
				r.set(i, aa.get(j))
				i++
			}
			vm.da(aa)
		} else {
			r.set(i, v)
			i++
		}
	}
	return r
}

func one(vm *VM, a Any) Any {
	if aa, is := a.(*Args); is {
		v := aa.get(0)
		vm.da(aa)
		return v
	}
	return a
}

func call(vm *VM, f Any, aa *Args) *Args {
	return f.(Func)(vm, aa)
}

func find(t Any, key Any) Any {
	if t != nil {
		return t.Lib().Search(key)
	} else {
		return protoDef.Search(key)
	}
	return nil
}

func field(t Any, key Any) Any {
	if t != nil {
		if m, is := t.(*Map); is {
			return m.Get(key)
		}
		if l, is := t.(*List); is {
			return l.Get(key)
		}
		panic(fmt.Sprintf("invalid retrieve operation: %v", t))
	}
	return nil
}

func method(t Any, key Any) (Any, Any) {
	return t, find(t, key)
}

func store(t Any, key Any, val Any) Any {
	if t != nil {
		if m, is := t.(*Map); is {
			m.Set(key, val)
			return val
		}
		if l, is := t.(*List); is {
			l.Set(key, val)
			return val
		}
		panic(fmt.Sprintf("invalid store operation: %v", t))
	}
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

func extract(vm *VM, src Any) *Args {
	if src != nil {
		l := src.(*List).data
		aa := vm.ga(len(l))
		for i, v := range l {
			aa.set(i, v)
		}
		return aa
	}
	return nil
}

func trymethod(t Any, k string, def Any) Any {
	t, m := method(t, Text(k))
	if m != nil {
		if _, is := m.(Func); is {
			vm := &VM{}
			return call(vm, m, join(vm, t)).get(0)
		}
		return m
	}
	return def
}

func tostring(s Any) string {
	return trymethod(s, "string", Text("(nil)")).(TextIsh).Text().String()
}

func totext(s Any) string {
	t := trymethod(s, "text", s)
	return t.(TextIsh).Text().String()
}

func tobool(s Any) bool {
	return bool(trymethod(s, "bool", Bool(false)).(BoolIsh).Bool())
}

func toregexp(vm *VM, p string) *regexp.Regexp {

	if vm.reCache == nil {
		vm.reCache = map[string]*regexp.Regexp{}
	}

	re, have := vm.reCache[p]
	if !have {
		re = regexp.MustCompile(p)
		vm.reCache[p] = re
	}

	return re
}

func length(v Any) Any {
	if l, is := v.(LenIsh); is {
		return Int(l.Len())
	}
	return trymethod(v, "len", nil)
}

func noop(a Any) {
}

func ifnil(a, b Any) Any {
	if a == nil {
		return b
	}
	return a
}

var Nprint Any = Func(func(vm *VM, aa *Args) *Args {
	parts := []string{}
	for i := 0; i < aa.len(); i++ {
		parts = append(parts, tostring(aa.get(i)))
	}
	fmt.Printf("%s\n", strings.Join(parts, "\t"))
	return aa
})

var Nlog Any = Func(func(vm *VM, aa *Args) *Args {
	parts := []string{}
	for i := 0; i < aa.len(); i++ {
		parts = append(parts, tostring(aa.get(i)))
	}
	fmt.Fprintf(os.Stderr, "%s\n", strings.Join(parts, "\t"))
	return aa
})

var Nio *Map
var Ntime *Map
var Nsync *Map

var Ntype Any = Func(func(vm *VM, aa *Args) *Args {
	v := aa.get(0)
	vm.da(aa)
	if v != nil {
		return join(vm, Text(v.Type()))
	}
	return join(vm, Text("nil"))
})

var Nerror Any = Func(func(vm *VM, aa *Args) *Args {
	msg := aa.get(0)
	vm.da(aa)
	if msg != nil {
		return join(vm, NewStatus(errors.New(tostring(msg))))
	}
	return join(vm, NewStatus(nil))
})

var Nsetprototype Any = Func(func(vm *VM, aa *Args) *Args {
	if l, is := aa.get(0).(Searchable); is {
		l.(Linkable).Link(aa.get(1).(Searchable))
		vm.da(aa)
		aa = vm.ga(1)
		aa.set(0, l.(Any))
		return aa
	}
	panic(fmt.Sprintf("cannot set prototype"))
})

var Ngetprototype Any = Func(func(vm *VM, aa *Args) *Args {
	v := aa.get(0)
	vm.da(aa)
	if v == nil {
		return join(vm, protoDef)
	}
	if m, is := v.(*Map); is {
		return join(vm, m.meta.(Any))
	}
	return join(vm, v.Lib().(Any))
})

type ReadRuner interface {
	ReadRune() (rune, int, error)
}

type Flusher interface {
	Flush() error
}

type Closer interface {
	Close() error
}

func init() {
	protoDef = NewMap(MapData{
		Text("string"): Func(func(vm *VM, aa *Args) *Args {
			a := aa.get(0)
			vm.da(aa)
			if a == nil {
				return join(vm, Text("(nil)"))
			}
			return join(vm, Text(a.String()))
		}),
		Text("text"): Func(func(vm *VM, aa *Args) *Args {
			a := aa.get(0)
			vm.da(aa)
			return join(vm, a.(TextIsh).Text())
		}),
	})

	protoInt = NewMap(MapData{
		Text("huge"): Int(math.MaxInt64),
		Text("tiny"): Int(math.MinInt64),
	})
	protoInt.meta = protoDef

	protoDec = NewMap(MapData{
		Text("huge"): Dec(math.MaxFloat64),
		Text("tiny"): Dec(math.SmallestNonzeroFloat64),
	})
	protoDec.meta = protoDef

	protoMap = NewMap(MapData{
		Text("keys"): Func(func(vm *VM, aa *Args) *Args {
			keys := []Any{}
			for k, _ := range aa.get(0).(*Map).data {
				keys = append(keys, k)
			}
			vm.da(aa)
			return join(vm, NewList(keys))
		}),
	})
	protoMap.meta = protoDef

	protoList = NewMap(MapData{
		Text("insert"): Func(func(vm *VM, aa *Args) *Args {
			l := aa.get(0).(*List)
			p := int(aa.get(1).(IntIsh).Int())
			v := aa.get(2)
			vm.da(aa)
			if p >= len(l.data) {
				l.data = append(l.data, v)
				return join(vm, l)
			}
			if p <= 0 {
				l.data = append([]Any{v}, l.data...)
				return join(vm, l)
			}
			l.data = append(l.data[0:p], append([]Any{v}, l.data[p:]...)...)
			return join(vm, l)
		}),
		Text("remove"): Func(func(vm *VM, aa *Args) *Args {
			l := aa.get(0).(*List)
			p := int(aa.get(1).(IntIsh).Int())
			vm.da(aa)
			if p >= len(l.data) || p < 0 {
				return join(vm, nil)
			}
			v := l.data[p]
			l.data = append(l.data[0:p], l.data[p+1:]...)
			return join(vm, v)
		}),
		Text("join"): Func(func(vm *VM, aa *Args) *Args {
			l := aa.get(0).(*List)
			j := ifnil(aa.get(1), Text(""))
			vm.da(aa)
			var ls []string
			for _, s := range l.data {
				ls = append(ls, totext(s))
			}
			return join(vm, Text(strings.Join(ls, totext(j))))
		}),
		Text("sort"): Func(func(vm *VM, aa *Args) *Args {
			l := aa.get(0).(*List)
			f := aa.get(1).(Func)
			vm.da(aa)
			sort.SliceStable(l.data, func(a, b int) bool {
				return truth(one(vm, f(vm, join(vm, l.data[a], l.data[b]))))
			})
			return join(vm, l)
		}),
	})
	protoList.meta = protoDef

	protoChan = NewMap(MapData{
		Text("read"): Func(func(vm *VM, aa *Args) *Args {
			c := aa.get(0).(*Chan)
			vm.da(aa)
			return join(vm, <-c.c)
		}),
		Text("write"): Func(func(vm *VM, aa *Args) *Args {
			c := aa.get(0).(*Chan)
			a := aa.get(1)
			vm.da(aa)
			rs := func() (err error) {
				defer func() {
					if r := recover(); r != nil {
						err = errors.New("channel closed")
					}
				}()
				c.c <- a
				return
			}()
			return join(vm, NewStatus(rs))
		}),
		Text("close"): Func(func(vm *VM, aa *Args) *Args {
			c := aa.get(0).(*Chan)
			vm.da(aa)
			close(c.c)
			return nil
		}),
	})
	protoChan.meta = protoDef
	protoGroup = NewMap(MapData{
		Text("run"): Func(func(vm *VM, aa *Args) *Args {
			g := aa.get(0).(*Group)
			f := aa.get(1).(Func)
			ab := vm.ga(aa.len() - 2)
			for i := 0; i < aa.len()-2; i++ {
				ab.set(i, aa.get(i+2))
			}
			g.Run(f, ab)
			vm.da(aa)
			return join(vm, NewStatus(nil))
		}),
		Text("wait"): Func(func(vm *VM, aa *Args) *Args {
			g := aa.get(0).(*Group)
			vm.da(aa)
			g.Wait()
			return join(vm, NewStatus(nil))
		}),
	})
	protoGroup.meta = protoDef
	protoText = NewMap(MapData{
		Text("blob"): Func(func(vm *VM, aa *Args) *Args {
			s := totext(aa.get(0))
			vm.da(aa)
			return join(vm, Blob(s))
		}),
		Text("json"): Func(func(vm *VM, aa *Args) *Args {
			s := totext(aa.get(0))
			vm.da(aa)
			var m interface{}
			err := json.Unmarshal([]byte(s), &m)
			var walk func(v interface{}) Any
			walk = func(v interface{}) Any {
				switch v.(type) {
				case int:
					return Int(v.(int))
				case int64:
					return Int(v.(int64))
				case float64:
					return Dec(v.(float64))
				case string:
					return Text(v.(string))
				case []interface{}:
					vals := []Any{}
					for _, v := range v.([]interface{}) {
						vals = append(vals, walk(v))
					}
					return NewList(vals)
				case map[string]interface{}:
					pairs := NewMap(MapData{})
					for k, v := range v.(map[string]interface{}) {
						pairs.Set(Text(k), walk(v))
					}
					return pairs
				}
				return nil
			}
			return join(vm, NewStatus(err), walk(m))
		}),

		Text("split"): Func(func(vm *VM, aa *Args) *Args {
			s := totext(aa.get(0))
			j := totext(aa.get(1))
			vm.da(aa)
			l := []Any{}
			for _, p := range strings.Split(s, j) {
				l = append(l, Text(p))
			}
			return join(vm, NewList(l))
		}),

		Text("match"): Func(func(vm *VM, aa *Args) *Args {
			s := totext(aa.get(0))
			p := totext(aa.get(1))
			vm.da(aa)
			re := toregexp(vm, p)
			m := re.FindStringSubmatch(s)
			if m != nil {
				l := []Any{}
				for _, ss := range m {
					l = append(l, Text(ss))
				}
				return join(vm, NewList(l))
			}
			return join(vm, nil)
		}),
	})
	protoText.meta = protoDef

	protoTick = NewMap(MapData{
		Text("read"): Func(func(vm *VM, aa *Args) *Args {
			ti := aa.get(0).(Ticker)
			vm.da(aa)
			return join(vm, ti.Read())
		}),
		Text("stop"): Func(func(vm *VM, aa *Args) *Args {
			ti := aa.get(0).(Ticker)
			vm.da(aa)
			return join(vm, ti.Stop())
		}),
	})
	protoTick.meta = protoDef

	protoInst = protoDef

	libTime = NewMap(MapData{
		Text("ms"): Int(int64(time.Millisecond)),
		Text("ticker"): Func(func(vm *VM, aa *Args) *Args {
			d := int64(aa.get(0).(Int))
			vm.da(aa)
			return join(vm, NewTicker(time.Duration(d)))
		}),
	})
	Ntime = libTime

	libSync = NewMap(MapData{
		Text("group"): Func(func(vm *VM, aa *Args) *Args {
			vm.da(aa)
			return join(vm, NewGroup())
		}),
		Text("channel"): Func(func(vm *VM, aa *Args) *Args {
			n := int(aa.get(0).(IntIsh).Int())
			vm.da(aa)
			return join(vm, NewChan(n))
		}),
	})
	Nsync = libSync

	libIO = NewMap(MapData{
		Text("stdin"):  NewStream(os.Stdin),
		Text("stdout"): NewStream(os.Stdout),
		Text("stderr"): NewStream(os.Stderr),

		Text("open"): Func(func(vm *VM, aa *Args) *Args {
			path := totext(aa.get(0))
			modes := totext(aa.get(1))
			mode := os.O_RDONLY
			switch modes {
			case "r":
				mode = os.O_RDONLY
			case "w":
				mode = os.O_WRONLY
			case "a":
				mode = os.O_APPEND
			case "r+":
				mode = os.O_RDWR
			case "w+":
				mode = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
			case "a+":
				mode = os.O_RDONLY | os.O_CREATE | os.O_APPEND
			default:
				panic(fmt.Sprintf("unknown file acces mode: %s", modes))
			}
			file, err := os.OpenFile(path, mode, 0644)
			if err != nil {
				return join(vm, NewStatus(err))
			}
			return join(vm, NewStatus(nil), NewStream(file))
		}),
	})
	Nio = libIO

	protoBlob = NewMap(MapData{
		Text("text"): Func(func(vm *VM, aa *Args) *Args {
			bin := aa.get(0).(BlobIsh).Blob()
			vm.da(aa)
			return join(vm, bin.Text())
		}),
	})
	protoBlob.meta = protoDef

	protoStream = NewMap(MapData{

		Text("read"): Func(func(vm *VM, aa *Args) *Args {
			stream := aa.get(0).(Stream)
			limit := ifnil(aa.get(1), Int(1024*1024)).(IntIsh).Int()
			vm.da(aa)
			if _, is := stream.s.(io.Reader); !is {
				return join(vm, NewStatus(errors.New("not a reader")))
			}
			buff := make([]byte, int(limit))
			length, err := stream.s.(io.Reader).Read(buff)
			return join(vm, NewStatus(err), Blob(buff[:length]))
		}),

		Text("readall"): Func(func(vm *VM, aa *Args) *Args {
			stream := aa.get(0).(Stream)
			vm.da(aa)
			if _, is := stream.s.(io.Reader); !is {
				return join(vm, NewStatus(errors.New("not a reader")))
			}
			buff, err := ioutil.ReadAll(stream.s.(io.Reader))
			return join(vm, NewStatus(err), Blob(buff))
		}),

		Text("readrune"): Func(func(vm *VM, aa *Args) *Args {
			stream := aa.get(0).(Stream)
			vm.da(aa)
			if r, is := stream.s.(ReadRuner); is {
				char, _, err := r.ReadRune()
				return join(vm, NewStatus(err), Rune(char))
			}
			return join(vm, NewStatus(errors.New("not a reader")))
		}),

		Text("write"): Func(func(vm *VM, aa *Args) *Args {
			stream := aa.get(0).(Stream)
			data := aa.get(1).(BlobIsh).Blob()
			vm.da(aa)
			if _, is := stream.s.(io.Writer); !is {
				return join(vm, NewStatus(errors.New("not a writer")))
			}
			length, err := stream.s.(io.Writer).Write([]byte(data))
			return join(vm, NewStatus(err), Int(length))
		}),

		Text("flush"): Func(func(vm *VM, aa *Args) *Args {
			stream := aa.get(0).(Stream)
			vm.da(aa)
			return join(vm, NewStatus(stream.Flush()))
		}),

		Text("close"): Func(func(vm *VM, aa *Args) *Args {
			stream := aa.get(0).(Stream)
			vm.da(aa)
			return join(vm, NewStatus(stream.Close()))
		}),
	})
	protoStream.meta = protoDef
}

const S42 Text = Text("init")
const S1 Text = Text("stdin")
const S5 Text = Text("write")
const S49 Text = Text("processes")
const S44 Text = Text("process")
const S57 Text = Text("serve")
const S31 Text = Text("open")
const S34 Text = Text("slurp")
const S45 Text = Text("to")
const S32 Text = Text("readall")
const S41 Text = Text("job")
const S33 Text = Text("close")
const S55 Text = Text("contains")
const S51 Text = Text("drop")
const S29 Text = Text("readrune")
const S48 Text = Text("name")
const S43 Text = Text("host")
const S53 Text = Text("dispatch")
const S19 Text = Text("set")
const S36 Text = Text("id")
const S25 Text = Text("lock")
const S40 Text = Text("todo")
const S7 Text = Text("flush")
const S8 Text = Text("type")
const S12 Text = Text("min")
const S14 Text = Text("push")
const S37 Text = Text("run")
const S4 Text = Text("stdout")
const S9 Text = Text("len")
const S59 Text = Text("text")
const S24 Text = Text("read")
const S52 Text = Text("send")
const S17 Text = Text("shift")
const S20 Text = Text("get")
const S26 Text = Text("jobs")
const S47 Text = Text("body")
const S50 Text = Text("exit")
const S54 Text = Text("devices")
const S6 Text = Text("join")
const S15 Text = Text("pop")
const S35 Text = Text("group")
const S21 Text = Text("keys")
const S28 Text = Text("queue")
const S11 Text = Text("max")
const S18 Text = Text("extend")
const S3 Text = Text("string")
const S10 Text = Text("iterate")
const S22 Text = Text("ticker")
const S23 Text = Text("stop")
const S46 Text = Text("from")
const S56 Text = Text("json")
const S58 Text = Text("websocket")
const S2 Text = Text("remove")
const S13 Text = Text("insert")
const S39 Text = Text("wait")
const S27 Text = Text("channel")
const S38 Text = Text("log")
const S16 Text = Text("shove")
const S30 Text = Text("readline")

func main() {

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

	{
		var Nprint Any
		noop(Nprint)
		var NKB Any
		noop(NKB)
		var NMB Any
		noop(NMB)
		var Nsuper Any
		noop(Nsuper)
		var NGHz Any
		noop(NGHz)
		var Ncreate Any
		noop(Ncreate)
		var Nplay Any
		noop(Nplay)
		var Nmessage Any
		noop(Nmessage)
		var Nnil Any
		noop(Nnil)
		var NKHz Any
		noop(NKHz)
		var NTHz Any
		noop(NTHz)
		var Naddress Any
		noop(Naddress)
		var Nmap Any
		noop(Nmap)
		var NMHz Any
		noop(NMHz)
		var Nstring Any
		noop(Nstring)
		var Nservices Any
		noop(Nservices)
		var Ncommon Any
		noop(Ncommon)
		var Ninteger Any
		noop(Ninteger)
		var Nstream Any
		noop(Nstream)
		var Ntrue Any
		noop(Ntrue)
		var NGB Any
		noop(NGB)
		var Ndecimal Any
		noop(Ndecimal)
		var Nlist Any
		noop(Nlist)
		var Nstarnix Any
		noop(Nstarnix)
		var Nchassis Any
		noop(Nchassis)
		var Nfalse Any
		noop(Nfalse)
		var NTB Any
		noop(NTB)
		var Nnode Any
		noop(Nnode)
		func() Any { a := one(vm, call(vm, Ngetprototype, join(vm, Nnil))); Nsuper = a; return a }()
		func() Any { a := one(vm, call(vm, Ngetprototype, join(vm, Int(0)))); Ninteger = a; return a }()
		func() Any { a := one(vm, call(vm, Ngetprototype, join(vm, Int(0)))); Ndecimal = a; return a }()
		func() Any { a := one(vm, call(vm, Ngetprototype, join(vm, Text("")))); Nstring = a; return a }()
		func() Any { a := one(vm, call(vm, Ngetprototype, join(vm, NewList([]Any{})))); Nlist = a; return a }()
		func() Any {
			a := one(vm, call(vm, Ngetprototype, join(vm, NewMap(MapData{}))))
			Nmap = a
			return a
		}()
		func() Any {
			a := one(vm, call(vm, Ngetprototype, join(vm, find(Nio, S1 /* stdin */))))
			Nstream = a
			return a
		}()
		func() Any { a := Bool(lt(Int(0), Int(1))); Ntrue = a; return a }()
		func() Any { a := Bool(lt(Int(1), Int(0))); Nfalse = a; return a }()
		func() Any {
			a := one(vm, func() *Args {
				t, m := method(NewList([]Any{}), S2 /* remove */)
				return call(vm, m, join(vm, t, Int(0)))
			}())
			Nnil = a
			return a
		}()
		func() Any {
			a := one(vm, Func(func(vm *VM, aa *Args) *Args {
				Na := aa.agg(0)
				noop(Na)
				vm.da(aa)
				{
					var Ns Any
					noop(Ns)
					loop(func() {
						it := iterate(Na)
						for {
							aa := it(vm, nil)
							if aa.get(0) == nil {
								vm.da(aa)
								break
							}
							vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
								Ni := aa.get(0)
								noop(Ni)
								Nv := aa.get(1)
								noop(Nv)
								vm.da(aa)
								{
									func() Any {
										a := one(vm, func() *Args {
											t, m := method(Nv, S3 /* string */)
											return call(vm, m, join(vm, t, nil))
										}())
										Ns = a
										return a
									}()
									if !eq(one(vm, call(vm, Ntype, join(vm, Ns))), Text("text")) {
										{
											func() Any { a := one(vm, call(vm, find(Nsuper, S3 /* string */), join(vm, Nv))); Ns = a; return a }()
										}
									}
									func() Any { a := Ns; store(Na, Ni, a); return a }()
								}
								return nil
							}), aa))
						}
					})
					vm.da(func() *Args {
						t, m := method(find(Nio, S4 /* stdout */), S5 /* write */)
						return call(vm, m, join(vm, t, concat(one(vm, func() *Args {
							t, m := method(Na, S6 /* join */)
							return call(vm, m, join(vm, t, Text("\t")))
						}()), Text("\n"))))
					}())
					vm.da(func() *Args {
						t, m := method(find(Nio, S4 /* stdout */), S7 /* flush */)
						return call(vm, m, join(vm, t, nil))
					}())
				}
				return nil
			}))
			Nprint = a
			return a
		}()
		vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
			vm.da(aa)
			{
				var Np Any
				noop(Np)
				func() Any { a := Nsuper; Np = a; return a }()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nany := aa.get(0)
						noop(Nany)
						vm.da(aa)
						{
							return join(vm, call(vm, Ntype, join(vm, Nany)))
						}
						return nil
					}))
					store(Np, S8 /* type */, a)
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nany := aa.get(0)
						noop(Nany)
						vm.da(aa)
						{
							return join(vm, length(Nany))
						}
						return nil
					}))
					store(Np, S9 /* len */, a)
					return a
				}()
			}
			return nil
		}), join(vm, nil)))
		vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
			vm.da(aa)
			{
				var Np Any
				noop(Np)
				func() Any { a := Ninteger; Np = a; return a }()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nlimit := aa.get(0)
						noop(Nlimit)
						vm.da(aa)
						{
							var Ni Any
							noop(Ni)
							var Nn Any
							noop(Nn)
							func() Any { a := Int(0); Ni = a; return a }()
							return join(vm, Func(func(vm *VM, aa *Args) *Args {
								vm.da(aa)
								{
									if lt(Ni, Nlimit) {
										{
											vm.da(func() *Args { aa := join(vm, Ni, add(Ni, Int(1))); Nn = aa.get(0); Ni = aa.get(1); return aa }())
											return join(vm, Nn)
										}
									}
								}
								return nil
							}))
						}
						return nil
					}))
					store(Np, S10 /* iterate */, a)
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Na := aa.get(0)
						noop(Na)
						Nb := aa.get(1)
						noop(Nb)
						vm.da(aa)
						{
							return join(vm, func() Any {
								var a Any
								a = func() Any {
									var a Any
									a = Bool(gt(Na, Nb))
									if truth(a) {
										var b Any
										b = Na
										if truth(b) {
											return b
										}
									}
									return nil
								}()
								if !truth(a) {
									a = Nb
								}
								return a
							}())
						}
						return nil
					}))
					store(Np, S11 /* max */, a)
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Na := aa.get(0)
						noop(Na)
						Nb := aa.get(1)
						noop(Nb)
						vm.da(aa)
						{
							return join(vm, func() Any {
								var a Any
								a = func() Any {
									var a Any
									a = Bool(lt(Na, Nb))
									if truth(a) {
										var b Any
										b = Na
										if truth(b) {
											return b
										}
									}
									return nil
								}()
								if !truth(a) {
									a = Nb
								}
								return a
							}())
						}
						return nil
					}))
					store(Np, S12 /* min */, a)
					return a
				}()
			}
			return nil
		}), join(vm, nil)))
		vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
			vm.da(aa)
			{
				var Np Any
				noop(Np)
				func() Any { a := Ndecimal; Np = a; return a }()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Na := aa.get(0)
						noop(Na)
						Nb := aa.get(1)
						noop(Nb)
						vm.da(aa)
						{
							return join(vm, func() Any {
								var a Any
								a = func() Any {
									var a Any
									a = Bool(gt(Na, Nb))
									if truth(a) {
										var b Any
										b = Na
										if truth(b) {
											return b
										}
									}
									return nil
								}()
								if !truth(a) {
									a = Nb
								}
								return a
							}())
						}
						return nil
					}))
					store(Np, S11 /* max */, a)
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Na := aa.get(0)
						noop(Na)
						Nb := aa.get(1)
						noop(Nb)
						vm.da(aa)
						{
							return join(vm, func() Any {
								var a Any
								a = func() Any {
									var a Any
									a = Bool(lt(Na, Nb))
									if truth(a) {
										var b Any
										b = Na
										if truth(b) {
											return b
										}
									}
									return nil
								}()
								if !truth(a) {
									a = Nb
								}
								return a
							}())
						}
						return nil
					}))
					store(Np, S12 /* min */, a)
					return a
				}()
			}
			return nil
		}), join(vm, nil)))
		vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
			vm.da(aa)
			{
				var Np Any
				noop(Np)
				func() Any { a := Nlist; Np = a; return a }()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nlist := aa.get(0)
						noop(Nlist)
						Nval := aa.get(1)
						noop(Nval)
						vm.da(aa)
						{
							return join(vm, func() *Args {
								t, m := method(Nlist, S13 /* insert */)
								return call(vm, m, join(vm, t, length(Nlist), Nval))
							}())
						}
						return nil
					}))
					store(Np, S14 /* push */, a)
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nlist := aa.get(0)
						noop(Nlist)
						vm.da(aa)
						{
							return join(vm, func() *Args {
								t, m := method(Nlist, S2 /* remove */)
								return call(vm, m, join(vm, t, sub(one(vm, length(Nlist)), Int(1))))
							}())
						}
						return nil
					}))
					store(Np, S15 /* pop */, a)
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nlist := aa.get(0)
						noop(Nlist)
						Nval := aa.get(1)
						noop(Nval)
						vm.da(aa)
						{
							return join(vm, func() *Args {
								t, m := method(Nlist, S13 /* insert */)
								return call(vm, m, join(vm, t, Int(0), Nval))
							}())
						}
						return nil
					}))
					store(Np, S16 /* shove */, a)
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nlist := aa.get(0)
						noop(Nlist)
						vm.da(aa)
						{
							return join(vm, func() *Args {
								t, m := method(Nlist, S2 /* remove */)
								return call(vm, m, join(vm, t, Int(0)))
							}())
						}
						return nil
					}))
					store(Np, S17 /* shift */, a)
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nlist := aa.get(0)
						noop(Nlist)
						vm.da(aa)
						{
							var Ni Any
							noop(Ni)
							var Nn Any
							noop(Nn)
							func() Any { a := Int(0); Ni = a; return a }()
							return join(vm, Func(func(vm *VM, aa *Args) *Args {
								vm.da(aa)
								{
									if lt(Ni, one(vm, length(Nlist))) {
										{
											vm.da(func() *Args { aa := join(vm, Ni, add(Ni, Int(1))); Nn = aa.get(0); Ni = aa.get(1); return aa }())
											return join(vm, Nn, field(Nlist, Nn))
										}
									}
								}
								return nil
							}))
						}
						return nil
					}))
					store(Np, S10 /* iterate */, a)
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nlist := aa.get(0)
						noop(Nlist)
						Nsize := aa.get(1)
						noop(Nsize)
						Ndef := aa.get(2)
						noop(Ndef)
						vm.da(aa)
						{
							loop(func() {
								for lt(one(vm, length(Nlist)), Nsize) {
									vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
										vm.da(aa)
										{
											vm.da(func() *Args {
												t, m := method(Nlist, S14 /* push */)
												return call(vm, m, join(vm, t, Ndef))
											}())
										}
										return nil
									}), nil))
								}
							})
							return join(vm, Nlist)
						}
						return nil
					}))
					store(Np, S18 /* extend */, a)
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nlist := aa.get(0)
						noop(Nlist)
						Npos := aa.get(1)
						noop(Npos)
						Nval := aa.get(2)
						noop(Nval)
						vm.da(aa)
						{
							func() Any { a := Nval; store(Nlist, Npos, a); return a }()
						}
						return nil
					}))
					store(Np, S19 /* set */, a)
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nlist := aa.get(0)
						noop(Nlist)
						Npos := aa.get(1)
						noop(Npos)
						vm.da(aa)
						{
							return join(vm, field(Nlist, Npos))
						}
						return nil
					}))
					store(Np, S20 /* get */, a)
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Na := aa.get(0)
						noop(Na)
						Nb := aa.get(1)
						noop(Nb)
						vm.da(aa)
						{
							return join(vm, func() Any {
								var a Any
								a = func() Any {
									var a Any
									a = Bool(gt(one(vm, length(Na)), one(vm, length(Nb))))
									if truth(a) {
										var b Any
										b = Na
										if truth(b) {
											return b
										}
									}
									return nil
								}()
								if !truth(a) {
									a = Nb
								}
								return a
							}())
						}
						return nil
					}))
					store(Np, S11 /* max */, a)
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Na := aa.get(0)
						noop(Na)
						Nb := aa.get(1)
						noop(Nb)
						vm.da(aa)
						{
							return join(vm, func() Any {
								var a Any
								a = func() Any {
									var a Any
									a = Bool(lt(one(vm, length(Na)), one(vm, length(Nb))))
									if truth(a) {
										var b Any
										b = Na
										if truth(b) {
											return b
										}
									}
									return nil
								}()
								if !truth(a) {
									a = Nb
								}
								return a
							}())
						}
						return nil
					}))
					store(Np, S12 /* min */, a)
					return a
				}()
			}
			return nil
		}), join(vm, nil)))
		vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
			vm.da(aa)
			{
				var Np Any
				noop(Np)
				func() Any { a := Nmap; Np = a; return a }()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nm := aa.get(0)
						noop(Nm)
						vm.da(aa)
						{
							var Nkeys Any
							noop(Nkeys)
							var Nn Any
							noop(Nn)
							var Nkey Any
							noop(Nkey)
							var Ni Any
							noop(Ni)
							func() Any { a := Int(0); Ni = a; return a }()
							func() Any {
								a := one(vm, func() *Args {
									t, m := method(Nm, S21 /* keys */)
									return call(vm, m, join(vm, t, nil))
								}())
								Nkeys = a
								return a
							}()
							return join(vm, Func(func(vm *VM, aa *Args) *Args {
								vm.da(aa)
								{
									if lt(Ni, one(vm, length(Nkeys))) {
										{
											vm.da(func() *Args { aa := join(vm, Ni, add(Ni, Int(1))); Nn = aa.get(0); Ni = aa.get(1); return aa }())
											func() Any { a := one(vm, field(Nkeys, Nn)); Nkey = a; return a }()
											return join(vm, Nkey, field(Nm, Nkey))
										}
									}
								}
								return nil
							}))
						}
						return nil
					}))
					store(Np, S10 /* iterate */, a)
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nm := aa.get(0)
						noop(Nm)
						Npos := aa.get(1)
						noop(Npos)
						Nval := aa.get(2)
						noop(Nval)
						vm.da(aa)
						{
							func() Any { a := Nval; store(Nm, Npos, a); return a }()
						}
						return nil
					}))
					store(Np, S19 /* set */, a)
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nm := aa.get(0)
						noop(Nm)
						Npos := aa.get(1)
						noop(Npos)
						vm.da(aa)
						{
							return join(vm, field(Nm, Npos))
						}
						return nil
					}))
					store(Np, S20 /* get */, a)
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Na := aa.get(0)
						noop(Na)
						Nb := aa.get(1)
						noop(Nb)
						vm.da(aa)
						{
							return join(vm, func() Any {
								var a Any
								a = func() Any {
									var a Any
									a = Bool(gt(one(vm, length(Na)), one(vm, length(Nb))))
									if truth(a) {
										var b Any
										b = Na
										if truth(b) {
											return b
										}
									}
									return nil
								}()
								if !truth(a) {
									a = Nb
								}
								return a
							}())
						}
						return nil
					}))
					store(Np, S11 /* max */, a)
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Na := aa.get(0)
						noop(Na)
						Nb := aa.get(1)
						noop(Nb)
						vm.da(aa)
						{
							return join(vm, func() Any {
								var a Any
								a = func() Any {
									var a Any
									a = Bool(lt(one(vm, length(Na)), one(vm, length(Nb))))
									if truth(a) {
										var b Any
										b = Na
										if truth(b) {
											return b
										}
									}
									return nil
								}()
								if !truth(a) {
									a = Nb
								}
								return a
							}())
						}
						return nil
					}))
					store(Np, S12 /* min */, a)
					return a
				}()
			}
			return nil
		}), join(vm, nil)))
		vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
			vm.da(aa)
			{
				var Nti Any
				noop(Nti)
				var Ntick Any
				noop(Ntick)
				func() Any {
					a := one(vm, call(vm, find(Ntime, S22 /* ticker */), join(vm, Int(1000000))))
					Nti = a
					return a
				}()
				func() Any { a := one(vm, call(vm, Ngetprototype, join(vm, Nti))); Ntick = a; return a }()
				vm.da(func() *Args {
					t, m := method(Nti, S23 /* stop */)
					return call(vm, m, join(vm, t, nil))
				}())
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nt := aa.get(0)
						noop(Nt)
						vm.da(aa)
						{
							var Ni Any
							noop(Ni)
							var Nn Any
							noop(Nn)
							func() Any { a := Int(0); Ni = a; return a }()
							return join(vm, Func(func(vm *VM, aa *Args) *Args {
								vm.da(aa)
								{
									vm.da(func() *Args { aa := join(vm, Ni, add(Ni, Int(1))); Nn = aa.get(0); Ni = aa.get(1); return aa }())
									return join(vm, Nn, func() *Args {
										t, m := method(Nt, S24 /* read */)
										return call(vm, m, join(vm, t, nil))
									}())
								}
								return nil
							}))
						}
						return nil
					}))
					store(Ntick, S10 /* iterate */, a)
					return a
				}()
			}
			return nil
		}), join(vm, nil)))
		vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
			vm.da(aa)
			{
				var NprotoQueue Any
				noop(NprotoQueue)
				func() Any {
					a := one(vm, NewMap(MapData{
						S24 /* read */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
							Nq := aa.get(0)
							noop(Nq)
							vm.da(aa)
							{
								var Njob Any
								noop(Njob)
								vm.da(func() *Args {
									t, m := method(find(Nq, S25 /* lock */), S5 /* write */)
									return call(vm, m, join(vm, t, nil))
								}())
								func() Any {
									a := one(vm, func() *Args {
										t, m := method(find(Nq, S26 /* jobs */), S17 /* shift */)
										return call(vm, m, join(vm, t, nil))
									}())
									Njob = a
									return a
								}()
								vm.da(func() *Args {
									t, m := method(find(Nq, S25 /* lock */), S24 /* read */)
									return call(vm, m, join(vm, t, nil))
								}())
								return join(vm, Njob)
							}
							return nil
						})),
						S5 /* write */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
							Nq := aa.get(0)
							noop(Nq)
							Nfn := aa.get(1)
							noop(Nfn)
							vm.da(aa)
							{
								vm.da(func() *Args {
									t, m := method(find(Nq, S25 /* lock */), S5 /* write */)
									return call(vm, m, join(vm, t, nil))
								}())
								vm.da(func() *Args {
									t, m := method(find(Nq, S26 /* jobs */), S14 /* push */)
									return call(vm, m, join(vm, t, Nfn))
								}())
								vm.da(func() *Args {
									t, m := method(find(Nq, S25 /* lock */), S24 /* read */)
									return call(vm, m, join(vm, t, nil))
								}())
							}
							return nil
						})),
						S10 /* iterate */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
							Nq := aa.get(0)
							noop(Nq)
							vm.da(aa)
							{
								var Njobs Any
								noop(Njobs)
								vm.da(func() *Args {
									t, m := method(find(Nq, S25 /* lock */), S5 /* write */)
									return call(vm, m, join(vm, t, nil))
								}())
								func() Any { a := one(vm, find(Nq, S26 /* jobs */)); Njobs = a; return a }()
								func() Any { a := one(vm, NewList([]Any{})); store(Nq, S26 /* jobs */, a); return a }()
								vm.da(func() *Args {
									t, m := method(find(Nq, S25 /* lock */), S24 /* read */)
									return call(vm, m, join(vm, t, nil))
								}())
								return join(vm, Func(func(vm *VM, aa *Args) *Args {
									vm.da(aa)
									{
										return join(vm, func() *Args {
											t, m := method(Njobs, S17 /* shift */)
											return call(vm, m, join(vm, t, nil))
										}())
									}
									return nil
								}))
							}
							return nil
						}))}))
					NprotoQueue = a
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						vm.da(aa)
						{
							return join(vm, call(vm, Nsetprototype, join(vm, NewMap(MapData{
								S25 /* lock */ : one(vm, call(vm, find(Nsync, S27 /* channel */), join(vm, Int(1)))),
								S26 /* jobs */ : one(vm, NewList([]Any{}))}), NprotoQueue)))
						}
						return nil
					}))
					store(Nsync, S28 /* queue */, a)
					return a
				}()
			}
			return nil
		}), join(vm, nil)))
		vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
			vm.da(aa)
			{
				var Np Any
				noop(Np)
				func() Any { a := Nstream; Np = a; return a }()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Ns := aa.get(0)
						noop(Ns)
						vm.da(aa)
						{
							var Nline Any
							noop(Nline)
							var Nok Any
							noop(Nok)
							var Nc Any
							noop(Nc)
							func() Any { a := one(vm, NewList([]Any{})); Nline = a; return a }()
							loop(func() {
								for truth(one(vm, func() *Args {
									aa := join(vm, func() *Args {
										t, m := method(Ns, S29 /* readrune */)
										return call(vm, m, join(vm, t, nil))
									}())
									Nok = aa.get(0)
									Nc = aa.get(1)
									return aa
								}())) {
									vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
										vm.da(aa)
										{
											vm.da(func() *Args {
												t, m := method(Nline, S14 /* push */)
												return call(vm, m, join(vm, t, Nc))
											}())
											if !truth(one(vm, func() Any {
												var a Any
												a = Nc
												if !truth(a) {
													a = Bool(eq(Nc, Rune('\n')))
												}
												return a
											}())) {
												{
													loopbreak()
												}
											}
										}
										return nil
									}), nil))
								}
							})
							return join(vm, Bool(gt(one(vm, length(Nline)), Int(0))), func() *Args {
								t, m := method(Nline, S6 /* join */)
								return call(vm, m, join(vm, t, nil))
							}())
						}
						return nil
					}))
					store(Np, S30 /* readline */, a)
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Ns := aa.get(0)
						noop(Ns)
						vm.da(aa)
						{
							var Ndone Any
							noop(Ndone)
							var Nok Any
							noop(Nok)
							var Nline Any
							noop(Nline)
							func() Any { a := Nfalse; Ndone = a; return a }()
							return join(vm, Func(func(vm *VM, aa *Args) *Args {
								vm.da(aa)
								{
									if truth(one(vm, func() *Args {
										aa := join(vm, func() *Args {
											t, m := method(Ns, S30 /* readline */)
											return call(vm, m, join(vm, t, nil))
										}())
										Nok = aa.get(0)
										Nline = aa.get(1)
										return aa
									}())) {
										{
											return join(vm, Nline)
										}
									}
									return join(vm, Nnil)
								}
								return nil
							}))
						}
						return nil
					}))
					store(Np, S10 /* iterate */, a)
					return a
				}()
			}
			return nil
		}), join(vm, nil)))
		vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
			vm.da(aa)
			{
				var Np Any
				noop(Np)
				func() Any { a := Nio; Np = a; return a }()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Npath := aa.get(0)
						noop(Npath)
						vm.da(aa)
						{
							var Nok Any
							noop(Nok)
							var Nfile Any
							noop(Nfile)
							var Ncontent Any
							noop(Ncontent)
							if truth(one(vm, func() *Args {
								aa := join(vm, call(vm, find(Nio, S31 /* open */), join(vm, Npath, Text("r"))))
								Nok = aa.get(0)
								Nfile = aa.get(1)
								return aa
							}())) {
								{
									if truth(one(vm, func() *Args {
										aa := join(vm, func() *Args {
											t, m := method(Nfile, S32 /* readall */)
											return call(vm, m, join(vm, t, nil))
										}())
										Nok = aa.get(0)
										Ncontent = aa.get(1)
										return aa
									}())) {
										{
											vm.da(func() *Args {
												t, m := method(Nfile, S33 /* close */)
												return call(vm, m, join(vm, t, nil))
											}())
											return join(vm, Ncontent)
										}
									}
								}
							}
							return join(vm, Nnil)
						}
						return nil
					}))
					store(Np, S34 /* slurp */, a)
					return a
				}()
			}
			return nil
		}), join(vm, nil)))
		vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
			vm.da(aa)
			{
				var Np Any
				noop(Np)
				func() Any {
					a := one(vm, call(vm, Ngetprototype, join(vm, call(vm, find(Nsync, S27 /* channel */), join(vm, Int(1))))))
					Np = a
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nchan := aa.get(0)
						noop(Nchan)
						vm.da(aa)
						{
							return join(vm, Func(func(vm *VM, aa *Args) *Args {
								vm.da(aa)
								{
									return join(vm, func() *Args {
										t, m := method(Nchan, S24 /* read */)
										return call(vm, m, join(vm, t, nil))
									}())
								}
								return nil
							}))
						}
						return nil
					}))
					store(Np, S10 /* iterate */, a)
					return a
				}()
			}
			return nil
		}), join(vm, nil)))
		func() Any { a := Int(1024); NKB = a; return a }()
		func() Any { a := one(vm, mul(NKB, Int(1024))); NMB = a; return a }()
		func() Any { a := one(vm, mul(NMB, Int(1024))); NGB = a; return a }()
		func() Any { a := one(vm, mul(NGB, Int(1024))); NTB = a; return a }()
		func() Any { a := Int(1000); NKHz = a; return a }()
		func() Any { a := one(vm, mul(NKHz, Int(1000))); NMHz = a; return a }()
		func() Any { a := one(vm, mul(NMHz, Int(1000))); NGHz = a; return a }()
		func() Any { a := one(vm, mul(NGHz, Int(1000))); NTHz = a; return a }()
		func() Any {
			a := one(vm, call(vm, Func(func(vm *VM, aa *Args) *Args {
				vm.da(aa)
				{
					var Napi Any
					noop(Napi)
					var Ngroup Any
					noop(Ngroup)
					var Nseq Any
					noop(Nseq)
					var Nerr Any
					noop(Nerr)
					func() Any {
						a := one(vm, NewMap(MapData{}))
						Napi = a
						return a
					}()
					func() Any { a := one(vm, call(vm, find(Nsync, S35 /* group */), join(vm, nil))); Ngroup = a; return a }()
					func() Any {
						a := one(vm, call(vm, find(Nsync, S27 /* channel */), join(vm, Int(0))))
						Nseq = a
						return a
					}()
					func() Any {
						a := one(vm, Func(func(vm *VM, aa *Args) *Args {
							vm.da(aa)
							{
								return join(vm, func() *Args {
									t, m := method(Nseq, S24 /* read */)
									return call(vm, m, join(vm, t, nil))
								}())
							}
							return nil
						}))
						store(Napi, S36 /* id */, a)
						return a
					}()
					vm.da(func() *Args {
						t, m := method(Ngroup, S37 /* run */)
						return call(vm, m, join(vm, t, Func(func(vm *VM, aa *Args) *Args {
							vm.da(aa)
							{
								var Nn Any
								noop(Nn)
								func() Any { a := Int(1); Nn = a; return a }()
								loop(func() {
									for truth(one(vm, func() *Args {
										t, m := method(Nseq, S5 /* write */)
										return call(vm, m, join(vm, t, Nn))
									}())) {
										vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
											vm.da(aa)
											{
												func() Any { a := add(Nn, Int(1)); Nn = a; return a }()
											}
											return nil
										}), nil))
									}
								})
							}
							return nil
						})))
					}())
					func() Any {
						a := one(vm, call(vm, find(Nsync, S27 /* channel */), join(vm, Int(100))))
						Nerr = a
						return a
					}()
					func() Any {
						a := one(vm, Func(func(vm *VM, aa *Args) *Args {
							Nmsg := aa.get(0)
							noop(Nmsg)
							vm.da(aa)
							{
								vm.da(func() *Args {
									t, m := method(Nerr, S5 /* write */)
									return call(vm, m, join(vm, t, Nmsg))
								}())
							}
							return nil
						}))
						store(Napi, S38 /* log */, a)
						return a
					}()
					vm.da(func() *Args {
						t, m := method(Ngroup, S37 /* run */)
						return call(vm, m, join(vm, t, Func(func(vm *VM, aa *Args) *Args {
							vm.da(aa)
							{
								loop(func() {
									it := iterate(Nerr)
									for {
										aa := it(vm, nil)
										if aa.get(0) == nil {
											vm.da(aa)
											break
										}
										vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
											Nentry := aa.get(0)
											noop(Nentry)
											vm.da(aa)
											{
												vm.da(call(vm, Nlog, join(vm, Nentry)))
											}
											return nil
										}), aa))
									}
								})
							}
							return nil
						})))
					}())
					func() Any {
						a := one(vm, Func(func(vm *VM, aa *Args) *Args {
							vm.da(aa)
							{
								vm.da(func() *Args {
									t, m := method(Nseq, S33 /* close */)
									return call(vm, m, join(vm, t, nil))
								}())
								vm.da(func() *Args {
									t, m := method(Nerr, S33 /* close */)
									return call(vm, m, join(vm, t, nil))
								}())
								vm.da(func() *Args {
									t, m := method(Ngroup, S39 /* wait */)
									return call(vm, m, join(vm, t, nil))
								}())
							}
							return nil
						}))
						store(Napi, S23 /* stop */, a)
						return a
					}()
					return join(vm, Napi)
				}
				return nil
			}), join(vm, nil)))
			Nservices = a
			return a
		}()
		func() Any {
			a := one(vm, NewMap(MapData{
				S26 /* jobs */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
					Nself := aa.get(0)
					noop(Nself)
					vm.da(aa)
					{
						loop(func() {
							it := iterate(one(vm, find(Nself, S40 /* todo */)))
							for {
								aa := it(vm, nil)
								if aa.get(0) == nil {
									vm.da(aa)
									break
								}
								vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
									Nfn := aa.get(0)
									noop(Nfn)
									vm.da(aa)
									{
										vm.da(call(vm, Nfn, join(vm, Nself)))
									}
									return nil
								}), aa))
							}
						})
					}
					return nil
				})),
				S41 /* job */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
					Nself := aa.get(0)
					noop(Nself)
					Nfn := aa.get(1)
					noop(Nfn)
					vm.da(aa)
					{
						vm.da(func() *Args {
							t, m := method(find(Nself, S40 /* todo */), S5 /* write */)
							return call(vm, m, join(vm, t, Nfn))
						}())
					}
					return nil
				})),
				S42 /* init */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
					Nself := aa.get(0)
					noop(Nself)
					vm.da(aa)
					{
						func() Any {
							a := one(vm, call(vm, find(Nservices, S36 /* id */), join(vm, nil)))
							store(Nself, S36 /* id */, a)
							return a
						}()
						func() Any {
							a := one(vm, call(vm, find(Nsync, S28 /* queue */), join(vm, nil)))
							store(Nself, S40 /* todo */, a)
							return a
						}()
					}
					return nil
				}))}))
			Ncommon = a
			return a
		}()
		func() Any {
			a := one(vm, call(vm, Func(func(vm *VM, aa *Args) *Args {
				vm.da(aa)
				{
					var Nproto Any
					noop(Nproto)
					func() Any {
						a := one(vm, NewMap(MapData{
							S43 /* host */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
								Nself := aa.get(0)
								noop(Nself)
								vm.da(aa)
								{
									return join(vm, field(Nself, Int(0)))
								}
								return nil
							})),
							S44 /* process */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
								Nself := aa.get(0)
								noop(Nself)
								vm.da(aa)
								{
									return join(vm, field(Nself, Int(1)))
								}
								return nil
							}))}))
						Nproto = a
						return a
					}()
					return join(vm, Func(func(vm *VM, aa *Args) *Args {
						Nhost := aa.get(0)
						noop(Nhost)
						Nprocess := aa.get(1)
						noop(Nprocess)
						vm.da(aa)
						{
							return join(vm, call(vm, Nsetprototype, join(vm, NewList([]Any{Nhost, Nprocess}), Nproto)))
						}
						return nil
					}))
				}
				return nil
			}), join(vm, nil)))
			Naddress = a
			return a
		}()
		func() Any {
			a := one(vm, call(vm, Func(func(vm *VM, aa *Args) *Args {
				vm.da(aa)
				{
					var Nproto Any
					noop(Nproto)
					func() Any {
						a := one(vm, NewMap(MapData{
							S45 /* to */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
								Nself := aa.get(0)
								noop(Nself)
								vm.da(aa)
								{
									return join(vm, field(Nself, Int(0)))
								}
								return nil
							})),
							S46 /* from */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
								Nself := aa.get(0)
								noop(Nself)
								vm.da(aa)
								{
									return join(vm, field(Nself, Int(1)))
								}
								return nil
							})),
							S47 /* body */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
								Nself := aa.get(0)
								noop(Nself)
								vm.da(aa)
								{
									return join(vm, field(Nself, Int(2)))
								}
								return nil
							}))}))
						Nproto = a
						return a
					}()
					return join(vm, Func(func(vm *VM, aa *Args) *Args {
						Nto := aa.get(0)
						noop(Nto)
						Nfrom := aa.get(1)
						noop(Nfrom)
						Nbody := aa.get(2)
						noop(Nbody)
						vm.da(aa)
						{
							return join(vm, call(vm, Nsetprototype, join(vm, NewList([]Any{Nto, Nfrom, Nbody}), Nproto)))
						}
						return nil
					}))
				}
				return nil
			}), join(vm, nil)))
			Nmessage = a
			return a
		}()
		func() Any {
			a := one(vm, NewMap(MapData{
				S43 /* host */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
					Nself := aa.get(0)
					noop(Nself)
					vm.da(aa)
					{
						return join(vm, find(Nself, S48 /* name */))
					}
					return nil
				})),
				S37 /* run */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
					Nself := aa.get(0)
					noop(Nself)
					Nname := aa.get(1)
					noop(Nname)
					Nqueue := aa.get(2)
					noop(Nqueue)
					vm.da(aa)
					{
						vm.da(func() *Args {
							t, m := method(Nself, S41 /* job */)
							return call(vm, m, join(vm, t, Func(func(vm *VM, aa *Args) *Args {
								vm.da(aa)
								{
									if lt(one(vm, length(find(Nself, S49 /* processes */))), Int(100)) {
										{
											func() Any { a := Nqueue; store(find(Nself, S49 /* processes */), Nname, a); return a }()
										}
									}
								}
								return nil
							})))
						}())
					}
					return nil
				})),
				S50 /* exit */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
					Nself := aa.get(0)
					noop(Nself)
					Nname := aa.get(1)
					noop(Nname)
					vm.da(aa)
					{
						vm.da(func() *Args {
							t, m := method(Nself, S41 /* job */)
							return call(vm, m, join(vm, t, Func(func(vm *VM, aa *Args) *Args {
								vm.da(aa)
								{
									vm.da(func() *Args {
										t, m := method(find(Nself, S49 /* processes */), S51 /* drop */)
										return call(vm, m, join(vm, t, Nname))
									}())
								}
								return nil
							})))
						}())
					}
					return nil
				})),
				S52 /* send */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
					Nself := aa.get(0)
					noop(Nself)
					Nname := aa.get(1)
					noop(Nname)
					Nmessage := aa.get(2)
					noop(Nmessage)
					vm.da(aa)
					{
						vm.da(func() *Args {
							t, m := method(Nself, S41 /* job */)
							return call(vm, m, join(vm, t, Func(func(vm *VM, aa *Args) *Args {
								vm.da(aa)
								{
									var Nqueue Any
									noop(Nqueue)
									if truth(one(vm, func() Any { a := one(vm, field(find(Nself, S49 /* processes */), Nname)); Nqueue = a; return a }())) {
										{
											vm.da(func() *Args {
												t, m := method(Nqueue, S5 /* write */)
												return call(vm, m, join(vm, t, Nmessage))
											}())
										}
									}
								}
								return nil
							})))
						}())
					}
					return nil
				})),
				S42 /* init */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
					Nself := aa.get(0)
					noop(Nself)
					vm.da(aa)
					{
						vm.da(func() *Args {
							t, m := method(Nself, S38 /* log */)
							return call(vm, m, join(vm, t, Text("boot")))
						}())
						func() Any {
							a := one(vm, NewMap(MapData{
								S53 /* dispatch */ : one(vm, call(vm, Func(func(vm *VM, aa *Args) *Args {
									vm.da(aa)
									{
										var Nqueue Any
										noop(Nqueue)
										var Nrun Any
										noop(Nrun)
										func() Any { a := one(vm, call(vm, find(Nsync, S28 /* queue */), join(vm, nil))); Nqueue = a; return a }()
										func() Any {
											a := one(vm, Func(func(vm *VM, aa *Args) *Args {
												vm.da(aa)
												{
													var Npq Any
													noop(Npq)
													var Nto Any
													noop(Nto)
													var Nprocess Any
													noop(Nprocess)
													vm.da(func() *Args {
														t, m := method(Nself, S41 /* job */)
														return call(vm, m, join(vm, t, Nrun))
													}())
													vm.da(func() *Args {
														t, m := method(Nself, S38 /* log */)
														return call(vm, m, join(vm, t, Text("dispatch")))
													}())
													loop(func() {
														it := iterate(Nqueue)
														for {
															aa := it(vm, nil)
															if aa.get(0) == nil {
																vm.da(aa)
																break
															}
															vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
																Nmsg := aa.get(0)
																noop(Nmsg)
																vm.da(aa)
																{
																	if truth(one(vm, func() Any {
																		var a Any
																		a = func() Any {
																			var a Any
																			a = func() Any {
																				a := one(vm, func() *Args {
																					t, m := method(Nmsg, S45 /* to */)
																					return call(vm, m, join(vm, t, nil))
																				}())
																				Nto = a
																				return a
																			}()
																			if truth(a) {
																				var b Any
																				b = Bool(eq(one(vm, func() *Args {
																					t, m := method(Nto, S43 /* host */)
																					return call(vm, m, join(vm, t, nil))
																				}()), one(vm, func() *Args {
																					t, m := method(Nself, S43 /* host */)
																					return call(vm, m, join(vm, t, nil))
																				}())))
																				if truth(b) {
																					return b
																				}
																			}
																			return nil
																		}()
																		if truth(a) {
																			var b Any
																			b = func() Any {
																				a := one(vm, func() *Args {
																					t, m := method(Nto, S44 /* process */)
																					return call(vm, m, join(vm, t, nil))
																				}())
																				Nprocess = a
																				return a
																			}()
																			if truth(b) {
																				return b
																			}
																		}
																		return nil
																	}())) {
																		{
																			if truth(one(vm, func() Any { a := one(vm, field(find(Nself, S49 /* processes */), Nprocess)); Npq = a; return a }())) {
																				{
																					vm.da(func() *Args {
																						t, m := method(Npq, S5 /* write */)
																						return call(vm, m, join(vm, t, Nmsg))
																					}())
																				}
																			} else {
																				{
																					vm.da(func() *Args {
																						t, m := method(Nself, S38 /* log */)
																						return call(vm, m, join(vm, t, Text("drop message"), Nmessage))
																					}())
																				}
																			}
																		}
																	}
																}
																return nil
															}), aa))
														}
													})
												}
												return nil
											}))
											Nrun = a
											return a
										}()
										vm.da(func() *Args {
											t, m := method(Nself, S41 /* job */)
											return call(vm, m, join(vm, t, Nrun))
										}())
										return join(vm, Nqueue)
									}
									return nil
								}), join(vm, nil)))}))
							store(Nself, S49 /* processes */, a)
							return a
						}()
					}
					return nil
				})),
				S38 /* log */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
					Nself := aa.get(0)
					noop(Nself)
					Nmsg := aa.agg(1)
					noop(Nmsg)
					vm.da(aa)
					{
						vm.da(call(vm, find(Nservices, S38 /* log */), join(vm, func() *Args {
							t, m := method(Nmsg, S16 /* shove */)
							return call(vm, m, join(vm, t, find(Nself, S36 /* id */)))
						}())))
					}
					return nil
				}))}))
			Nstarnix = a
			return a
		}()
		func() Any {
			a := one(vm, NewMap(MapData{
				S42 /* init */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
					Nself := aa.get(0)
					noop(Nself)
					vm.da(aa)
					{
						func() Any { a := one(vm, NewList([]Any{})); store(Nself, S54 /* devices */, a); return a }()
					}
					return nil
				})),
				S55 /* contains */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
					Nself := aa.get(0)
					noop(Nself)
					Ndev := aa.get(1)
					noop(Ndev)
					vm.da(aa)
					{
						var Nid Any
						noop(Nid)
						var Nfound Any
						noop(Nfound)
						func() Any {
							a := one(vm, func() *Args {
								t, m := method(Ndev, S36 /* id */)
								return call(vm, m, join(vm, t, nil))
							}())
							Nid = a
							return a
						}()
						func() Any { a := Nfalse; Nfound = a; return a }()
						loop(func() {
							it := iterate(one(vm, find(Nself, S54 /* devices */)))
							for {
								aa := it(vm, nil)
								if aa.get(0) == nil {
									vm.da(aa)
									break
								}
								vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
									Ni := aa.get(0)
									noop(Ni)
									Nd := aa.get(1)
									noop(Nd)
									vm.da(aa)
									{
										if eq(one(vm, func() *Args {
											t, m := method(Nd, S36 /* id */)
											return call(vm, m, join(vm, t, nil))
										}()), Nid) {
											{
												func() Any { a := Ntrue; Nfound = a; return a }()
												loopbreak()
											}
										}
									}
									return nil
								}), aa))
							}
						})
						return join(vm, Nfound)
					}
					return nil
				})),
				S13 /* insert */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
					Nself := aa.get(0)
					noop(Nself)
					Ndev := aa.get(1)
					noop(Ndev)
					vm.da(aa)
					{
						vm.da(call(vm, find(Nself, S41 /* job */), join(vm, Func(func(vm *VM, aa *Args) *Args {
							vm.da(aa)
							{
								if !truth(one(vm, func() *Args {
									t, m := method(Nself, S55 /* contains */)
									return call(vm, m, join(vm, t, Ndev))
								}())) {
									{
										vm.da(func() *Args {
											t, m := method(find(Nself, S54 /* devices */), S14 /* push */)
											return call(vm, m, join(vm, t, Ndev))
										}())
									}
								}
							}
							return nil
						}))))
					}
					return nil
				})),
				S2 /* remove */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
					Nself := aa.get(0)
					noop(Nself)
					Ndev := aa.get(1)
					noop(Ndev)
					vm.da(aa)
					{
						vm.da(call(vm, find(Nself, S41 /* job */), join(vm, Func(func(vm *VM, aa *Args) *Args {
							vm.da(aa)
							{
								var Nid Any
								noop(Nid)
								func() Any {
									a := one(vm, func() *Args {
										t, m := method(Ndev, S36 /* id */)
										return call(vm, m, join(vm, t, nil))
									}())
									Nid = a
									return a
								}()
								loop(func() {
									it := iterate(one(vm, find(Nself, S54 /* devices */)))
									for {
										aa := it(vm, nil)
										if aa.get(0) == nil {
											vm.da(aa)
											break
										}
										vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
											Ni := aa.get(0)
											noop(Ni)
											Nd := aa.get(1)
											noop(Nd)
											vm.da(aa)
											{
												if eq(one(vm, func() *Args {
													t, m := method(Nd, S36 /* id */)
													return call(vm, m, join(vm, t, nil))
												}()), Nid) {
													{
														vm.da(func() *Args {
															t, m := method(find(Nself, S54 /* devices */), S2 /* remove */)
															return call(vm, m, join(vm, t, Ni))
														}())
														loopbreak()
													}
												}
											}
											return nil
										}), aa))
									}
								})
							}
							return nil
						}))))
					}
					return nil
				}))}))
			Nchassis = a
			return a
		}()
		func() Any {
			a := one(vm, Func(func(vm *VM, aa *Args) *Args {
				Nroles := aa.get(0)
				noop(Nroles)
				vm.da(aa)
				{
					var Nnode Any
					noop(Nnode)
					func() Any {
						a := one(vm, func() *Args {
							t, m := method(func() Any {
								var a Any
								a = Nroles
								if !truth(a) {
									a = NewList([]Any{})
								}
								return a
							}(), S16 /* shove */)
							return call(vm, m, join(vm, t, Ncommon))
						}())
						Nroles = a
						return a
					}()
					func() Any {
						a := one(vm, call(vm, Nsetprototype, join(vm, NewMap(MapData{}), Nroles)))
						Nnode = a
						return a
					}()
					loop(func() {
						it := iterate(Nroles)
						for {
							aa := it(vm, nil)
							if aa.get(0) == nil {
								vm.da(aa)
								break
							}
							vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
								Ni := aa.get(0)
								noop(Ni)
								Nrole := aa.get(1)
								noop(Nrole)
								vm.da(aa)
								{
									if truth(one(vm, find(Nrole, S42 /* init */))) {
										{
											vm.da(call(vm, find(Nrole, S42 /* init */), join(vm, Nnode)))
										}
									}
								}
								return nil
							}), aa))
						}
					})
					return join(vm, Nnode)
				}
				return nil
			}))
			Ncreate = a
			return a
		}()
		func() Any {
			a := one(vm, call(vm, Ncreate, join(vm, NewList([]Any{Nchassis, Nstarnix}))))
			Nnode = a
			return a
		}()
		vm.da(func() *Args {
			t, m := method(Nnode, S26 /* jobs */)
			return call(vm, m, join(vm, t, nil))
		}())
		vm.da(func() *Args {
			t, m := method(Nnode, S26 /* jobs */)
			return call(vm, m, join(vm, t, nil))
		}())
		vm.da(func() *Args {
			t, m := method(Nnode, S26 /* jobs */)
			return call(vm, m, join(vm, t, nil))
		}())
		func() Any {
			a := one(vm, Func(func(vm *VM, aa *Args) *Args {
				Ngroup := aa.get(0)
				noop(Ngroup)
				Ninbox := aa.get(1)
				noop(Ninbox)
				vm.da(aa)
				{
					var Noutbox Any
					noop(Noutbox)
					func() Any {
						a := one(vm, call(vm, find(Nsync, S27 /* channel */), join(vm, Int(8))))
						Noutbox = a
						return a
					}()
					vm.da(func() *Args {
						t, m := method(Ngroup, S37 /* run */)
						return call(vm, m, join(vm, t, Func(func(vm *VM, aa *Args) *Args {
							vm.da(aa)
							{
								var Ntry Any
								noop(Ntry)
								var Nmsg Any
								noop(Nmsg)
								loop(func() {
									it := iterate(Ninbox)
									for {
										aa := it(vm, nil)
										if aa.get(0) == nil {
											vm.da(aa)
											break
										}
										vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
											Nmsg := aa.get(0)
											noop(Nmsg)
											vm.da(aa)
											{
												if truth(one(vm, func() *Args {
													aa := join(vm, func() *Args {
														t, m := method(Nmsg, S56 /* json */)
														return call(vm, m, join(vm, t, nil))
													}())
													Ntry = aa.get(0)
													Nmsg = aa.get(1)
													return aa
												}())) {
													{
														vm.da(call(vm, Nprint, join(vm, Text("got"), Nmsg)))
														vm.da(func() *Args {
															t, m := method(Noutbox, S5 /* write */)
															return call(vm, m, join(vm, t, Text("ok")))
														}())
													}
												} else {
													{
														vm.da(call(vm, Nlog, join(vm, Ntry)))
													}
												}
											}
											return nil
										}), aa))
									}
								})
								vm.da(func() *Args {
									t, m := method(Noutbox, S33 /* close */)
									return call(vm, m, join(vm, t, nil))
								}())
							}
							return nil
						})))
					}())
					return join(vm, Noutbox)
				}
				return nil
			}))
			Nplay = a
			return a
		}()
		vm.da(call(vm, find(Nhttp, S57 /* serve */), join(vm, Text(":3000"), Text("static/"), NewMap(MapData{
			Text("/join"): one(vm, Func(func(vm *VM, aa *Args) *Args {
				Nreq := aa.get(0)
				noop(Nreq)
				vm.da(aa)
				{
					var Nok Any
					noop(Nok)
					var Nws Any
					noop(Nws)
					var Ngroup Any
					noop(Ngroup)
					var Ninbox Any
					noop(Ninbox)
					var Noutbox Any
					noop(Noutbox)
					if truth(one(vm, func() *Args {
						aa := join(vm, func() *Args {
							t, m := method(Nreq, S58 /* websocket */)
							return call(vm, m, join(vm, t, nil))
						}())
						Nok = aa.get(0)
						Nws = aa.get(1)
						return aa
					}())) {
						{
							func() Any { a := one(vm, call(vm, find(Nsync, S35 /* group */), join(vm, nil))); Ngroup = a; return a }()
							func() Any {
								a := one(vm, call(vm, find(Nsync, S27 /* channel */), join(vm, Int(8))))
								Ninbox = a
								return a
							}()
							func() Any { a := one(vm, call(vm, Nplay, join(vm, Ngroup, Ninbox))); Noutbox = a; return a }()
							vm.da(func() *Args {
								t, m := method(Ngroup, S37 /* run */)
								return call(vm, m, join(vm, t, Func(func(vm *VM, aa *Args) *Args {
									vm.da(aa)
									{
										var Nmode Any
										noop(Nmode)
										var Nmsg Any
										noop(Nmsg)
										var Ntry Any
										noop(Ntry)
										loop(func() {
											for truth(one(vm, func() *Args {
												aa := join(vm, func() *Args {
													t, m := method(Nws, S24 /* read */)
													return call(vm, m, join(vm, t, nil))
												}())
												Ntry = aa.get(0)
												Nmode = aa.get(1)
												Nmsg = aa.get(2)
												return aa
											}())) {
												vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
													vm.da(aa)
													{
														vm.da(func() *Args {
															t, m := method(Ninbox, S5 /* write */)
															return call(vm, m, join(vm, t, func() *Args {
																t, m := method(Nmsg, S59 /* text */)
																return call(vm, m, join(vm, t, nil))
															}()))
														}())
													}
													return nil
												}), nil))
											}
										})
										vm.da(call(vm, Nlog, join(vm, Ntry)))
										vm.da(func() *Args {
											t, m := method(Ninbox, S33 /* close */)
											return call(vm, m, join(vm, t, nil))
										}())
									}
									return nil
								})))
							}())
							vm.da(func() *Args {
								t, m := method(Ngroup, S37 /* run */)
								return call(vm, m, join(vm, t, Func(func(vm *VM, aa *Args) *Args {
									vm.da(aa)
									{
										loop(func() {
											it := iterate(Noutbox)
											for {
												aa := it(vm, nil)
												if aa.get(0) == nil {
													vm.da(aa)
													break
												}
												vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
													Nmsg := aa.get(0)
													noop(Nmsg)
													vm.da(aa)
													{
														vm.da(func() *Args {
															t, m := method(Nws, S5 /* write */)
															return call(vm, m, join(vm, t, find(Nws, S59 /* text */), Nmsg))
														}())
													}
													return nil
												}), aa))
											}
										})
									}
									return nil
								})))
							}())
							vm.da(func() *Args {
								t, m := method(Ngroup, S39 /* wait */)
								return call(vm, m, join(vm, t, nil))
							}())
						}
					}
				}
				return nil
			}))}))))
	}
}
