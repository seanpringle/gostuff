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
	return libDef
}

func (aa *Args) get(i int) Any {
	if i >= 32 {
		panic("maximum of 32 arguments per call")
	}
	if aa != nil && 1<<uint(i)&aa.used != 0 {
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
		aa.used = aa.used | 1<<uint(i)
	}
}

func (aa *Args) len() int {
	l := 0
	u := aa.used
	for u != 0 {
		l++
		u = u >> 1
	}
	return l
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
	return "ticker"
}

func (t Ticker) Lib() *Map {
	return libTick
}

func (t Ticker) Read() Any {
	return Time(<-t.Ticker.C)
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
	vm.da(aa)
	c := make(chan Any, int(n))
	return join(vm, Chan(c))
})

var Ngroup Any = Func(func(vm *VM, aa *Args) *Args {
	vm.da(aa)
	return join(vm, NewGroup())
})

var Ntime *Map

var Ntype Any = Func(func(vm *VM, aa *Args) *Args {
	return join(vm, Str{aa.get(0).Type()})
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
	if m, is := aa.get(0).(*Map); is {
		return join(vm, m.meta)
	}
	return join(vm, aa.get(0).Lib())
})

func init() {
	libDef = NewMap(MapData{
		Str{"len"}: Func(func(vm *VM, aa *Args) *Args {
			s := aa.get(0).(LenIsh)
			vm.da(aa)
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
			vm.da(aa)
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
			vm.da(aa)
			l.data = append(l.data, v)
			return join(vm, l)
		}),
		Str{"pop"}: Func(func(vm *VM, aa *Args) *Args {
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
		Str{"shove"}: Func(func(vm *VM, aa *Args) *Args {
			l := aa.get(0).(*List)
			v := aa.get(1)
			vm.da(aa)
			l.data = append([]Any{v}, l.data...)
			return join(vm, l)
		}),
		Str{"shift"}: Func(func(vm *VM, aa *Args) *Args {
			l := aa.get(0).(*List)
			vm.da(aa)
			var v Any
			if len(l.data) > 0 {
				v = l.data[0]
				l.data = l.data[1:]
			}
			return join(vm, v)
		}),
		Str{"join"}: Func(func(vm *VM, aa *Args) *Args {
			l := aa.get(0).(*List)
			j := aa.get(1)
			vm.da(aa)
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
			vm.da(aa)
			return join(vm, <-c)
		}),
		Str{"write"}: Func(func(vm *VM, aa *Args) *Args {
			c := aa.get(0).(Chan)
			a := aa.get(1)
			vm.da(aa)
			c <- a
			return join(vm, Bool(true))
		}),
		Str{"close"}: Func(func(vm *VM, aa *Args) *Args {
			c := aa.get(0).(Chan)
			vm.da(aa)
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
			vm.da(aa)
			return join(vm, Bool(true))
		}),
		Str{"wait"}: Func(func(vm *VM, aa *Args) *Args {
			g := aa.get(0).(*Group)
			vm.da(aa)
			g.Wait()
			return join(vm, Bool(true))
		}),
	})
	libGroup.meta = libDef
	libStr = NewMap(MapData{
		Str{"split"}: Func(func(vm *VM, aa *Args) *Args {
			s := tostring(aa.get(0))
			j := tostring(aa.get(1))
			vm.da(aa)
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
			vm.da(aa)
			return join(vm, ti.Read())
		}),
		Str{"stop"}: Func(func(vm *VM, aa *Args) *Args {
			ti := aa.get(0).(Ticker)
			vm.da(aa)
			return join(vm, ti.Stop())
		}),
	})
	libTick.meta = libDef
	libTime = NewMap(MapData{
		Str{"ms"}: Int(int64(time.Millisecond)),
		Str{"ticker"}: Func(func(vm *VM, aa *Args) *Args {
			d := int64(aa.get(0).(Int))
			vm.da(aa)
			return join(vm, NewTicker(time.Duration(d)))
		}),
	})
	libTime.meta = libDef
	Ntime = libTime
}

func main() {

	f, err := os.Create("cpuprofile")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	vm := &VM{}

	{
		var Nc Any
		var Nnil Any
		var Na Any
		var Nlen Any
		var Ntrue Any
		var Nfalse Any
		var Nt Any
		var Nblinker Any
		var Nblink Any
		var Nb Any
		var Ninc Any
		var Nl Any
		var Nm Any
		var Ns Any
		var Nhi Any
		var Ng Any
		func() Any { a := Bool(lt(Int(0), Int(1))); Ntrue = a; return a }()
		func() Any { a := Bool(lt(Int(1), Int(0))); Nfalse = a; return a }()
		func() Any {
			a := one(vm, func() *Args { t, m := method(NewList([]Any{}), Str{"pop"}); return call(vm, m, join(vm, t, nil)) }())
			Nnil = a
			return a
		}()
		vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
			vm.da(aa)
			{
				var Np Any
				func() Any { a := one(vm, call(vm, Ngetprototype, join(vm, Int(0)))); Np = a; return a }()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nself := aa.get(0)
						noop(Nself)
						vm.da(aa)
						{
							var Ni Any
							func() Any { a := Int(0); Ni = a; return a }()
							return join(vm, Func(func(vm *VM, aa *Args) *Args {
								vm.da(aa)
								{
									if lt(Ni, Nself) {
										{
											return join(vm, func() Any { v := one(vm, Ni); Ni = add(v, Int(1)); return v }())
										}
									}
								}
								return nil
							}))
						}
						return nil
					}))
					store(Np, Str{"iterate"}, a)
					return a
				}()
			}
			return nil
		}), join(vm, nil)))
		vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
			vm.da(aa)
			{
				var Nl Any
				func() Any { a := one(vm, call(vm, Ngetprototype, join(vm, NewList([]Any{})))); Nl = a; return a }()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nself := aa.get(0)
						noop(Nself)
						vm.da(aa)
						{
							var Ni Any
							func() Any { a := Int(0); Ni = a; return a }()
							return join(vm, Func(func(vm *VM, aa *Args) *Args {
								vm.da(aa)
								{
									if lt(Ni, one(vm, func() *Args { t, m := method(Nself, Str{"len"}); return call(vm, m, join(vm, t, nil)) }())) {
										{
											return join(vm, func() Any { v := one(vm, Ni); Ni = add(v, Int(1)); return v }(), func() *Args { t, m := method(Nself, Str{"get"}); return call(vm, m, join(vm, t, sub(Ni, Int(1)))) }())
										}
									}
								}
								return nil
							}))
						}
						return nil
					}))
					store(Nl, Str{"iterate"}, a)
					return a
				}()
			}
			return nil
		}), join(vm, nil)))
		vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
			vm.da(aa)
			{
				var Nm Any
				func() Any { a := one(vm, call(vm, Ngetprototype, join(vm, NewMap(MapData{})))); Nm = a; return a }()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nself := aa.get(0)
						noop(Nself)
						vm.da(aa)
						{
							var Ni Any
							var Nkeys Any
							var Nkey Any
							func() Any { a := Int(0); Ni = a; return a }()
							func() Any {
								a := one(vm, func() *Args { t, m := method(Nself, Str{"keys"}); return call(vm, m, join(vm, t, nil)) }())
								Nkeys = a
								return a
							}()
							return join(vm, Func(func(vm *VM, aa *Args) *Args {
								vm.da(aa)
								{
									if lt(Ni, one(vm, func() *Args { t, m := method(Nkeys, Str{"len"}); return call(vm, m, join(vm, t, nil)) }())) {
										{
											func() Any {
												a := one(vm, func() *Args {
													t, m := method(Nkeys, Str{"get"})
													return call(vm, m, join(vm, t, func() Any { v := one(vm, Ni); Ni = add(v, Int(1)); return v }()))
												}())
												Nkey = a
												return a
											}()
											return join(vm, Nkey, func() *Args { t, m := method(Nself, Str{"get"}); return call(vm, m, join(vm, t, Nkey)) }())
										}
									}
								}
								return nil
							}))
						}
						return nil
					}))
					store(Nm, Str{"iterate"}, a)
					return a
				}()
			}
			return nil
		}), join(vm, nil)))
		vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
			vm.da(aa)
			{
				var Nti Any
				var Ntick Any
				func() Any {
					a := one(vm, call(vm, find(Ntime, Str{"ticker"}), join(vm, Int(1000000))))
					Nti = a
					return a
				}()
				func() Any { a := one(vm, call(vm, Ngetprototype, join(vm, Nti))); Ntick = a; return a }()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nself := aa.get(0)
						noop(Nself)
						vm.da(aa)
						{
							var Ni Any
							func() Any { a := Int(0); Ni = a; return a }()
							return join(vm, Func(func(vm *VM, aa *Args) *Args {
								vm.da(aa)
								{
									return join(vm, func() Any { v := one(vm, Ni); Ni = add(v, Int(1)); return v }(), func() *Args { t, m := method(Nself, Str{"read"}); return call(vm, m, join(vm, t, nil)) }())
								}
								return nil
							}))
						}
						return nil
					}))
					store(Ntick, Str{"iterate"}, a)
					return a
				}()
				vm.da(func() *Args { t, m := method(Nti, Str{"stop"}); return call(vm, m, join(vm, t, nil)) }())
			}
			return nil
		}), join(vm, nil)))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			aa := join(vm, call(vm, Nprint, join(vm, Int(1), Str{"hi"})))
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
			a := one(vm, NewMap(MapData{Str{"a"}: Int(1), Str{"__*&^"}: Int(2), Str{"c"}: NewMap(MapData{Str{"d"}: Func(func(vm *VM, aa *Args) *Args {
				vm.da(aa)
				{
					return join(vm, Str{"hello world"})
				}
				return nil
			})})}))
			Nt = a
			return a
		}()
		vm.da(call(vm, Nprint, join(vm, call(vm, find(find(Nt, Str{"c"}), Str{"d"}), join(vm, nil)))))
		func() Any { a := Int(42); store(Nt, Str{"a"}, a); return a }()
		vm.da(call(vm, Nprint, join(vm, Nt)))
		vm.da(call(vm, Nprint, join(vm, Str{""}, func() *Args { t, m := method(Nt, Str{"keys"}); return call(vm, m, join(vm, t, nil)) }())))
		vm.da(call(vm, Nprint, join(vm, add(one(vm, mul(Int(2), Int(2))), Int(3)))))
		func() Any {
			a := one(vm, NewMap(MapData{Str{"g"}: Func(func(vm *VM, aa *Args) *Args {
				vm.da(aa)
				{
					return join(vm, Str{"hello world"})
				}
				return nil
			})}))
			Nt = a
			return a
		}()
		func() Any {
			a := one(vm, Func(func(vm *VM, aa *Args) *Args {
				Nself := aa.get(0)
				noop(Nself)
				vm.da(aa)
				{
					return join(vm, func() *Args { t, m := method(Nself, Str{"g"}); return call(vm, m, join(vm, t, nil)) }())
				}
				return nil
			}))
			store(Nt, Str{"m"}, a)
			return a
		}()
		vm.da(call(vm, Nprint, join(vm, func() *Args { t, m := method(Nt, Str{"m"}); return call(vm, m, join(vm, t, nil)) }())))
		func() Any { a := one(vm, Str{"goodbye world"}); Ns = a; return a }()
		vm.da(call(vm, Nprint, join(vm, func() *Args { t, m := method(Ns, Str{"len"}); return call(vm, m, join(vm, t, nil)) }())))
		vm.da(call(vm, Nprint, join(vm, call(vm, Ntype, join(vm, Ns)))))
		vm.da(call(vm, Nprint, join(vm, NewList([]Any{Int(1), Int(2), Int(7)}))))
		func() Any { a := one(vm, NewMap(MapData{})); Na = a; return a }()
		vm.da(call(vm, Nprint, join(vm, Na)))
		vm.da(func() *Args { t, m := method(Na, Str{"set"}); return call(vm, m, join(vm, t, Str{"1"}, Int(1))) }())
		vm.da(call(vm, Nprint, join(vm, Na)))
		func() Any { a := one(vm, NewMap(MapData{})); Nb = a; return a }()
		vm.da(func() *Args { t, m := method(Na, Str{"set"}); return call(vm, m, join(vm, t, Nb, Int(2))) }())
		vm.da(call(vm, Nprint, join(vm, Na)))
		vm.da(func() *Args { t, m := method(Nb, Str{"set"}); return call(vm, m, join(vm, t, Str{"2"}, Int(2))) }())
		vm.da(call(vm, Nprint, join(vm, Na)))
		func() Any { a := one(vm, NewList([]Any{Int(1), Int(2), Int(3)})); Nl = a; return a }()
		vm.da(call(vm, Nprint, join(vm, Nl)))
		vm.da(func() *Args { t, m := method(Nl, Str{"push"}); return call(vm, m, join(vm, t, Int(4))) }())
		vm.da(call(vm, Nprint, join(vm, Nl)))
		vm.da(call(vm, Nprint, join(vm, func() *Args { t, m := method(Nl, Str{"pop"}); return call(vm, m, join(vm, t, nil)) }())))
		vm.da(call(vm, Nprint, join(vm, Nl)))
		vm.da(call(vm, Nprint, join(vm, add(one(vm, Str{"a"}), one(vm, Str{"b"})))))
		func() Any { a := one(vm, Str{"hi"}); Nlen = a; return a }()
		vm.da(call(vm, Nprint, join(vm, Str{"yo"}, func() *Args { t, m := method(Nl, Str{"len"}); return call(vm, m, join(vm, t, nil)) }())))
		vm.da(call(vm, Nprint, join(vm, func() *Args { t, m := method(Str{"a,b,c"}, Str{"split"}); return call(vm, m, join(vm, t, Str{","})) }())))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(one(vm, func() *Args { t, m := method(Str{"a,b,c"}, Str{"split"}); return call(vm, m, join(vm, t, Str{","})) }()), Str{"join"})
			return call(vm, m, join(vm, t, Str{":"}))
		}())))
		func() Any { a := one(vm, call(vm, Nchan, join(vm, Int(10)))); Nc = a; return a }()
		vm.da(func() *Args { t, m := method(Nc, Str{"write"}); return call(vm, m, join(vm, t, Int(1))) }())
		vm.da(func() *Args { t, m := method(Nc, Str{"write"}); return call(vm, m, join(vm, t, Int(2))) }())
		vm.da(func() *Args { t, m := method(Nc, Str{"write"}); return call(vm, m, join(vm, t, Int(3))) }())
		vm.da(call(vm, Nprint, join(vm, func() *Args { t, m := method(Nc, Str{"read"}); return call(vm, m, join(vm, t, nil)) }())))
		vm.da(call(vm, Nprint, join(vm, func() *Args { t, m := method(Nc, Str{"read"}); return call(vm, m, join(vm, t, nil)) }())))
		vm.da(call(vm, Nprint, join(vm, func() *Args { t, m := method(Nc, Str{"read"}); return call(vm, m, join(vm, t, nil)) }())))
		func() Any {
			a := one(vm, Func(func(vm *VM, aa *Args) *Args {
				Ng := aa.get(0)
				noop(Ng)
				vm.da(aa)
				{
					vm.da(call(vm, Nprint, join(vm, Str{"hi"})))
				}
				return nil
			}))
			Nhi = a
			return a
		}()
		func() Any { a := one(vm, call(vm, Ngroup, join(vm, nil))); Ng = a; return a }()
		vm.da(func() *Args { t, m := method(Ng, Str{"run"}); return call(vm, m, join(vm, t, Nhi)) }())
		vm.da(func() *Args { t, m := method(Ng, Str{"run"}); return call(vm, m, join(vm, t, Nhi)) }())
		vm.da(func() *Args { t, m := method(Ng, Str{"run"}); return call(vm, m, join(vm, t, Nhi)) }())
		vm.da(func() *Args { t, m := method(Ng, Str{"wait"}); return call(vm, m, join(vm, t, nil)) }())
		vm.da(call(vm, Nprint, join(vm, Str{"done"})))
		vm.da(call(vm, Nprint, join(vm, func() *Args { t, m := method(Nb, Str{"get"}); return call(vm, m, join(vm, t, Str{"hi"})) }())))
		vm.da(call(vm, Nprint, join(vm, func() Any {
			var a Any
			a = func() Any {
				var a Any
				a = Ntrue
				if truth(a) {
					var b Any
					b = Str{"yes"}
					if truth(b) {
						return b
					}
				}
				return nil
			}()
			if !truth(a) {
				a = Str{"no"}
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
						vm.da(call(vm, Nprint, join(vm, Ni, Str{":"}, Nv)))
					}
					return nil
				}), aa))
			}
		})
		loop(func() {
			it := iterate(one(vm, NewMap(MapData{Str{"tom"}: Int(1), Str{"dick"}: Int(2), Str{"harry"}: Int(43)})))
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
						vm.da(call(vm, Nprint, join(vm, Nk, Str{"=>"}, Nv)))
					}
					return nil
				}), aa))
			}
		})
		func() Any { a := Int(1); Na = a; return a }()
		vm.da(call(vm, Nprint, join(vm, func() Any { v := one(vm, Na); Na = add(v, Int(1)); return v }())))
		vm.da(call(vm, Nprint, join(vm, func() Any { v := one(vm, Na); Na = add(v, Int(1)); return v }())))
		vm.da(call(vm, Nprint, join(vm, func() Any { v := one(vm, Na); Na = add(v, Int(1)); return v }())))
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
		func() Any {
			a := one(vm, Func(func(vm *VM, aa *Args) *Args {
				vm.da(aa)
				{
					var Nlock Any
					var Njobs Any
					func() Any { a := one(vm, call(vm, Nchan, join(vm, Int(1)))); Nlock = a; return a }()
					func() Any { a := one(vm, NewList([]Any{})); Njobs = a; return a }()
					return join(vm, Func(func(vm *VM, aa *Args) *Args {
						Nfn := aa.get(0)
						noop(Nfn)
						vm.da(aa)
						{
							vm.da(func() *Args { t, m := method(Nlock, Str{"write"}); return call(vm, m, join(vm, t, nil)) }())
							if eq(Nfn, Nnil) {
								{
									loop(func() {
										it := iterate(Njobs)
										for {
											aa := it(vm, nil)
											if aa.get(0) == nil {
												vm.da(aa)
												break
											}
											vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
												Ni := aa.get(0)
												noop(Ni)
												Njob := aa.get(1)
												noop(Njob)
												vm.da(aa)
												{
													vm.da(call(vm, Njob, join(vm, nil)))
												}
												return nil
											}), aa))
										}
									})
									func() Any { a := one(vm, NewList([]Any{})); Njobs = a; return a }()
								}
							} else {
								{
									vm.da(func() *Args { t, m := method(Njobs, Str{"push"}); return call(vm, m, join(vm, t, Nfn)) }())
								}
							}
							vm.da(func() *Args { t, m := method(Nlock, Str{"read"}); return call(vm, m, join(vm, t, nil)) }())
						}
						return nil
					}))
				}
				return nil
			}))
			Nblinker = a
			return a
		}()
		func() Any { a := one(vm, call(vm, Nblinker, join(vm, nil))); Nblink = a; return a }()
		vm.da(call(vm, Nblink, join(vm, Func(func(vm *VM, aa *Args) *Args {
			vm.da(aa)
			{
				vm.da(call(vm, Nprint, join(vm, Str{"hello world"})))
			}
			return nil
		}))))
		vm.da(call(vm, Nblink, join(vm, Func(func(vm *VM, aa *Args) *Args {
			vm.da(aa)
			{
				vm.da(call(vm, Nprint, join(vm, Str{"hello world"})))
			}
			return nil
		}))))
		vm.da(call(vm, Nblink, join(vm, Func(func(vm *VM, aa *Args) *Args {
			vm.da(aa)
			{
				vm.da(call(vm, Nprint, join(vm, Str{"hello world"})))
			}
			return nil
		}))))
		vm.da(call(vm, Nprint, join(vm, Str{"and..."})))
		vm.da(call(vm, Nblink, join(vm, nil)))
		func() Any { a := one(vm, NewList([]Any{Int(1), Int(2), Int(3)})); Nl = a; return a }()
		vm.da(call(vm, Nprint, join(vm, field(Nl, Int(0)))))
		func() Any {
			a := one(vm, NewMap(MapData{Str{"a"}: Int(1), Str{"b"}: NewMap(MapData{Str{"c"}: Int(4)})}))
			Nm = a
			return a
		}()
		vm.da(call(vm, Nprint, join(vm, field(field(Nm, Str{"b"}), Str{"c"}))))
		func() Any { a := Int(5); store(field(Nm, Str{"b"}), Str{"c"}, a); return a }()
		vm.da(call(vm, Nprint, join(vm, field(field(Nm, Str{"b"}), Str{"c"}))))
	}
}
