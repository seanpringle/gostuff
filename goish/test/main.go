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
			n := len(l.data) - 1
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
			c <- a
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

func main() {

	f, err := os.Create("cpuprofile")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	{
		var Nfib Any
		func() Tup {
			aa := join(Func(func(aa Tup) Tup {
				Nn := get(aa, 0)
				noop(Nn)
				{
					if truth(Bool(lt(Nn, Int(2)))) {
						{
							return Tup{Int(1)}
						}
					}
					return Tup{add(Nfib.(Func)(Tup{sub(Nn, Int(2))})[0], Nfib.(Func)(Tup{sub(Nn, Int(1))})[0])}
				}
				return nil
			}))
			Nfib = get(aa, 0)
			return aa
		}()
		Nprint.(Func)(Tup{Nfib.(Func)(Tup{Int(32)})[0]})
	}
}
