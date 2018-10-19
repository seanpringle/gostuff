package main

import "fmt"
import "math"
import "strings"
import "strconv"
import "sync"
import "time"
import "log"
import "os"
import "runtime/pprof"

type Any interface {
	Type() string
	Lib() *Map
	String() string
}

type VM struct {
	args []*Args
}

func (vm *VM) ga(n int) *Args {
	if n > 16 {
		panic("maximum of 16 arguments per call")
	}
	l := len(vm.args)
	if l > 0 {
		aa := vm.args[l-1]
		vm.args = vm.args[:l-1]
		aa.used = n
		//for i := 0; i < n; i++ {
		//	aa.set(i, nil)
		//}
		return aa
	}
	return &Args{used: n}
}

func (vm *VM) da(a *Args) {
	if a != nil {
		for i := 0; i < a.used; i++ {
			a.cells[i] = nil
		}
		if len(vm.args) < 16 {
			vm.args = append(vm.args, a)
		}
	}
}

type Args struct {
	vm    *VM
	used  int
	cells [16]Any
}

func (aa *Args) Type() string {
	return "args"
}

func (aa *Args) String() string {
	return "args"
}

func (aa *Args) Lib() *Map {
	return libDef
}

func (aa *Args) get(i int) Any {
	if aa != nil && i < aa.used {
		return aa.cells[i]
	}
	return nil
}

func (aa *Args) set(i int, v Any) {
	if aa != nil && i < aa.used {
		aa.cells[i] = v
	}
}

func (aa *Args) len() int {
	if aa != nil {
		return aa.used
	}
	return 0
}

var libDef *Map
var libMap *Map
var libList *Map
var libStr *Map
var libChan *Map
var libGroup *Map
var libTime *Map
var libTick *Map

type Stringer = fmt.Stringer

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
	return libDef
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
	if as, is := a.(StrIsh); is {
		return Str{as.Str().s + tostring(b)}
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
			for j := 0; j < aa.len(); j++ {
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
		vm := &VM{}
		return call(vm, m, join(vm, t)).get(0)
	}
	return def
}

func tostring(s Any) string {
	return trymethod(s, "string", Str{"nil"}).(Str).s
}

func noop(a Any) {
}

var Nprint Any = Func(func(vm *VM, aa *Args) *Args {
	parts := []string{}
	for i := 0; i < aa.len(); i++ {
		parts = append(parts, tostring(aa.get(i)))
	}
	fmt.Printf("%s\n", strings.Join(parts, " "))
	return aa
})

var Nchan Any = Func(func(vm *VM, aa *Args) *Args {
	n := int64(aa.get(0).(IntIsh).Int())
	c := make(chan Any, int(n))
	return join(vm, Chan(c))
})

var Ngroup Any = Func(func(vm *VM, aa *Args) *Args {
	return join(vm, NewGroup())
})

var Ntime *Map

var Ntype Any = Func(func(vm *VM, aa *Args) *Args {
	return join(vm, Str{aa.get(0).Type()})
})

var Nsetprototype Any = Func(func(vm *VM, aa *Args) *Args {
	if m, is := aa.get(0).(*Map); is {
		m.meta = aa.get(1)
		return nil
	}
	panic(fmt.Sprintf("cannot set prototype"))
})

var Ngetprototype Any = Func(func(vm *VM, aa *Args) *Args {
	if m, is := aa.get(0).(*Map); is {
		return join(vm, m.meta)
	}
	return join(vm, aa.get(0).Lib())
})

