package main

import "fmt"
import "math"
import "strings"
import "strconv"
import "sync"

type Any interface {
	Type() string
	Lib() *Map
	String() string
}
type Tup []Any

var libDef *Map
var libMap *Map
var libList *Map
var libStr *Map
var libChan *Map
var libGroup *Map

type Stringer = fmt.Stringer

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
	if v, ok := t.data[key]; ok {
		return v
	}
	return nil
}

func (t *Map) Set(key Any, val Any) {
	if val == nil {
		delete(t.data, key)
	} else {
		t.data[key] = val
	}
}

func (t *Map) Find(key Any) Any {
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
	return nil
}

func (t *Map) Type() string {
	return "map"
}

func (t *Map) Len() int64 {
	return int64(len(t.data))
}

func (t *Map) String() string {
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
	if a == nil {
		return false
	}
	if b, is := a.(bool); is {
		return b
	}
	if ab, is := a.(BoolIsh); is {
		return ab.Bool().b
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

var n_print Any = Func(func(aa Tup) Tup {
	for _, a := range aa {
		fmt.Printf("%s", tostring(a))
	}
	return aa
})

var n_chan Any = Func(func(t Tup) Tup {
	n := get(t, 0).(IntIsh).Int().i64
	c := make(chan Any, int(n))
	return Tup{Chan(c)}
})

var n_group Any = Func(func(t Tup) Tup {
	return Tup{NewGroup()}
})

var n_len Any = Func(func(t Tup) Tup {
	s := get(t, 0).(LenIsh)
	return Tup{Int{s.Len()}}
})

var n_type Any = Func(func(t Tup) Tup {
	return Tup{Str{get(t, 0).Type()}}
})

var n_string Any = Func(func(t Tup) Tup {
	return Tup{Str{fmt.Sprintf("%v", get(t, 0))}}
})

var n_lib Any = Func(func(t Tup) Tup {
	return Tup{get(t, 0).Lib()}
})

var n_set Any = Func(func(t Tup) Tup {
	get(t, 0).(*Map).Set(get(t, 1), get(t, 2))
	return Tup{get(t, 1), get(t, 2)}
})

var n_get Any = Func(func(t Tup) Tup {
	return Tup{get(t, 0).(*Map).Get(get(t, 1))}
})

func init() {
	libDef = NewMap(MapData{
		Str{"len"}:    n_len,
		Str{"lib"}:    n_lib,
		Str{"type"}:   n_type,
		Str{"string"}: n_string,
	})
	libMap = NewMap(MapData{
		Str{"keys"}: Func(func(t Tup) Tup {
			keys := []Any{}
			for k, _ := range get(t, 0).(*Map).data {
				keys = append(keys, k)
			}
			return Tup{NewList(keys)}
		}),
		Str{"set"}: n_set,
		Str{"get"}: n_get,
		Str{"setmeta"}: Func(func(t Tup) Tup {
			get(t, 0).(*Map).meta = get(t, 1)
			return Tup{get(t, 1)}
		}),
		Str{"getmeta"}: Func(func(t Tup) Tup {
			return Tup{get(t, 0).(*Map).meta}
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
			v := l.data[n]
			l.data = l.data[0:n]
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
}

func main() {
	func() Tup {
		var n_b Any
		var n_t Any
		var n_s Any
		var n_len Any
		var n_hi Any
		var n_g Any
		var n_a Any
		var n_inc Any
		var n_i Any
		var n_l Any
		var n_c Any
		func() Tup {
			aa := join(call(n_print, join(Int{1}, Str{"hi"})))
			n_a = get(aa, 0)
			n_b = get(aa, 1)
			return aa
		}()
		call(n_print, join(func() Any {
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
		call(n_print, join(add(Int{5}, Int{6})))
		func() Tup {
			aa := join(Func(func(aa Tup) Tup {
				n_a := get(aa, 0)
				return func() Tup {

					return join(add(n_a, Int{1}))
					return nil
				}()
			}))
			n_inc = get(aa, 0)
			return aa
		}()
		call(n_print, join(call(n_inc, join(Int{42}))))
		call(n_print, join(func() Any {
			var a Any
			a = func() Any {
				var a Any
				a = Bool{eq(n_a, Int{1})}
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
		func() Tup { aa := join(Int{0}); n_i = get(aa, 0); return aa }()
		for truth(Bool{lt(n_i, Int{10})}) {
			func() Tup {

				call(n_print, join(n_i))
				func() Tup { aa := join(add(n_i, Int{1})); n_i = get(aa, 0); return aa }()
				return nil
			}()
		}
		func() Tup {
			aa := join(NewMap(MapData{Str{"a"}: Int{1}, Str{"__*&^"}: Int{2}, Str{"c"}: NewMap(MapData{Str{"d"}: Func(func(aa Tup) Tup {
				return func() Tup {

					return join(Str{"helloworld"})
					return nil
				}()
			})})}))
			n_t = get(aa, 0)
			return aa
		}()
		call(n_print, join(call(find(find(n_t, Str{"c"}), Str{"d"}), join())))
		func() Tup { aa := join(Int{42}); store(n_t, Str{"a"}, get(aa, 0)); return aa }()
		call(n_print, join(n_t))
		call(n_print, join(Str{"\n"}, func() Any { t, m := method(n_t, Str{"keys"}); return call(m, join(t)) }()))
		call(n_print, join(add(mul(Int{2}, Int{2}), Int{3})))
		func() Tup {
			aa := join(NewMap(MapData{Str{"g"}: Func(func(aa Tup) Tup {
				return func() Tup {

					return join(Str{"helloworld"})
					return nil
				}()
			})}))
			n_t = get(aa, 0)
			return aa
		}()
		func() Tup {
			aa := join(Func(func(aa Tup) Tup {
				n_self := get(aa, 0)
				return func() Tup {

					return join(func() Any { t, m := method(n_self, Str{"g"}); return call(m, join(t)) }())
					return nil
				}()
			}))
			store(n_t, Str{"m"}, get(aa, 0))
			return aa
		}()
		call(n_print, join(func() Any { t, m := method(n_t, Str{"m"}); return call(m, join(t)) }()))
		func() Tup { aa := join(Str{"goodbyeworld"}); n_s = get(aa, 0); return aa }()
		call(n_print, join(func() Any { t, m := method(n_s, Str{"len"}); return call(m, join(t)) }()))
		call(n_print, join(func() Any { t, m := method(n_s, Str{"type"}); return call(m, join(t)) }()))
		call(n_print, join(NewList([]Any{Int{1}, Int{2}, Int{7}})))
		call(n_print, join(Str{"\n"}))
		func() Tup { aa := join(NewMap(MapData{})); n_a = get(aa, 0); return aa }()
		call(n_print, join(n_a, Str{"\n"}))
		func() Any { t, m := method(n_a, Str{"set"}); return call(m, join(t, Str{"1"}, Int{1})) }()
		call(n_print, join(n_a, Str{"\n"}))
		func() Tup { aa := join(NewMap(MapData{})); n_b = get(aa, 0); return aa }()
		func() Any { t, m := method(n_a, Str{"set"}); return call(m, join(t, n_b, Int{2})) }()
		call(n_print, join(n_a, Str{"\n"}))
		func() Any { t, m := method(n_b, Str{"set"}); return call(m, join(t, Str{"2"}, Int{2})) }()
		call(n_print, join(n_a, Str{"\n"}))
		call(n_print, join(func() Any { t, m := method(n_a, Str{"getmeta"}); return call(m, join(t)) }(), Str{"\n"}))
		func() Tup { aa := join(NewList([]Any{Int{1}, Int{2}, Int{3}})); n_l = get(aa, 0); return aa }()
		call(n_print, join(n_l, Str{"\n"}))
		func() Any { t, m := method(n_l, Str{"push"}); return call(m, join(t, Int{4})) }()
		call(n_print, join(n_l, Str{"\n"}))
		call(n_print, join(func() Any { t, m := method(n_l, Str{"pop"}); return call(m, join(t)) }()))
		call(n_print, join(n_l, Str{"\n"}))
		call(n_print, join(add(Str{"a"}, Str{"b"}), Str{"\n"}))
		func() Tup { aa := join(Str{"hi"}); n_len = get(aa, 0); return aa }()
		call(n_print, join(Str{"yo"}, func() Any { t, m := method(n_l, Str{"len"}); return call(m, join(t)) }(), Str{"\n"}))
		call(n_print, join(func() Any { t, m := method(Str{"a,b,c"}, Str{"split"}); return call(m, join(t, Str{","})) }()))
		call(n_print, join(func() Any {
			t, m := method(get(func() Any { t, m := method(Str{"a,b,c"}, Str{"split"}); return call(m, join(t, Str{","})) }().(Tup), 0), Str{"join"})
			return call(m, join(t, Str{":"}))
		}()))
		func() Tup { aa := join(call(n_chan, join(Int{10}))); n_c = get(aa, 0); return aa }()
		func() Any { t, m := method(n_c, Str{"write"}); return call(m, join(t, Int{1})) }()
		func() Any { t, m := method(n_c, Str{"write"}); return call(m, join(t, Int{2})) }()
		func() Any { t, m := method(n_c, Str{"write"}); return call(m, join(t, Int{3})) }()
		call(n_print, join(func() Any { t, m := method(n_c, Str{"read"}); return call(m, join(t)) }()))
		call(n_print, join(func() Any { t, m := method(n_c, Str{"read"}); return call(m, join(t)) }()))
		call(n_print, join(func() Any { t, m := method(n_c, Str{"read"}); return call(m, join(t)) }()))
		func() Tup {
			aa := join(Func(func(aa Tup) Tup {
				return func() Tup {

					call(n_print, join(Str{"hi\n"}))
					return nil
				}()
			}))
			n_hi = get(aa, 0)
			return aa
		}()
		func() Tup { aa := join(call(n_group, join())); n_g = get(aa, 0); return aa }()
		func() Any { t, m := method(n_g, Str{"run"}); return call(m, join(t, n_hi)) }()
		func() Any { t, m := method(n_g, Str{"run"}); return call(m, join(t, n_hi)) }()
		func() Any { t, m := method(n_g, Str{"run"}); return call(m, join(t, n_hi)) }()
		func() Any { t, m := method(n_g, Str{"wait"}); return call(m, join(t)) }()
		call(n_print, join(Str{"done\n"}))
		call(n_print, join(func() Any { t, m := method(n_b, Str{"get"}); return call(m, join(t, Str{"hi"})) }()))
		return nil
	}()
}
