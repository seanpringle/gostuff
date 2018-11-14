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
		f(vm, aa)
		g.g.Done()
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
	return Stream{s}
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
	return nil
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

const S30 Text = Text("queue")
const S40 Text = Text("g")
const S44 Text = Text("run")
const S3 Text = Text("string")
const S4 Text = Text("stdout")
const S24 Text = Text("ticker")
const S25 Text = Text("stop")
const S28 Text = Text("jobs")
const S36 Text = Text("slurp")
const S46 Text = Text("tom")
const S9 Text = Text("type")
const S11 Text = Text("iterate")
const S23 Text = Text("keys")
const S21 Text = Text("set")
const S22 Text = Text("get")
const S14 Text = Text("insert")
const S16 Text = Text("pop")
const S20 Text = Text("extend")
const S19 Text = Text("clear")
const S34 Text = Text("readall")
const S41 Text = Text("m")
const S43 Text = Text("group")
const S52 Text = Text("huge")
const S5 Text = Text("write")
const S10 Text = Text("len")
const S18 Text = Text("shift")
const S32 Text = Text("readline")
const S45 Text = Text("wait")
const S47 Text = Text("dick")
const S49 Text = Text("b")
const S51 Text = Text("match")
const S6 Text = Text("join")
const S13 Text = Text("min")
const S17 Text = Text("shove")
const S26 Text = Text("read")
const S35 Text = Text("close")
const S39 Text = Text("a")
const S42 Text = Text("split")
const S48 Text = Text("harry")
const S1 Text = Text("stdin")
const S2 Text = Text("remove")
const S8 Text = Text("stderr")
const S50 Text = Text("sort")
const S15 Text = Text("push")
const S29 Text = Text("channel")
const S33 Text = Text("open")
const S31 Text = Text("readrune")
const S37 Text = Text("c")
const S38 Text = Text("d")
const S7 Text = Text("flush")
const S12 Text = Text("max")
const S27 Text = Text("lock")

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
		var Nnil Any
		noop(Nnil)
		var Ninc Any
		noop(Ninc)
		var Nlen Any
		noop(Nlen)
		var Ng Any
		noop(Ng)
		var Ninteger Any
		noop(Ninteger)
		var Ndecimal Any
		noop(Ndecimal)
		var Nlist Any
		noop(Nlist)
		var Na Any
		noop(Na)
		var Nb Any
		noop(Nb)
		var Nc Any
		noop(Nc)
		var Nhi Any
		noop(Nhi)
		var Nblink Any
		noop(Nblink)
		var Nm Any
		noop(Nm)
		var Nl Any
		noop(Nl)
		var Nstring Any
		noop(Nstring)
		var Nmap Any
		noop(Nmap)
		var Ntrue Any
		noop(Ntrue)
		var Nprint Any
		noop(Nprint)
		var Nlog Any
		noop(Nlog)
		var Nis Any
		noop(Nis)
		var Nt Any
		noop(Nt)
		var Nf Any
		noop(Nf)
		var Nsuper Any
		noop(Nsuper)
		var Nstream Any
		noop(Nstream)
		var Nfalse Any
		noop(Nfalse)
		var Ns Any
		noop(Ns)
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
										store(Na, Ni, a)
										return a
									}()
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
		func() Any {
			a := one(vm, Func(func(vm *VM, aa *Args) *Args {
				Na := aa.agg(0)
				noop(Na)
				vm.da(aa)
				{
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
										store(Na, Ni, a)
										return a
									}()
								}
								return nil
							}), aa))
						}
					})
					vm.da(func() *Args {
						t, m := method(find(Nio, S8 /* stderr */), S5 /* write */)
						return call(vm, m, join(vm, t, concat(one(vm, func() *Args {
							t, m := method(Na, S6 /* join */)
							return call(vm, m, join(vm, t, Text("\t")))
						}()), Text("\n"))))
					}())
					vm.da(func() *Args {
						t, m := method(find(Nio, S8 /* stderr */), S7 /* flush */)
						return call(vm, m, join(vm, t, nil))
					}())
				}
				return nil
			}))
			Nlog = a
			return a
		}()
		func() Any {
			a := one(vm, Func(func(vm *VM, aa *Args) *Args {
				Nclass := aa.get(0)
				noop(Nclass)
				Nobject := aa.get(1)
				noop(Nobject)
				vm.da(aa)
				{
					var Nproto Any
					noop(Nproto)
					func() Any { a := one(vm, call(vm, Ngetprototype, join(vm, Nobject))); Nproto = a; return a }()
					if eq(one(vm, call(vm, Ntype, join(vm, Nproto))), Text("list")) {
						{
							loop(func() {
								it := iterate(Nproto)
								for {
									aa := it(vm, nil)
									if aa.get(0) == nil {
										vm.da(aa)
										break
									}
									vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
										Ni := aa.get(0)
										noop(Ni)
										Nitem := aa.get(1)
										noop(Nitem)
										vm.da(aa)
										{
											if eq(Nitem, Nclass) {
												{
													return join(vm, Ntrue)
												}
											}
										}
										return nil
									}), aa))
								}
							})
						}
					}
					return join(vm, Bool(eq(Nclass, Nproto)))
				}
				return nil
			}))
			Nis = a
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
					store(Np, S9 /* type */, a)
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
					store(Np, S10 /* len */, a)
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
					store(Np, S11 /* iterate */, a)
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
					store(Np, S12 /* max */, a)
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
					store(Np, S13 /* min */, a)
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
					store(Np, S12 /* max */, a)
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
					store(Np, S13 /* min */, a)
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
								t, m := method(Nlist, S14 /* insert */)
								return call(vm, m, join(vm, t, length(Nlist), Nval))
							}())
						}
						return nil
					}))
					store(Np, S15 /* push */, a)
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
					store(Np, S16 /* pop */, a)
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
								t, m := method(Nlist, S14 /* insert */)
								return call(vm, m, join(vm, t, Int(0), Nval))
							}())
						}
						return nil
					}))
					store(Np, S17 /* shove */, a)
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
					store(Np, S18 /* shift */, a)
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nlist := aa.get(0)
						noop(Nlist)
						vm.da(aa)
						{
							loop(func() {
								it := iterate(one(vm, length(Nlist)))
								for {
									aa := it(vm, nil)
									if aa.get(0) == nil {
										vm.da(aa)
										break
									}
									vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
										vm.da(aa)
										{
											vm.da(func() *Args {
												t, m := method(Nlist, S16 /* pop */)
												return call(vm, m, join(vm, t, nil))
											}())
										}
										return nil
									}), aa))
								}
							})
						}
						return nil
					}))
					store(Np, S19 /* clear */, a)
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
					store(Np, S11 /* iterate */, a)
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
												t, m := method(Nlist, S15 /* push */)
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
					store(Np, S20 /* extend */, a)
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
					store(Np, S21 /* set */, a)
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
					store(Np, S22 /* get */, a)
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
					store(Np, S12 /* max */, a)
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
					store(Np, S13 /* min */, a)
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
							var Nkey Any
							noop(Nkey)
							var Ni Any
							noop(Ni)
							var Nkeys Any
							noop(Nkeys)
							var Nn Any
							noop(Nn)
							func() Any { a := Int(0); Ni = a; return a }()
							func() Any {
								a := one(vm, func() *Args {
									t, m := method(Nm, S23 /* keys */)
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
					store(Np, S11 /* iterate */, a)
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
					store(Np, S21 /* set */, a)
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
					store(Np, S22 /* get */, a)
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
					store(Np, S12 /* max */, a)
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
					store(Np, S13 /* min */, a)
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
					a := one(vm, call(vm, find(Ntime, S24 /* ticker */), join(vm, Int(1000000))))
					Nti = a
					return a
				}()
				func() Any { a := one(vm, call(vm, Ngetprototype, join(vm, Nti))); Ntick = a; return a }()
				vm.da(func() *Args {
					t, m := method(Nti, S25 /* stop */)
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
										t, m := method(Nt, S26 /* read */)
										return call(vm, m, join(vm, t, nil))
									}())
								}
								return nil
							}))
						}
						return nil
					}))
					store(Ntick, S11 /* iterate */, a)
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
						S11 /* iterate */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
							Nq := aa.get(0)
							noop(Nq)
							vm.da(aa)
							{
								var Njobs Any
								noop(Njobs)
								vm.da(func() *Args {
									t, m := method(find(Nq, S27 /* lock */), S5 /* write */)
									return call(vm, m, join(vm, t, nil))
								}())
								func() Any { a := one(vm, find(Nq, S28 /* jobs */)); Njobs = a; return a }()
								func() Any { a := one(vm, NewList([]Any{})); store(Nq, S28 /* jobs */, a); return a }()
								vm.da(func() *Args {
									t, m := method(find(Nq, S27 /* lock */), S26 /* read */)
									return call(vm, m, join(vm, t, nil))
								}())
								return join(vm, Func(func(vm *VM, aa *Args) *Args {
									vm.da(aa)
									{
										return join(vm, func() *Args {
											t, m := method(Njobs, S18 /* shift */)
											return call(vm, m, join(vm, t, nil))
										}())
									}
									return nil
								}))
							}
							return nil
						})),
						S26 /* read */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
							Nq := aa.get(0)
							noop(Nq)
							vm.da(aa)
							{
								var Njob Any
								noop(Njob)
								vm.da(func() *Args {
									t, m := method(find(Nq, S27 /* lock */), S5 /* write */)
									return call(vm, m, join(vm, t, nil))
								}())
								func() Any {
									a := one(vm, func() *Args {
										t, m := method(find(Nq, S28 /* jobs */), S18 /* shift */)
										return call(vm, m, join(vm, t, nil))
									}())
									Njob = a
									return a
								}()
								vm.da(func() *Args {
									t, m := method(find(Nq, S27 /* lock */), S26 /* read */)
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
									t, m := method(find(Nq, S27 /* lock */), S5 /* write */)
									return call(vm, m, join(vm, t, nil))
								}())
								vm.da(func() *Args {
									t, m := method(find(Nq, S28 /* jobs */), S15 /* push */)
									return call(vm, m, join(vm, t, Nfn))
								}())
								vm.da(func() *Args {
									t, m := method(find(Nq, S27 /* lock */), S26 /* read */)
									return call(vm, m, join(vm, t, nil))
								}())
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
								S27 /* lock */ : one(vm, call(vm, find(Nsync, S29 /* channel */), join(vm, Int(1)))),
								S28 /* jobs */ : one(vm, NewList([]Any{}))}), NprotoQueue)))
						}
						return nil
					}))
					store(Nsync, S30 /* queue */, a)
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
							var Nc Any
							noop(Nc)
							var Nline Any
							noop(Nline)
							var Nok Any
							noop(Nok)
							func() Any { a := one(vm, NewList([]Any{})); Nline = a; return a }()
							loop(func() {
								for truth(one(vm, func() *Args {
									aa := join(vm, func() *Args {
										t, m := method(Ns, S31 /* readrune */)
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
												t, m := method(Nline, S15 /* push */)
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
					store(Np, S32 /* readline */, a)
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
											t, m := method(Ns, S32 /* readline */)
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
					store(Np, S11 /* iterate */, a)
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
								aa := join(vm, call(vm, find(Nio, S33 /* open */), join(vm, Npath, Text("r"))))
								Nok = aa.get(0)
								Nfile = aa.get(1)
								return aa
							}())) {
								{
									if truth(one(vm, func() *Args {
										aa := join(vm, func() *Args {
											t, m := method(Nfile, S34 /* readall */)
											return call(vm, m, join(vm, t, nil))
										}())
										Nok = aa.get(0)
										Ncontent = aa.get(1)
										return aa
									}())) {
										{
											vm.da(func() *Args {
												t, m := method(Nfile, S35 /* close */)
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
					store(Np, S36 /* slurp */, a)
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
					a := one(vm, call(vm, Ngetprototype, join(vm, call(vm, find(Nsync, S29 /* channel */), join(vm, Int(1))))))
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
										t, m := method(Nchan, S26 /* read */)
										return call(vm, m, join(vm, t, nil))
									}())
								}
								return nil
							}))
						}
						return nil
					}))
					store(Np, S11 /* iterate */, a)
					return a
				}()
			}
			return nil
		}), join(vm, nil)))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			aa := join(vm, call(vm, Nprint, join(vm, Int(1), Text("hi"))))
			Na = aa.get(0)
			Nb = aa.get(1)
			return aa
		}())))
		vm.da(call(vm, Nprint, join(vm, func() Any {
			var a Any
			a = func() Any {
				var a Any
				a = Int(1)
				if truth(a) {
					var b Any
					b = Int(0)
					if truth(b) {
						return b
					}
				}
				return nil
			}()
			if !truth(a) {
				a = Int(3)
			}
			return a
		}())))
		vm.da(call(vm, Nprint, join(vm, add(Int(5), Int(6)))))
		func() Any {
			a := one(vm, Func(func(vm *VM, aa *Args) *Args {
				Na := aa.get(0)
				noop(Na)
				vm.da(aa)
				{
					return join(vm, add(Na, Int(1)))
				}
				return nil
			}))
			Ninc = a
			return a
		}()
		vm.da(call(vm, Nprint, join(vm, call(vm, Ninc, join(vm, Int(42))))))
		vm.da(call(vm, Nprint, join(vm, func() Any {
			var a Any
			a = func() Any {
				var a Any
				a = Bool(eq(Na, Int(1)))
				if truth(a) {
					var b Any
					b = Int(7)
					if truth(b) {
						return b
					}
				}
				return nil
			}()
			if !truth(a) {
				a = Int(9)
			}
			return a
		}())))
		func() Any {
			a := one(vm, NewMap(MapData{
				Text("__*&^"): Int(2),
				S37 /* c */ : one(vm, NewMap(MapData{
					S38 /* d */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
						vm.da(aa)
						{
							return join(vm, Text("hello world"))
						}
						return nil
					}))})),
				S39 /* a */ : Int(1)}))
			Nt = a
			return a
		}()
		vm.da(call(vm, Nprint, join(vm, call(vm, find(one(vm, find(Nt, S37 /* c */)), S38 /* d */), join(vm, nil)))))
		func() Any { a := Int(42); store(Nt, S39 /* a */, a); return a }()
		vm.da(call(vm, Nprint, join(vm, Nt)))
		vm.da(call(vm, Nprint, join(vm, Text(""), func() *Args {
			t, m := method(Nt, S23 /* keys */)
			return call(vm, m, join(vm, t, nil))
		}())))
		vm.da(call(vm, Nprint, join(vm, add(one(vm, mul(Int(2), Int(2))), Int(3)))))
		func() Any {
			a := one(vm, NewMap(MapData{
				S40 /* g */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
					vm.da(aa)
					{
						return join(vm, Text("hello world"))
					}
					return nil
				}))}))
			Nt = a
			return a
		}()
		func() Any {
			a := one(vm, Func(func(vm *VM, aa *Args) *Args {
				Nself := aa.get(0)
				noop(Nself)
				vm.da(aa)
				{
					return join(vm, func() *Args {
						t, m := method(Nself, S40 /* g */)
						return call(vm, m, join(vm, t, nil))
					}())
				}
				return nil
			}))
			store(Nt, S41 /* m */, a)
			return a
		}()
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(Nt, S41 /* m */)
			return call(vm, m, join(vm, t, nil))
		}())))
		func() Any { a := Text("goodbye world"); Ns = a; return a }()
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(Ns, S10 /* len */)
			return call(vm, m, join(vm, t, nil))
		}())))
		vm.da(call(vm, Nprint, join(vm, call(vm, Ntype, join(vm, Ns)))))
		vm.da(call(vm, Nprint, join(vm, NewList([]Any{Int(1), Int(2), Int(7)}))))
		func() Any {
			a := one(vm, NewMap(MapData{}))
			Na = a
			return a
		}()
		vm.da(call(vm, Nprint, join(vm, Na)))
		vm.da(func() *Args {
			t, m := method(Na, S21 /* set */)
			return call(vm, m, join(vm, t, Text("1"), Int(1)))
		}())
		vm.da(call(vm, Nprint, join(vm, Na)))
		func() Any {
			a := one(vm, NewMap(MapData{}))
			Nb = a
			return a
		}()
		vm.da(func() *Args {
			t, m := method(Na, S21 /* set */)
			return call(vm, m, join(vm, t, Nb, Int(2)))
		}())
		vm.da(call(vm, Nprint, join(vm, Na)))
		vm.da(func() *Args {
			t, m := method(Nb, S21 /* set */)
			return call(vm, m, join(vm, t, Text("2"), Int(2)))
		}())
		vm.da(call(vm, Nprint, join(vm, Na)))
		func() Any { a := one(vm, NewList([]Any{Int(1), Int(2), Int(3)})); Nl = a; return a }()
		vm.da(call(vm, Nprint, join(vm, Nl)))
		vm.da(func() *Args {
			t, m := method(Nl, S15 /* push */)
			return call(vm, m, join(vm, t, Int(4)))
		}())
		vm.da(call(vm, Nprint, join(vm, Nl)))
		vm.da(call(vm, Nprint, join(vm, call(vm, Ngetprototype, join(vm, Nl)))))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(Nl, S16 /* pop */)
			return call(vm, m, join(vm, t, nil))
		}())))
		vm.da(call(vm, Nprint, join(vm, Nl)))
		vm.da(call(vm, Nprint, join(vm, concat(Text("a"), Text("b")))))
		func() Any { a := Text("hi"); Nlen = a; return a }()
		vm.da(call(vm, Nprint, join(vm, Text("yo"), func() *Args {
			t, m := method(Nl, S10 /* len */)
			return call(vm, m, join(vm, t, nil))
		}())))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(Text("a,b,c"), S42 /* split */)
			return call(vm, m, join(vm, t, Text(",")))
		}())))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(one(vm, func() *Args {
				t, m := method(Text("a,b,c"), S42 /* split */)
				return call(vm, m, join(vm, t, Text(",")))
			}()), S6 /* join */)
			return call(vm, m, join(vm, t, Text(":")))
		}())))
		func() Any { a := one(vm, call(vm, find(Nsync, S29 /* channel */), join(vm, Int(10)))); Nc = a; return a }()
		vm.da(func() *Args {
			t, m := method(Nc, S5 /* write */)
			return call(vm, m, join(vm, t, Int(1)))
		}())
		vm.da(func() *Args {
			t, m := method(Nc, S5 /* write */)
			return call(vm, m, join(vm, t, Int(2)))
		}())
		vm.da(func() *Args {
			t, m := method(Nc, S5 /* write */)
			return call(vm, m, join(vm, t, Int(3)))
		}())
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(Nc, S26 /* read */)
			return call(vm, m, join(vm, t, nil))
		}())))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(Nc, S26 /* read */)
			return call(vm, m, join(vm, t, nil))
		}())))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(Nc, S26 /* read */)
			return call(vm, m, join(vm, t, nil))
		}())))
		func() Any {
			a := one(vm, Func(func(vm *VM, aa *Args) *Args {
				Ng := aa.get(0)
				noop(Ng)
				vm.da(aa)
				{
					vm.da(call(vm, Nprint, join(vm, Text("hi"))))
				}
				return nil
			}))
			Nhi = a
			return a
		}()
		func() Any { a := one(vm, call(vm, find(Nsync, S43 /* group */), join(vm, nil))); Ng = a; return a }()
		vm.da(func() *Args {
			t, m := method(Ng, S44 /* run */)
			return call(vm, m, join(vm, t, Nhi))
		}())
		vm.da(func() *Args {
			t, m := method(Ng, S44 /* run */)
			return call(vm, m, join(vm, t, Nhi))
		}())
		vm.da(func() *Args {
			t, m := method(Ng, S44 /* run */)
			return call(vm, m, join(vm, t, Nhi))
		}())
		vm.da(func() *Args {
			t, m := method(Ng, S45 /* wait */)
			return call(vm, m, join(vm, t, nil))
		}())
		vm.da(call(vm, Nprint, join(vm, Text("done"))))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(Nb, S22 /* get */)
			return call(vm, m, join(vm, t, Text("hi")))
		}())))
		vm.da(call(vm, Nprint, join(vm, func() Any {
			var a Any
			a = func() Any {
				var a Any
				a = Ntrue
				if truth(a) {
					var b Any
					b = Text("yes")
					if truth(b) {
						return b
					}
				}
				return nil
			}()
			if !truth(a) {
				a = Text("no")
			}
			return a
		}())))
		loop(func() {
			it := iterate(Int(10))
			for {
				aa := it(vm, nil)
				if aa.get(0) == nil {
					vm.da(aa)
					break
				}
				vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
					Ni := aa.get(0)
					noop(Ni)
					vm.da(aa)
					{
						vm.da(call(vm, Nprint, join(vm, Ni)))
					}
					return nil
				}), aa))
			}
		})
		loop(func() {
			it := iterate(one(vm, NewList([]Any{Int(1), Int(2), Int(3)})))
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
						vm.da(call(vm, Nprint, join(vm, Ni, Text(":"), Nv)))
					}
					return nil
				}), aa))
			}
		})
		loop(func() {
			it := iterate(one(vm, NewMap(MapData{
				S46 /* tom */ :   Int(1),
				S47 /* dick */ :  Int(2),
				S48 /* harry */ : Int(43)})))
			for {
				aa := it(vm, nil)
				if aa.get(0) == nil {
					vm.da(aa)
					break
				}
				vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
					Nk := aa.get(0)
					noop(Nk)
					Nv := aa.get(1)
					noop(Nv)
					vm.da(aa)
					{
						vm.da(call(vm, Nprint, join(vm, Nk, Text("=>"), Nv)))
					}
					return nil
				}), aa))
			}
		})
		func() Any { a := Int(1); Na = a; return a }()
		vm.da(call(vm, Nprint, join(vm, func() Any { a := add(Na, Int(1)); Na = a; return a }())))
		vm.da(call(vm, Nprint, join(vm, func() Any { a := add(Na, Int(1)); Na = a; return a }())))
		vm.da(call(vm, Nprint, join(vm, func() Any { a := add(Na, Int(1)); Na = a; return a }())))
		loop(func() {
			it := iterate(Int(10))
			for {
				aa := it(vm, nil)
				if aa.get(0) == nil {
					vm.da(aa)
					break
				}
				vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
					Ni := aa.get(0)
					noop(Ni)
					vm.da(aa)
					{
						if eq(Ni, Int(5)) {
							{
								loopbreak()
							}
						}
						vm.da(call(vm, Nprint, join(vm, Ni)))
					}
					return nil
				}), aa))
			}
		})
		func() Any { a := one(vm, call(vm, find(Nsync, S30 /* queue */), join(vm, nil))); Nblink = a; return a }()
		vm.da(func() *Args {
			t, m := method(Nblink, S5 /* write */)
			return call(vm, m, join(vm, t, Func(func(vm *VM, aa *Args) *Args {
				vm.da(aa)
				{
					vm.da(call(vm, Nprint, join(vm, Text("hello world"))))
				}
				return nil
			})))
		}())
		vm.da(func() *Args {
			t, m := method(Nblink, S5 /* write */)
			return call(vm, m, join(vm, t, Func(func(vm *VM, aa *Args) *Args {
				vm.da(aa)
				{
					vm.da(call(vm, Nprint, join(vm, Text("hello world"))))
				}
				return nil
			})))
		}())
		vm.da(func() *Args {
			t, m := method(Nblink, S5 /* write */)
			return call(vm, m, join(vm, t, Func(func(vm *VM, aa *Args) *Args {
				vm.da(aa)
				{
					vm.da(call(vm, Nprint, join(vm, Text("hello world"))))
				}
				return nil
			})))
		}())
		vm.da(call(vm, Nprint, join(vm, Text("and..."))))
		loop(func() {
			it := iterate(Nblink)
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
						vm.da(call(vm, Nprint, join(vm, Nfn, call(vm, Nfn, join(vm, nil)))))
					}
					return nil
				}), aa))
			}
		})
		func() Any { a := one(vm, NewList([]Any{Int(1), Int(2), Int(3)})); Nl = a; return a }()
		vm.da(call(vm, Nprint, join(vm, field(Nl, Int(0)))))
		func() Any {
			a := one(vm, NewMap(MapData{
				S49 /* b */ : one(vm, NewMap(MapData{
					S37 /* c */ : Int(4)})),
				S39 /* a */ : Int(1)}))
			Nm = a
			return a
		}()
		vm.da(call(vm, Nprint, join(vm, field(field(Nm, Text("b")), Text("c")))))
		func() Any { a := Int(5); store(field(Nm, Text("b")), Text("c"), a); return a }()
		vm.da(call(vm, Nprint, join(vm, field(field(Nm, Text("b")), Text("c")))))
		vm.da(call(vm, Nprint, join(vm, Text("length"), length(Nl), length(Nm))))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(Int(0), S12 /* max */)
			return call(vm, m, join(vm, t, Int(2)))
		}())))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(NewList([]Any{Int(2), Int(4), Int(6), Int(8), Int(3)}), S50 /* sort */)
			return call(vm, m, join(vm, t, Func(func(vm *VM, aa *Args) *Args {
				Na := aa.get(0)
				noop(Na)
				Nb := aa.get(1)
				noop(Nb)
				vm.da(aa)
				{
					return join(vm, Bool(lt(Na, Nb)))
				}
				return nil
			})))
		}())))
		vm.da(call(vm, Nprint, join(vm, Text(`a
`+"`"+`multi`+"`"+`
line
string
`))))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(Text("abc"), S51 /* match */)
			return call(vm, m, join(vm, t, Text("[aeiou]")))
		}())))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(Text("abc"), S51 /* match */)
			return call(vm, m, join(vm, t, Text("[aeiou]")))
		}())))
		vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
			vm.da(aa)
			{
				vm.da(call(vm, Nprint, join(vm, Text("hi"))))
			}
			return nil
		}), join(vm, nil)))
		vm.da(call(vm, Nprint, join(vm, find(one(vm, call(vm, Ngetprototype, join(vm, Int(0)))), S52 /* huge */))))
		vm.da(call(vm, Nprint, join(vm, find(one(vm, call(vm, Ngetprototype, join(vm, Dec(1)))), S52 /* huge */))))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(NewList([]Any{}), S20 /* extend */)
			return call(vm, m, join(vm, t, Int(3)))
		}())))
		vm.da(call(vm, Nprint, join(vm, mul(add(Int(1), Int(2)), Int(3)))))
		vm.da(call(vm, Nprint, join(vm, mul(Int(3), add(Int(1), Int(2))))))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(NewList([]Any{Int(1), Int(2), Int(3)}), S14 /* insert */)
			return call(vm, m, join(vm, t, Int(1), Int(7)))
		}())))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(NewList([]Any{Int(1), Int(2), Int(3)}), S14 /* insert */)
			return call(vm, m, join(vm, t, Int(0), Int(7)))
		}())))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(NewList([]Any{Int(1), Int(2), Int(3)}), S14 /* insert */)
			return call(vm, m, join(vm, t, Int(4), Int(7)))
		}())))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(NewList([]Any{Int(1), Int(2), Int(3)}), S2 /* remove */)
			return call(vm, m, join(vm, t, Int(0)))
		}())))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(NewList([]Any{Int(1), Int(2), Int(3)}), S2 /* remove */)
			return call(vm, m, join(vm, t, Int(1)))
		}())))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(NewList([]Any{Int(1), Int(2), Int(3)}), S2 /* remove */)
			return call(vm, m, join(vm, t, Int(3)))
		}())))
		func() Any { a := one(vm, NewList([]Any{Int(1), Int(2), Int(3)})); Nl = a; return a }()
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(Nl, S2 /* remove */)
			return call(vm, m, join(vm, t, Int(2)))
		}(), Nl)))
		func() Any {
			a := one(vm, NewMap(MapData{
				S49 /* b */ : one(vm, NewList([]Any{Int(1), Int(2), Int(3)}))}))
			Na = a
			return a
		}()
		vm.da(call(vm, Nprint, join(vm, length(Na))))
		if truth(one(vm, func() Any {
			var a Any
			a = func() Any {
				a := one(vm, NewMap(MapData{
					S49 /* b */ : Int(1)}))
				Na = a
				return a
			}()
			if truth(a) {
				var b Any
				b = Bool(eq(one(vm, find(Na, S49 /* b */)), Int(1)))
				if truth(b) {
					return b
				}
			}
			return nil
		}())) {
			{
				vm.da(call(vm, Nprint, join(vm, Text("yes"))))
			}
		}
		func() Any {
			a := one(vm, Func(func(vm *VM, aa *Args) *Args {
				Na := aa.agg(0)
				noop(Na)
				vm.da(aa)
				{
					vm.da(call(vm, Nprint, join(vm, Na)))
				}
				return nil
			}))
			Nf = a
			return a
		}()
		vm.da(call(vm, Nprint, join(vm, extract(vm, one(vm, NewList([]Any{Int(1), Int(2), Int(3)}))))))
		defer func() { call(vm, Nprint, join(vm, Text("deferred!"))) }()
		defer func() { call(vm, Nprint, join(vm, Text("deferred! 2"))) }()
		loop(func() {
			it := iterate(Int(3))
			for {
				aa := it(vm, nil)
				if aa.get(0) == nil {
					vm.da(aa)
					break
				}
				vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
					Ni := aa.get(0)
					noop(Ni)
					vm.da(aa)
					{
						defer func() { call(vm, Nprint, join(vm, Text("defer"), Ni)) }()
					}
					return nil
				}), aa))
			}
		})
	}
}