func init() {
	libDef = NewMap(MapData{
		Str{"len"}: Func(func(vm *VM, aa *Args) *Args {
			s := aa.get(0).(LenIsh)
			return join(vm, Int(s.Len()))
		}),
		Str{"type"}: Ntype,
		Str{"string"}: Func(func(vm *VM, aa *Args) *Args {
			return join(vm, Str{fmt.Sprintf("%v", aa.get(0))})
		}),
	})
	libMap = NewMap(MapData{
		Str{"keys"}: Func(func(vm *VM, aa *Args) *Args {
			keys := []Any{}
			for k, _ := range aa.get(0).(*Map).data {
				keys = append(keys, k)
			}
			return join(vm, NewList(keys))
		}),
		Str{"set"}: Func(func(vm *VM, aa *Args) *Args {
			aa.get(0).(*Map).Set(aa.get(1), aa.get(2))
			return join(vm, aa.get(0))
		}),
		Str{"get"}: Func(func(vm *VM, aa *Args) *Args {
			return join(vm, aa.get(0).(*Map).Get(aa.get(1)))
		}),
	})
	libMap.meta = libDef
	libList = NewMap(MapData{
		Str{"push"}: Func(func(vm *VM, aa *Args) *Args {
			l := aa.get(0).(*List)
			v := aa.get(1)
			l.data = append(l.data, v)
			return join(vm, l)
		}),
		Str{"pop"}: Func(func(vm *VM, aa *Args) *Args {
			l := aa.get(0).(*List)
			n := len(l.data) - 1
			var v Any
			if len(l.data) < n {
				v = l.data[n]
				l.data = l.data[0:n]
			}
			return join(vm, v)
		}),
		Str{"join"}: Func(func(vm *VM, aa *Args) *Args {
			l := aa.get(0).(*List)
			j := aa.get(1)
			var ls []string
			for _, s := range l.data {
				ls = append(ls, tostring(s))
			}
			return join(vm, Str{strings.Join(ls, tostring(j))})
		}),
		Str{"set"}: Func(func(vm *VM, aa *Args) *Args {
			aa.get(0).(*List).Set(aa.get(1), aa.get(2))
			return join(vm, aa.get(1), aa.get(2))
		}),
		Str{"get"}: Func(func(vm *VM, aa *Args) *Args {
			return join(vm, aa.get(0).(*List).Get(aa.get(1)))
		}),
	})
	libList.meta = libDef
	libChan = NewMap(MapData{
		Str{"read"}: Func(func(vm *VM, aa *Args) *Args {
			c := aa.get(0).(Chan)
			return join(vm, <-c)
		}),
		Str{"write"}: Func(func(vm *VM, aa *Args) *Args {
			c := aa.get(0).(Chan)
			a := aa.get(1)
			c <- a
			return join(vm, Bool(true))
		}),
		Str{"close"}: Func(func(vm *VM, aa *Args) *Args {
			c := aa.get(0).(Chan)
			close(c)
			return nil
		}),
	})
	libChan.meta = libDef
	libGroup = NewMap(MapData{
		Str{"run"}: Func(func(vm *VM, aa *Args) *Args {
			g := aa.get(0).(*Group)
			f := aa.get(1).(Func)
			ab := vm.ga(aa.len() - 2)
			for i := 0; i < aa.len()-2; i++ {
				ab.set(i, aa.get(i+2))
			}
			g.Run(f, ab)
			return join(vm, Bool(true))
		}),
		Str{"wait"}: Func(func(vm *VM, aa *Args) *Args {
			g := aa.get(0).(*Group)
			g.Wait()
			return join(vm, Bool(true))
		}),
	})
	libGroup.meta = libDef
	libStr = NewMap(MapData{
		Str{"split"}: Func(func(vm *VM, aa *Args) *Args {
			s := tostring(aa.get(0))
			j := tostring(aa.get(1))
			l := []Any{}
			for _, p := range strings.Split(s, j) {
				l = append(l, Str{p})
			}
			return join(vm, NewList(l))
		}),
	})
	libStr.meta = libDef
	libTick = NewMap(MapData{
		Str{"read"}: Func(func(vm *VM, aa *Args) *Args {
			ti := aa.get(0).(Ticker)
			return join(vm, ti.Read())
		}),
	})
	libTick.meta = libDef
	libTime = NewMap(MapData{
		Str{"ms"}: Int(int64(time.Millisecond)),
		Str{"ticker"}: Func(func(vm *VM, aa *Args) *Args {
			d := int64(aa.get(0).(Int))
			return join(vm, NewTicker(time.Duration(d)))
		}),
	})
	libTime.meta = libDef
	Ntime = libTime
}
