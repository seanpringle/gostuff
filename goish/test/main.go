package main

import "fmt"
import "math"
import "strings"
import "strconv"
import "sync"
import "time"

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

type IterIsh interface {
	Iterate() Func
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

func (i Int) Iterate() Func {
	step := int64(0)
	limit := i.i64
	return Func(func(t Tup) Tup {
		if step < limit {
			n := step
			step++
			return Tup{Int{n}}
		}
		return Tup{nil}
	})
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

type Time time.Time

func (t Time) Bool() Bool {
	var z time.Time
	return Bool{time.Time(t) != z}
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
	return Bool{t.Ticker != nil}
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
	return func() (rs bool) {
		defer func() {
			recover()
		}()
		rs = a == b
		return
	}()
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

func one(a Any) Any {
	if t, is := a.(Tup); is {
		return get(t, 0)
	}
	return a
}

func get(aa Tup, i int) Any {
	if aa != nil && i < len(aa) {
		return aa[i]
	}
	return nil
}

type funcReturn Tup

func call(f Any, aa Tup) (rs Any) {
	if fn, is := f.(Func); is {
		defer func() {
			if r := recover(); r != nil {
				if fr, is := r.(funcReturn); is {
					rs = Tup(fr)
					return
				}
				panic(r)
			}
		}()
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

func iterate(o Any) Func {
	if oi, is := o.(IterIsh); is {
		return oi.Iterate()
	}
	if oi := trymethod(o, "iterate", nil); oi != nil {
		return oi.(Func)
	}
	panic(fmt.Sprintf("not iterable: %v", o))
}

type loopBreak int

func loop(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			if _, is := r.(loopBreak); is {
				return
			}
			panic(r)
		}
	}()

	fn()
}

func block(fn func()) Tup {
	fn()
	return nil
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

var Ntime *Map

var Ntype Any = Func(func(t Tup) Tup {
	return Tup{Str{get(t, 0).Type()}}
})

func init() {
	libDef = NewMap(MapData{
		Str{"len"}: Func(func(t Tup) Tup {
			s := get(t, 0).(LenIsh)
			return Tup{Int{s.Len()}}
		}),
		Str{"lib"}: Func(func(t Tup) Tup {
			return Tup{get(t, 0).Lib()}
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
		Str{"chain"}: Func(func(t Tup) Tup {
			get(t, 0).(*Map).meta = t[1:]
			return Tup{get(t, 0)}
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
			return Tup{Bool{true}}
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
		Str{"ms"}: Int{int64(time.Millisecond)},
		Str{"ticker"}: Func(func(t Tup) Tup {
			d := get(t, 0).(Int).i64
			return Tup{NewTicker(time.Duration(d))}
		}),
	})
	libTime.meta = libDef
	Ntime = libTime
}

func main() {
	block(func() {
		var Na Any
		var Nlen Any
		var Nhi Any
		var Nlist Any
		var Nb Any
		var Nl Any
		var Ng Any
		var Nfn Any
		var Ninc Any
		var Ntrue Any
		var Nfalse Any
		var Nnil Any
		var Nt Any
		var Ns Any
		var Nc Any
		func() Tup {
			aa := join(call(Nprint, join(Int{1}, Str{"hi"})))
			Na = get(aa, 0)
			Nb = get(aa, 1)
			return aa
		}()
		call(Nprint, join(func() Any {
			var a Any
			a = func() Any {
				var a Any
				a = Int{1}
				if truth(a) {
					var b Any
					b = Int{0}
					if truth(b) {
						return b
					}
				}
				return nil
			}()
			if !truth(a) {
				a = Int{3}
			}
			return a
		}()))
		call(Nprint, join(add(one(Int{5}), one(Int{6}))))
		func() Tup {
			aa := join(Func(func(aa Tup) Tup {
				Na := get(aa, 0)
				noop(Na)
				return block(func() { ; panic(funcReturn(join(add(one(Na), one(Int{1}))))) })
			}))
			Ninc = get(aa, 0)
			return aa
		}()
		call(Nprint, join(call(Ninc, join(Int{42}))))
		call(Nprint, join(func() Any {
			var a Any
			a = func() Any {
				var a Any
				a = Bool{eq(one(Na), one(Int{1}))}
				if truth(a) {
					var b Any
					b = Int{7}
					if truth(b) {
						return b
					}
				}
				return nil
			}()
			if !truth(a) {
				a = Int{9}
			}
			return a
		}()))
		func() Tup {
			aa := join(NewMap(MapData{Str{"a"}: Int{1}, Str{"__*&^"}: Int{2}, Str{"c"}: NewMap(MapData{Str{"d"}: Func(func(aa Tup) Tup { ; return block(func() { ; panic(funcReturn(join(Str{"hello world"}))) }) })})}))
			Nt = get(aa, 0)
			return aa
		}()
		call(Nprint, join(call(find(find(Nt, Str{"c"}), Str{"d"}), join())))
		func() Tup { aa := join(Int{42}); store(Nt, Str{"a"}, get(aa, 0)); return aa }()
		call(Nprint, join(Nt))
		call(Nprint, join(Str{"\n"}, func() Any { t, m := method(Nt, Str{"keys"}); return call(m, join(t)) }()))
		call(Nprint, join(add(one(mul(one(Int{2}), one(Int{2}))), one(Int{3}))))
		func() Tup {
			aa := join(NewMap(MapData{Str{"g"}: Func(func(aa Tup) Tup { ; return block(func() { ; panic(funcReturn(join(Str{"hello world"}))) }) })}))
			Nt = get(aa, 0)
			return aa
		}()
		func() Tup {
			aa := join(Func(func(aa Tup) Tup {
				Nself := get(aa, 0)
				noop(Nself)
				return block(func() {
					panic(funcReturn(join(func() Any { t, m := method(Nself, Str{"g"}); return call(m, join(t)) }())))
				})
			}))
			store(Nt, Str{"m"}, get(aa, 0))
			return aa
		}()
		call(Nprint, join(func() Any { t, m := method(Nt, Str{"m"}); return call(m, join(t)) }()))
		func() Tup { aa := join(Str{"goodbye world"}); Ns = get(aa, 0); return aa }()
		call(Nprint, join(func() Any { t, m := method(Ns, Str{"len"}); return call(m, join(t)) }()))
		call(Nprint, join(func() Any { t, m := method(Ns, Str{"type"}); return call(m, join(t)) }()))
		call(Nprint, join(NewList([]Any{Int{1}, Int{2}, Int{7}})))
		call(Nprint, join(Str{"\n"}))
		func() Tup { aa := join(NewMap(MapData{})); Na = get(aa, 0); return aa }()
		call(Nprint, join(Na, Str{"\n"}))
		func() Any { t, m := method(Na, Str{"set"}); return call(m, join(t, Str{"1"}, Int{1})) }()
		call(Nprint, join(Na, Str{"\n"}))
		func() Tup { aa := join(NewMap(MapData{})); Nb = get(aa, 0); return aa }()
		func() Any { t, m := method(Na, Str{"set"}); return call(m, join(t, Nb, Int{2})) }()
		call(Nprint, join(Na, Str{"\n"}))
		func() Any { t, m := method(Nb, Str{"set"}); return call(m, join(t, Str{"2"}, Int{2})) }()
		call(Nprint, join(Na, Str{"\n"}))
		func() Tup { aa := join(NewList([]Any{Int{1}, Int{2}, Int{3}})); Nl = get(aa, 0); return aa }()
		call(Nprint, join(Nl, Str{"\n"}))
		func() Any { t, m := method(Nl, Str{"push"}); return call(m, join(t, Int{4})) }()
		call(Nprint, join(Nl, Str{"\n"}))
		call(Nprint, join(func() Any { t, m := method(Nl, Str{"pop"}); return call(m, join(t)) }()))
		call(Nprint, join(Nl, Str{"\n"}))
		call(Nprint, join(add(one(Str{"a"}), one(Str{"b"})), Str{"\n"}))
		func() Tup { aa := join(Str{"hi"}); Nlen = get(aa, 0); return aa }()
		call(Nprint, join(Str{"yo"}, func() Any { t, m := method(Nl, Str{"len"}); return call(m, join(t)) }(), Str{"\n"}))
		call(Nprint, join(func() Any { t, m := method(Str{"a,b,c"}, Str{"split"}); return call(m, join(t, Str{","})) }()))
		call(Nprint, join(func() Any {
			t, m := method(get(func() Any { t, m := method(Str{"a,b,c"}, Str{"split"}); return call(m, join(t, Str{","})) }().(Tup), 0), Str{"join"})
			return call(m, join(t, Str{":"}))
		}()))
		func() Tup { aa := join(call(Nchan, join(Int{10}))); Nc = get(aa, 0); return aa }()
		func() Any { t, m := method(Nc, Str{"write"}); return call(m, join(t, Int{1})) }()
		func() Any { t, m := method(Nc, Str{"write"}); return call(m, join(t, Int{2})) }()
		func() Any { t, m := method(Nc, Str{"write"}); return call(m, join(t, Int{3})) }()
		call(Nprint, join(func() Any { t, m := method(Nc, Str{"read"}); return call(m, join(t)) }()))
		call(Nprint, join(func() Any { t, m := method(Nc, Str{"read"}); return call(m, join(t)) }()))
		call(Nprint, join(func() Any { t, m := method(Nc, Str{"read"}); return call(m, join(t)) }()))
		func() Tup {
			aa := join(Func(func(aa Tup) Tup {
				Ng := get(aa, 0)
				noop(Ng)
				return block(func() { ; call(Nprint, join(Str{"hi\n"})) })
			}))
			Nhi = get(aa, 0)
			return aa
		}()
		func() Tup { aa := join(call(Ngroup, join())); Ng = get(aa, 0); return aa }()
		func() Any { t, m := method(Ng, Str{"run"}); return call(m, join(t, Nhi)) }()
		func() Any { t, m := method(Ng, Str{"run"}); return call(m, join(t, Nhi)) }()
		func() Any { t, m := method(Ng, Str{"run"}); return call(m, join(t, Nhi)) }()
		func() Any { t, m := method(Ng, Str{"wait"}); return call(m, join(t)) }()
		call(Nprint, join(Str{"done\n"}))
		call(Nprint, join(func() Any { t, m := method(Nb, Str{"get"}); return call(m, join(t, Str{"hi"})) }()))
		call(Nprint, join(func() Any {
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
		}()))
		func() Tup { aa := join(Bool{lt(one(Int{0}), one(Int{1}))}); Ntrue = get(aa, 0); return aa }()
		func() Tup { aa := join(Bool{lt(one(Int{1}), one(Int{0}))}); Nfalse = get(aa, 0); return aa }()
		func() Tup {
			aa := join(func() Any { t, m := method(NewList([]Any{}), Str{"pop"}); return call(m, join(t)) }())
			Nnil = get(aa, 0)
			return aa
		}()
		loop(func() {
			it := iterate(Int{10})
			for {
				t := it(nil)
				if get(t, 0) == nil {
					break
				}
				call(Func(func(aa Tup) Tup {
					Ni := get(aa, 0)
					noop(Ni)
					return block(func() { ; call(Nprint, join(Ni, Str{"\n"})) })
				}), t)
			}
		})
		func() Tup {
			aa := join(func() Any { t, m := method(NewList([]Any{}), Str{"lib"}); return call(m, join(t)) }())
			Nlist = get(aa, 0)
			return aa
		}()
		call(Nprint, join(Nlist, Str{"\n"}))
		func() Tup {
			aa := join(Func(func(aa Tup) Tup {
				Nself := get(aa, 0)
				noop(Nself)
				return block(func() {
					var Ni Any
					var Nl Any
					func() Tup { aa := join(Int{0}); Ni = get(aa, 0); return aa }()
					panic(funcReturn(join(Func(func(aa Tup) Tup {
						return block(func() {
							func() Tup {
								aa := join(func() Any { t, m := method(Nself, Str{"len"}); return call(m, join(t)) }())
								Nl = get(aa, 0)
								return aa
							}()
							if truth(Bool{lt(one(Ni), one(Nl))}) {
								block(func() {
									func() Tup { aa := join(add(one(Ni), one(Int{1}))); Ni = get(aa, 0); return aa }()
									call(Nprint, join(Str{"len "}, Nl, Str{" "}, Ni, Str{" "}, func() Any { t, m := method(Nself, Str{"get"}); return call(m, join(t, sub(Ni, Int{1}))) }(), Str{"\n"}))
									panic(funcReturn(join(Ni, func() Any { t, m := method(Nself, Str{"get"}); return call(m, join(t, sub(Ni, Int{1}))) }())))
								})
							}
						})
					}))))
				})
			}))
			store(Nlist, Str{"iterate"}, get(aa, 0))
			return aa
		}()
		call(Nprint, join(Nlist, Str{"\n"}))
		func() Tup {
			aa := join(func() Any {
				t, m := method(NewList([]Any{Int{1}, Int{2}, Int{3}}), Str{"iterate"})
				return call(m, join(t))
			}())
			Nfn = get(aa, 0)
			return aa
		}()
		call(Nprint, join(call(Nfn, join()), Str{"\n"}))
		call(Nprint, join(call(Nfn, join()), Str{"\n"}))
		call(Nprint, join(call(Nfn, join()), Str{"\n"}))
	})
}
