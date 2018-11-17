package main

import "log"
import "runtime/pprof"
import "os"

const S15 Text = Text("json")
const S21 Text = Text("shift")
const S25 Text = Text("get")
const S2 Text = Text("remove")
const S4 Text = Text("stdout")
const S5 Text = Text("write")
const S7 Text = Text("flush")
const S8 Text = Text("stderr")
const S36 Text = Text("date")
const S37 Text = Text("readrune")
const S41 Text = Text("close")
const S33 Text = Text("queue")
const S1 Text = Text("stdin")
const S10 Text = Text("len")
const S11 Text = Text("iterate")
const S16 Text = Text("quote")
const S24 Text = Text("set")
const S32 Text = Text("channel")
const S39 Text = Text("open")
const S43 Text = Text("query")
const S6 Text = Text("join")
const S13 Text = Text("min")
const S14 Text = Text("text")
const S27 Text = Text("ticker")
const S29 Text = Text("read")
const S44 Text = Text("begin")
const S19 Text = Text("pop")
const S34 Text = Text("now")
const S40 Text = Text("readall")
const S42 Text = Text("slurp")
const S12 Text = Text("max")
const S18 Text = Text("push")
const S30 Text = Text("lock")
const S35 Text = Text("split")
const S17 Text = Text("insert")
const S20 Text = Text("shove")
const S22 Text = Text("clear")
const S28 Text = Text("stop")
const S31 Text = Text("jobs")
const S45 Text = Text("rollback")
const S3 Text = Text("string")
const S9 Text = Text("type")
const S23 Text = Text("extend")
const S26 Text = Text("keys")
const S38 Text = Text("readline")

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
		var Nmap Any
		noop(Nmap)
		var Nprint Any
		noop(Nprint)
		var Nlog Any
		noop(Nlog)
		var Nsuper Any
		noop(Nsuper)
		var Nfalse Any
		noop(Nfalse)
		var Ngood Any
		noop(Ngood)
		var Nis Any
		noop(Nis)
		var Nok Any
		noop(Nok)
		var Ndb Any
		noop(Ndb)
		var Ndecimal Any
		noop(Ndecimal)
		var Nstream Any
		noop(Nstream)
		var Nnil Any
		noop(Nnil)
		var Nres Any
		noop(Nres)
		var Ntx Any
		noop(Ntx)
		var Ninteger Any
		noop(Ninteger)
		var Nstring Any
		noop(Nstring)
		var Nlist Any
		noop(Nlist)
		var Ntrue Any
		noop(Ntrue)
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
				t, m := method(one(vm, NewList([]Any{})), S2 /* remove */)
				return call(vm, m, join(vm, t, Int(0)))
			}())
			Nnil = a
			return a
		}()
		func() Any { a := one(vm, call(vm, Nstatus, join(vm, Nnil))); Ngood = a; return a }()
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
						t, m := method(one(vm, find(Nio, S4 /* stdout */)), S5 /* write */)
						return call(vm, m, join(vm, t, concat(one(vm, func() *Args {
							t, m := method(Na, S6 /* join */)
							return call(vm, m, join(vm, t, Text("\t")))
						}()), Text("\n"))))
					}())
					vm.da(func() *Args {
						t, m := method(one(vm, find(Nio, S4 /* stdout */)), S7 /* flush */)
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
						t, m := method(one(vm, find(Nio, S8 /* stderr */)), S5 /* write */)
						return call(vm, m, join(vm, t, concat(one(vm, func() *Args {
							t, m := method(Na, S6 /* join */)
							return call(vm, m, join(vm, t, Text("\t")))
						}()), Text("\n"))))
					}())
					vm.da(func() *Args {
						t, m := method(one(vm, find(Nio, S8 /* stderr */)), S7 /* flush */)
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
					if eq(Nclass, Nproto) {
						{
							return join(vm, Ntrue)
						}
					}
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
					return join(vm, Nfalse)
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
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nn := aa.get(0)
						noop(Nn)
						vm.da(aa)
						{
							return join(vm, func() *Args {
								t, m := method(Nn, S14 /* text */)
								return call(vm, m, join(vm, t, nil))
							}())
						}
						return nil
					}))
					store(Np, S15 /* json */, a)
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
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nn := aa.get(0)
						noop(Nn)
						vm.da(aa)
						{
							return join(vm, func() *Args {
								t, m := method(Nn, S14 /* text */)
								return call(vm, m, join(vm, t, nil))
							}())
						}
						return nil
					}))
					store(Np, S15 /* json */, a)
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
				func() Any { a := Nstring; Np = a; return a }()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nn := aa.get(0)
						noop(Nn)
						vm.da(aa)
						{
							return join(vm, func() *Args {
								t, m := method(Nn, S16 /* quote */)
								return call(vm, m, join(vm, t, nil))
							}())
						}
						return nil
					}))
					store(Np, S15 /* json */, a)
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
								t, m := method(Nlist, S17 /* insert */)
								return call(vm, m, join(vm, t, length(Nlist), Nval))
							}())
						}
						return nil
					}))
					store(Np, S18 /* push */, a)
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
					store(Np, S19 /* pop */, a)
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
								t, m := method(Nlist, S17 /* insert */)
								return call(vm, m, join(vm, t, Int(0), Nval))
							}())
						}
						return nil
					}))
					store(Np, S20 /* shove */, a)
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
					store(Np, S21 /* shift */, a)
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
												t, m := method(Nlist, S19 /* pop */)
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
					store(Np, S22 /* clear */, a)
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nlist := aa.get(0)
						noop(Nlist)
						vm.da(aa)
						{
							var Nn Any
							noop(Nn)
							var Ni Any
							noop(Ni)
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
												t, m := method(Nlist, S18 /* push */)
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
					store(Np, S23 /* extend */, a)
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
					store(Np, S24 /* set */, a)
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
					store(Np, S25 /* get */, a)
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
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nl := aa.get(0)
						noop(Nl)
						vm.da(aa)
						{
							var Nparts Any
							noop(Nparts)
							func() Any { a := one(vm, NewList([]Any{})); Nparts = a; return a }()
							loop(func() {
								it := iterate(Nl)
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
											vm.da(func() *Args {
												t, m := method(Nparts, S18 /* push */)
												return call(vm, m, join(vm, t, func() *Args {
													t, m := method(Nitem, S15 /* json */)
													return call(vm, m, join(vm, t, nil))
												}()))
											}())
										}
										return nil
									}), aa))
								}
							})
							return join(vm, func() *Args {
								t, m := method(one(vm, NewList([]Any{Text("["), one(vm, func() *Args {
									t, m := method(Nparts, S6 /* join */)
									return call(vm, m, join(vm, t, Text(",")))
								}()), Text("]")})), S6 /* join */)
								return call(vm, m, join(vm, t, nil))
							}())
						}
						return nil
					}))
					store(Np, S15 /* json */, a)
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
							var Ni Any
							noop(Ni)
							var Nkeys Any
							noop(Nkeys)
							var Nn Any
							noop(Nn)
							var Nkey Any
							noop(Nkey)
							func() Any { a := Int(0); Ni = a; return a }()
							func() Any {
								a := one(vm, func() *Args {
									t, m := method(Nm, S26 /* keys */)
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
					store(Np, S24 /* set */, a)
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
					store(Np, S25 /* get */, a)
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
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nm := aa.get(0)
						noop(Nm)
						vm.da(aa)
						{
							var Nparts Any
							noop(Nparts)
							var Nkq Any
							noop(Nkq)
							func() Any { a := one(vm, NewList([]Any{})); Nparts = a; return a }()
							loop(func() {
								it := iterate(Nm)
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
											func() Any {
												a := one(vm, func() *Args {
													t, m := method(one(vm, func() *Args {
														t, m := method(Nk, S14 /* text */)
														return call(vm, m, join(vm, t, nil))
													}()), S16 /* quote */)
													return call(vm, m, join(vm, t, nil))
												}())
												Nkq = a
												return a
											}()
											if eq(Nv, Nnil) {
												{
													vm.da(func() *Args {
														t, m := method(Nparts, S18 /* push */)
														return call(vm, m, join(vm, t, concat(Nkq, Text(": null"))))
													}())
													return join(vm, nil)
												}
											}
											if !eq(one(vm, find(Nv, S15 /* json */)), Nnil) {
												{
													vm.da(func() *Args {
														t, m := method(Nparts, S18 /* push */)
														return call(vm, m, join(vm, t, concat(concat(Nkq, Text(":")), one(vm, func() *Args {
															t, m := method(Nv, S15 /* json */)
															return call(vm, m, join(vm, t, nil))
														}()))))
													}())
													return join(vm, nil)
												}
											}
										}
										return nil
									}), aa))
								}
							})
							return join(vm, concat(concat(Text("{"), one(vm, func() *Args {
								t, m := method(Nparts, S6 /* join */)
								return call(vm, m, join(vm, t, Text(",")))
							}())), Text("}")))
						}
						return nil
					}))
					store(Np, S15 /* json */, a)
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
					a := one(vm, call(vm, find(Ntime, S27 /* ticker */), join(vm, Int(1000000))))
					Nti = a
					return a
				}()
				func() Any { a := one(vm, call(vm, Ngetprototype, join(vm, Nti))); Ntick = a; return a }()
				vm.da(func() *Args {
					t, m := method(Nti, S28 /* stop */)
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
										t, m := method(Nt, S29 /* read */)
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
									t, m := method(one(vm, find(Nq, S30 /* lock */)), S5 /* write */)
									return call(vm, m, join(vm, t, nil))
								}())
								func() Any { a := one(vm, find(Nq, S31 /* jobs */)); Njobs = a; return a }()
								func() Any { a := one(vm, NewList([]Any{})); store(Nq, S31 /* jobs */, a); return a }()
								vm.da(func() *Args {
									t, m := method(one(vm, find(Nq, S30 /* lock */)), S29 /* read */)
									return call(vm, m, join(vm, t, nil))
								}())
								return join(vm, Func(func(vm *VM, aa *Args) *Args {
									vm.da(aa)
									{
										return join(vm, func() *Args {
											t, m := method(Njobs, S21 /* shift */)
											return call(vm, m, join(vm, t, nil))
										}())
									}
									return nil
								}))
							}
							return nil
						})),
						S29 /* read */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
							Nq := aa.get(0)
							noop(Nq)
							vm.da(aa)
							{
								var Njob Any
								noop(Njob)
								vm.da(func() *Args {
									t, m := method(one(vm, find(Nq, S30 /* lock */)), S5 /* write */)
									return call(vm, m, join(vm, t, nil))
								}())
								func() Any {
									a := one(vm, func() *Args {
										t, m := method(one(vm, find(Nq, S31 /* jobs */)), S21 /* shift */)
										return call(vm, m, join(vm, t, nil))
									}())
									Njob = a
									return a
								}()
								vm.da(func() *Args {
									t, m := method(one(vm, find(Nq, S30 /* lock */)), S29 /* read */)
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
									t, m := method(one(vm, find(Nq, S30 /* lock */)), S5 /* write */)
									return call(vm, m, join(vm, t, nil))
								}())
								vm.da(func() *Args {
									t, m := method(one(vm, find(Nq, S31 /* jobs */)), S18 /* push */)
									return call(vm, m, join(vm, t, Nfn))
								}())
								vm.da(func() *Args {
									t, m := method(one(vm, find(Nq, S30 /* lock */)), S29 /* read */)
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
								S30 /* lock */ : one(vm, call(vm, find(Nsync, S32 /* channel */), join(vm, Int(1)))),
								S31 /* jobs */ : one(vm, NewList([]Any{}))}), NprotoQueue)))
						}
						return nil
					}))
					store(Nsync, S33 /* queue */, a)
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
					a := one(vm, call(vm, Ngetprototype, join(vm, call(vm, find(Ntime, S34 /* now */), join(vm, nil)))))
					Np = a
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nt := aa.get(0)
						noop(Nt)
						vm.da(aa)
						{
							return join(vm, func() *Args {
								t, m := method(one(vm, func() *Args {
									t, m := method(Nt, S14 /* text */)
									return call(vm, m, join(vm, t, nil))
								}()), S15 /* json */)
								return call(vm, m, join(vm, t, nil))
							}())
						}
						return nil
					}))
					store(Np, S15 /* json */, a)
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nt := aa.get(0)
						noop(Nt)
						vm.da(aa)
						{
							return join(vm, func() *Args {
								t, m := method(one(vm, func() *Args {
									t, m := method(one(vm, func() *Args {
										t, m := method(Nt, S14 /* text */)
										return call(vm, m, join(vm, t, nil))
									}()), S35 /* split */)
									return call(vm, m, join(vm, t, Text(" ")))
								}()), S21 /* shift */)
								return call(vm, m, join(vm, t, nil))
							}())
						}
						return nil
					}))
					store(Np, S36 /* date */, a)
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
										t, m := method(Ns, S37 /* readrune */)
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
												t, m := method(Nline, S18 /* push */)
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
					store(Np, S38 /* readline */, a)
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Ns := aa.get(0)
						noop(Ns)
						vm.da(aa)
						{
							var Nok Any
							noop(Nok)
							var Nline Any
							noop(Nline)
							var Ndone Any
							noop(Ndone)
							func() Any { a := Nfalse; Ndone = a; return a }()
							return join(vm, Func(func(vm *VM, aa *Args) *Args {
								vm.da(aa)
								{
									if truth(one(vm, func() *Args {
										aa := join(vm, func() *Args {
											t, m := method(Ns, S38 /* readline */)
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
							var Nis Any
							noop(Nis)
							var Nfile Any
							noop(Nfile)
							var Ncontent Any
							noop(Ncontent)
							if truth(one(vm, func() *Args {
								aa := join(vm, call(vm, find(Nio, S39 /* open */), join(vm, Npath, Text("r"))))
								Nis = aa.get(0)
								Nfile = aa.get(1)
								return aa
							}())) {
								{
									if truth(one(vm, func() *Args {
										aa := join(vm, func() *Args {
											t, m := method(Nfile, S40 /* readall */)
											return call(vm, m, join(vm, t, nil))
										}())
										Nis = aa.get(0)
										Ncontent = aa.get(1)
										return aa
									}())) {
										{
											vm.da(func() *Args {
												t, m := method(Nfile, S41 /* close */)
												return call(vm, m, join(vm, t, nil))
											}())
											return join(vm, func() *Args {
												t, m := method(Ncontent, S14 /* text */)
												return call(vm, m, join(vm, t, nil))
											}())
										}
									}
								}
							}
							return join(vm, Nnil)
						}
						return nil
					}))
					store(Np, S42 /* slurp */, a)
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
					a := one(vm, call(vm, Ngetprototype, join(vm, call(vm, find(Nsync, S32 /* channel */), join(vm, Int(1))))))
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
										t, m := method(Nchan, S29 /* read */)
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
		vm.da(call(vm, Nassert, join(vm, func() *Args {
			aa := join(vm, call(vm, find(Nsql, S39 /* open */), join(vm, Text("postgres"), Text("..."))))
			Nok = aa.get(0)
			Ndb = aa.get(1)
			return aa
		}())))
		defer func() {
			func() *Args {
				t, m := method(Ndb, S41 /* close */)
				return call(vm, m, join(vm, t, nil))
			}()
		}()
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			aa := join(vm, func() *Args {
				t, m := method(Ndb, S43 /* query */)
				return call(vm, m, join(vm, t, Text("select *, $2 as c9 from t1 where c2 = $1"), Text("hello"), Nnil))
			}())
			Nok = aa.get(0)
			Nres = aa.get(1)
			return aa
		}())))
		defer func() {
			func() *Args {
				t, m := method(Nres, S41 /* close */)
				return call(vm, m, join(vm, t, nil))
			}()
		}()
		loop(func() {
			it := iterate(Nres)
			for {
				aa := it(vm, nil)
				if aa.get(0) == nil {
					vm.da(aa)
					break
				}
				vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
					Nrow := aa.get(0)
					noop(Nrow)
					vm.da(aa)
					{
						loop(func() {
							it := iterate(Nrow)
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
										vm.da(call(vm, Nprint, join(vm, call(vm, Ntype, join(vm, Nv)), Nk, Nv)))
									}
									return nil
								}), aa))
							}
						})
					}
					return nil
				}), aa))
			}
		})
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			aa := join(vm, func() *Args {
				t, m := method(Ndb, S44 /* begin */)
				return call(vm, m, join(vm, t, nil))
			}())
			Nok = aa.get(0)
			Ntx = aa.get(1)
			return aa
		}())))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(Ntx, S45 /* rollback */)
			return call(vm, m, join(vm, t, nil))
		}())))
	}
}
