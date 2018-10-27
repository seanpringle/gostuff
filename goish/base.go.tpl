package main

import "fmt"
import "math"
import "strings"
import "strconv"
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

type Any interface {
	Type() string
	Lib() *Map
	String() string
}

type VM struct {
	argsC [8]*Args
	argsN int
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

func (aa *Args) Lib() *Map {
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
var protoStr *Map
var protoChan *Map
var protoGroup *Map
var protoInst *Map
var protoTick *Map
var protoByte *Map
var protoStream *Map

var libIO *Map
var libTime *Map
var libSync *Map

var onInit []func()

type Stringer = fmt.Stringer

type BoolIsh interface {
	Bool() Bool
}

type StrIsh interface {
	Str() Str
}

type ByteIsh interface {
	Byte() Byte
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
	return protoDef
}

type Str string

func (s Str) Str() Str {
	return s
}

func (s Str) String() string {
	return string(s)
}

func (s Str) Type() string {
	return "string"
}

func (s Str) Len() int64 { // characters
	l := int64(0)
	for range string(s) {
		l++
	}
	return l
}

func (s Str) Lib() *Map {
	return protoStr
}

func (s Str) Bool() Bool {
	return Bool(len(s) > 0)
}

func (s Str) Byte() Byte {
	return []byte(string(s))
}

type Byte []byte

func (b Byte) Type() string {
	return "bytes"
}

func (b Byte) Lib() *Map {
	return protoByte
}

func (b Byte) String() string {
	return hex.Dump([]byte(b))
}

func (b Byte) Str() Str {
	return Str(string(b))
}

func (b Byte) Len() int64 {
	return int64(len(b))
}

func (b Byte) Bool() Bool {
	return Bool(len(b) > 0)
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
	return protoInt
}

func (i Int) Bool() Bool {
	return Bool(int64(i) != 0)
}

func (i Int) Dec() Dec {
	return Dec(float64(int64(i)))
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

func (i *SInt) Lib() *Map {
	return protoInt
}

func (i *SInt) Bool() Bool {
	return i.i.Bool()
}

func (i *SInt) Dec() Dec {
	return Dec(float64(int64(i.i)))
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
	return protoDec
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

func (r Rune) Str() Str {
	return Str(r.String())
}

func (r Rune) Lib() *Map {
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

func (e Status) Lib() *Map {
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

func (t Instant) Lib() *Map {
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

func (t Ticker) Lib() *Map {
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
	meta Any
}

func NewMap(data MapData) *Map {
	return &Map{
		data: data,
		meta: protoMap,
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
			if l, is := t.meta.(*List); is {
				for _, i := range l.data {
					r := find(i, key)
					if r != nil {
						return r
					}
				}
			}
			if mt, is := t.meta.(*Map); is {
				return mt.Find(key)
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
				v = Str("map")
			}
			if _, is := v.(*List); is {
				v = Str("slice")
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
}

func NewList(data []Any) *List {
	return &List{data: data}
}

func (t *List) Lib() *Map {
	return protoList
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
			v = Str("map")
		}
		if _, is := v.(*List); is {
			v = Str("slice")
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

func (f Func) Lib() *Map {
	return protoDef
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

func (c *Chan) Lib() *Map {
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

func (g *Group) Lib() *Map {
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

func (s Stream) Lib() *Map {
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
	if n < IntCache {
		return &SInts[n]
	}
	return n
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
	if _, is := a.(StrIsh); is {
		return Str(tostring(a) + tostring(b))
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
	if as, is := a.(Str); is {
		if bs, is := b.(Str); is {
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
		return t.Lib().Find(key)
	} else {
		return protoDef.Find(key)
	}
	return nil
}

func field(t Any, key Any) Any {
	if t != nil {
		if _, is := t.(*Map); is {
			return find(t, key)
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

func trymethod(t Any, k string, def Any) Any {
	t, m := method(t, Str(k))
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
	return trymethod(s, "string", Str("(nil)")).(StrIsh).Str().String()
}

func tobool(s Any) bool {
	return bool(trymethod(s, "bool", Bool(false)).(BoolIsh).Bool())
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

var Nio *Map
var Ntime *Map
var Nsync *Map

var Ntype Any = Func(func(vm *VM, aa *Args) *Args {
	v := aa.get(0)
	vm.da(aa)
	if aa != nil {
		return join(vm, Str(v.Type()))
	}
	return join(vm, Str("nil"))
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
	if m, is := aa.get(0).(*Map); is {
		m.meta = aa.get(1)
		vm.da(aa)
		aa = vm.ga(1)
		aa.set(0, m)
		return aa
	}
	panic(fmt.Sprintf("cannot set prototype"))
})

var Ngetprototype Any = Func(func(vm *VM, aa *Args) *Args {
	v := aa.get(0)
	if v == nil {
		return join(vm, protoDef)
	}
	if m, is := v.(*Map); is {
		return join(vm, m.meta)
	}
	return join(vm, v.Lib())
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
		Str("string"): Func(func(vm *VM, aa *Args) *Args {
			a := aa.get(0)
			vm.da(aa)
			if a == nil {
				return join(vm, Str("(nil)"))
			}
			return join(vm, Str(a.String()))
		}),
		Str("status"): Func(func(vm *VM, aa *Args) *Args {
			a := aa.get(0)
			vm.da(aa)
			if a == nil {
				return join(vm, NewStatus(nil))
			}
			return join(vm, NewStatus(errors.New(tostring(a))))
		}),
	})
	protoInt = NewMap(MapData{})
	protoInt.meta = protoDef
	protoDec = NewMap(MapData{})
	protoDec.meta = protoDef
	protoMap = NewMap(MapData{
		Str("keys"): Func(func(vm *VM, aa *Args) *Args {
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
		Str("push"): Func(func(vm *VM, aa *Args) *Args {
			l := aa.get(0).(*List)
			v := aa.get(1)
			vm.da(aa)
			l.data = append(l.data, v)
			return join(vm, l)
		}),
		Str("pop"): Func(func(vm *VM, aa *Args) *Args {
			l := aa.get(0).(*List)
			n := len(l.data) - 1
			vm.da(aa)
			var v Any
			if len(l.data) < n {
				v = l.data[n]
				l.data = l.data[0:n]
			}
			return join(vm, v)
		}),
		Str("shove"): Func(func(vm *VM, aa *Args) *Args {
			l := aa.get(0).(*List)
			v := aa.get(1)
			vm.da(aa)
			l.data = append([]Any{v}, l.data...)
			return join(vm, l)
		}),
		Str("shift"): Func(func(vm *VM, aa *Args) *Args {
			l := aa.get(0).(*List)
			vm.da(aa)
			var v Any
			if len(l.data) > 0 {
				v = l.data[0]
				l.data = l.data[1:]
			}
			return join(vm, v)
		}),
		Str("join"): Func(func(vm *VM, aa *Args) *Args {
			l := aa.get(0).(*List)
			j := ifnil(aa.get(1), Str(""))
			vm.da(aa)
			var ls []string
			for _, s := range l.data {
				ls = append(ls, tostring(s))
			}
			return join(vm, Str(strings.Join(ls, tostring(j))))
		}),
		Str("sort"): Func(func(vm *VM, aa *Args) *Args {
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
		Str("read"): Func(func(vm *VM, aa *Args) *Args {
			c := aa.get(0).(*Chan)
			vm.da(aa)
			return join(vm, <-c.c)
		}),
		Str("write"): Func(func(vm *VM, aa *Args) *Args {
			c := aa.get(0).(*Chan)
			a := aa.get(1)
			vm.da(aa)
			c.c <- a
			return join(vm, NewStatus(nil))
		}),
		Str("close"): Func(func(vm *VM, aa *Args) *Args {
			c := aa.get(0).(*Chan)
			vm.da(aa)
			close(c.c)
			return nil
		}),
	})
	protoChan.meta = protoDef
	protoGroup = NewMap(MapData{
		Str("run"): Func(func(vm *VM, aa *Args) *Args {
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
		Str("wait"): Func(func(vm *VM, aa *Args) *Args {
			g := aa.get(0).(*Group)
			vm.da(aa)
			g.Wait()
			return join(vm, NewStatus(nil))
		}),
	})
	protoGroup.meta = protoDef
	protoStr = NewMap(MapData{
		Str("bytes"): Func(func(vm *VM, aa *Args) *Args {
			s := tostring(aa.get(0))
			vm.da(aa)
			return join(vm, Byte(s))
		}),
		Str("split"): Func(func(vm *VM, aa *Args) *Args {
			s := tostring(aa.get(0))
			j := tostring(aa.get(1))
			vm.da(aa)
			l := []Any{}
			for _, p := range strings.Split(s, j) {
				l = append(l, Str(p))
			}
			return join(vm, NewList(l))
		}),
	})
	protoStr.meta = protoDef
	protoTick = NewMap(MapData{
		Str("read"): Func(func(vm *VM, aa *Args) *Args {
			ti := aa.get(0).(Ticker)
			vm.da(aa)
			return join(vm, ti.Read())
		}),
		Str("stop"): Func(func(vm *VM, aa *Args) *Args {
			ti := aa.get(0).(Ticker)
			vm.da(aa)
			return join(vm, ti.Stop())
		}),
	})
	protoTick.meta = protoDef

	protoInst = protoDef

	libTime = NewMap(MapData{
		Str("ms"): Int(int64(time.Millisecond)),
		Str("ticker"): Func(func(vm *VM, aa *Args) *Args {
			d := int64(aa.get(0).(Int))
			vm.da(aa)
			return join(vm, NewTicker(time.Duration(d)))
		}),
	})
	Ntime = libTime

	libSync = NewMap(MapData{
		Str("group"): Func(func(vm *VM, aa *Args) *Args {
			vm.da(aa)
			return join(vm, NewGroup())
		}),
		Str("channel"): Func(func(vm *VM, aa *Args) *Args {
			n := int(aa.get(0).(IntIsh).Int())
			vm.da(aa)
			return join(vm, NewChan(n))
		}),
	})
	Nsync = libSync

	libIO = NewMap(MapData{
		Str("stdin"):  NewStream(os.Stdin),
		Str("stdout"): NewStream(os.Stdout),
		Str("stderr"): NewStream(os.Stderr),

		Str("open"): Func(func(vm *VM, aa *Args) *Args {
			path := tostring(aa.get(0))
			modes := tostring(aa.get(1))
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

	protoByte = protoDef

	protoStream = NewMap(MapData{

		Str("read"): Func(func(vm *VM, aa *Args) *Args {
			stream := aa.get(0).(Stream)
			limit := ifnil(aa.get(1), Int(1024*1024)).(IntIsh).Int()
			vm.da(aa)
			if _, is := stream.s.(io.Reader); !is {
				return join(vm, NewStatus(errors.New("not a reader")))
			}
			buff := make([]byte, int(limit))
			length, err := stream.s.(io.Reader).Read(buff)
			return join(vm, NewStatus(err), Byte(buff[:length]))
		}),

		Str("readall"): Func(func(vm *VM, aa *Args) *Args {
			stream := aa.get(0).(Stream)
			vm.da(aa)
			if _, is := stream.s.(io.Reader); !is {
				return join(vm, NewStatus(errors.New("not a reader")))
			}
			buff, err := ioutil.ReadAll(stream.s.(io.Reader))
			return join(vm, NewStatus(err), Byte(buff))
		}),

		Str("readrune"): Func(func(vm *VM, aa *Args) *Args {
			stream := aa.get(0).(Stream)
			vm.da(aa)
			if r, is := stream.s.(ReadRuner); is {
				char, _, err := r.ReadRune()
				return join(vm, NewStatus(err), Rune(char))
			}
			return join(vm, NewStatus(errors.New("not a reader")))
		}),

		Str("write"): Func(func(vm *VM, aa *Args) *Args {
			stream := aa.get(0).(Stream)
			data := aa.get(1).(ByteIsh).Byte()
			vm.da(aa)
			if _, is := stream.s.(io.Writer); !is {
				return join(vm, NewStatus(errors.New("not a writer")))
			}
			length, err := stream.s.(io.Writer).Write([]byte(data))
			return join(vm, NewStatus(err), Int(length))
		}),

		Str("flush"): Func(func(vm *VM, aa *Args) *Args {
			stream := aa.get(0).(Stream)
			vm.da(aa)
			return join(vm, NewStatus(stream.Flush()))
		}),

		Str("close"): Func(func(vm *VM, aa *Args) *Args {
			stream := aa.get(0).(Stream)
			vm.da(aa)
			return join(vm, NewStatus(stream.Close()))
		}),
	})
	protoStream.meta = protoDef
}
