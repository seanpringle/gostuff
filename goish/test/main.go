package main

import "log"
import "runtime/pprof"
import "os"

const S1 Text = Text("stdin")
const S6 Text = Text("stdout")
const S11 Text = Text("now")
const S20 Text = Text("max")
const S21 Text = Text("min")
const S22 Text = Text("text")
const S44 Text = Text("slurp")
const S7 Text = Text("join")
const S13 Text = Text("YMDHIS")
const S35 Text = Text("stop")
const S47 Text = Text("a")
const S27 Text = Text("pop")
const S28 Text = Text("shift")
const S29 Text = Text("clear")
const S30 Text = Text("extend")
const S36 Text = Text("lock")
const S40 Text = Text("readline")
const S52 Text = Text("tom")
const S10 Text = Text("shove")
const S12 Text = Text("format")
const S5 Text = Text("write")
const S17 Text = Text("type")
const S55 Text = Text("b")
const S9 Text = Text("read")
const S56 Text = Text("sort")
const S2 Text = Text("remove")
const S3 Text = Text("channel")
const S15 Text = Text("group")
const S23 Text = Text("json")
const S50 Text = Text("split")
const S16 Text = Text("run")
const S59 Text = Text("year")
const S14 Text = Text("stderr")
const S48 Text = Text("g")
const S54 Text = Text("harry")
const S4 Text = Text("string")
const S8 Text = Text("flush")
const S24 Text = Text("quote")
const S34 Text = Text("ticker")
const S51 Text = Text("wait")
const S25 Text = Text("insert")
const S18 Text = Text("len")
const S31 Text = Text("set")
const S32 Text = Text("get")
const S45 Text = Text("c")
const S53 Text = Text("dick")
const S58 Text = Text("huge")
const S46 Text = Text("d")
const S26 Text = Text("push")
const S38 Text = Text("queue")
const S39 Text = Text("readrune")
const S42 Text = Text("readall")
const S43 Text = Text("close")
const S37 Text = Text("jobs")
const S41 Text = Text("open")
const S49 Text = Text("m")
const S57 Text = Text("match")
const S19 Text = Text("iterate")
const S33 Text = Text("keys")

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
		var Nlist Any
		noop(Nlist)
		var Nprint Any
		noop(Nprint)
		var Nl Any
		noop(Nl)
		var Nsolo Any
		noop(Nsolo)
		var Na Any
		noop(Na)
		var Ng Any
		noop(Ng)
		var Ntrue Any
		noop(Ntrue)
		var Nstring Any
		noop(Nstring)
		var Nm Any
		noop(Nm)
		var Nmap Any
		noop(Nmap)
		var Nt Any
		noop(Nt)
		var Nsuper Any
		noop(Nsuper)
		var Ninteger Any
		noop(Ninteger)
		var Nb Any
		noop(Nb)
		var Nlen Any
		noop(Nlen)
		var Nhi Any
		noop(Nhi)
		var Nis Any
		noop(Nis)
		var Nlog Any
		noop(Nlog)
		var Ninc Any
		noop(Ninc)
		var Nc Any
		noop(Nc)
		var Nf Any
		noop(Nf)
		var Nstream Any
		noop(Nstream)
		var Nnil Any
		noop(Nnil)
		var Ns Any
		noop(Ns)
		var Nblink Any
		noop(Nblink)
		var Ndecimal Any
		noop(Ndecimal)
		var Nfalse Any
		noop(Nfalse)
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
		func() Any {
			a := one(vm, Func(func(vm *VM, aa *Args) *Args {
				Na := aa.agg(0)
				noop(Na)
				vm.da(aa)
				{
					var Nlock Any
					noop(Nlock)
					func() Any {
						a := one(vm, call(vm, find(Nsync, S3 /* channel */), join(vm, Int(1))))
						Nlock = a
						return a
					}()
					return join(vm, Func(func(vm *VM, aa *Args) *Args {
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
													t, m := method(Nv, S4 /* string */)
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
								t, m := method(Nlock, S5 /* write */)
								return call(vm, m, join(vm, t, nil))
							}())
							vm.da(func() *Args {
								t, m := method(one(vm, find(Nio, S6 /* stdout */)), S5 /* write */)
								return call(vm, m, join(vm, t, concat(one(vm, func() *Args {
									t, m := method(Na, S7 /* join */)
									return call(vm, m, join(vm, t, Text(" ")))
								}()), Text("\n"))))
							}())
							vm.da(func() *Args {
								t, m := method(one(vm, find(Nio, S6 /* stdout */)), S8 /* flush */)
								return call(vm, m, join(vm, t, nil))
							}())
							vm.da(func() *Args {
								t, m := method(Nlock, S9 /* read */)
								return call(vm, m, join(vm, t, nil))
							}())
						}
						return nil
					}))
				}
				return nil
			}))
			Nprint = a
			return a
		}()
		func() Any {
			a := one(vm, call(vm, Func(func(vm *VM, aa *Args) *Args {
				vm.da(aa)
				{
					var Nlock Any
					noop(Nlock)
					func() Any {
						a := one(vm, call(vm, find(Nsync, S3 /* channel */), join(vm, Int(1))))
						Nlock = a
						return a
					}()
					return join(vm, Func(func(vm *VM, aa *Args) *Args {
						Na := aa.agg(0)
						noop(Na)
						vm.da(aa)
						{
							vm.da(func() *Args {
								t, m := method(Na, S10 /* shove */)
								return call(vm, m, join(vm, t, func() *Args {
									t, m := method(one(vm, call(vm, find(Ntime, S11 /* now */), join(vm, nil))), S12 /* format */)
									return call(vm, m, join(vm, t, find(Ntime, S13 /* YMDHIS */)))
								}()))
							}())
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
													t, m := method(Nv, S4 /* string */)
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
								t, m := method(Nlock, S5 /* write */)
								return call(vm, m, join(vm, t, nil))
							}())
							vm.da(func() *Args {
								t, m := method(one(vm, find(Nio, S14 /* stderr */)), S5 /* write */)
								return call(vm, m, join(vm, t, concat(one(vm, func() *Args {
									t, m := method(Na, S7 /* join */)
									return call(vm, m, join(vm, t, Text(" ")))
								}()), Text("\n"))))
							}())
							vm.da(func() *Args {
								t, m := method(one(vm, find(Nio, S14 /* stderr */)), S8 /* flush */)
								return call(vm, m, join(vm, t, nil))
							}())
							vm.da(func() *Args {
								t, m := method(Nlock, S9 /* read */)
								return call(vm, m, join(vm, t, nil))
							}())
						}
						return nil
					}))
				}
				return nil
			}), join(vm, nil)))
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
		func() Any {
			a := one(vm, Func(func(vm *VM, aa *Args) *Args {
				Nfn := aa.get(0)
				noop(Nfn)
				vm.da(aa)
				{
					var Ng Any
					noop(Ng)
					func() Any { a := one(vm, call(vm, find(Nsync, S15 /* group */), join(vm, nil))); Ng = a; return a }()
					vm.da(func() *Args {
						t, m := method(Ng, S16 /* run */)
						return call(vm, m, join(vm, t, Nfn))
					}())
					return join(vm, Ng)
				}
				return nil
			}))
			Nsolo = a
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
					store(Np, S17 /* type */, a)
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
					store(Np, S18 /* len */, a)
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
					store(Np, S19 /* iterate */, a)
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
					store(Np, S20 /* max */, a)
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
					store(Np, S21 /* min */, a)
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nn := aa.get(0)
						noop(Nn)
						vm.da(aa)
						{
							return join(vm, func() *Args {
								t, m := method(Nn, S22 /* text */)
								return call(vm, m, join(vm, t, nil))
							}())
						}
						return nil
					}))
					store(Np, S23 /* json */, a)
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
					store(Np, S20 /* max */, a)
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
					store(Np, S21 /* min */, a)
					return a
				}()
				func() Any {
					a := one(vm, Func(func(vm *VM, aa *Args) *Args {
						Nn := aa.get(0)
						noop(Nn)
						vm.da(aa)
						{
							return join(vm, func() *Args {
								t, m := method(Nn, S22 /* text */)
								return call(vm, m, join(vm, t, nil))
							}())
						}
						return nil
					}))
					store(Np, S23 /* json */, a)
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
								t, m := method(Nn, S24 /* quote */)
								return call(vm, m, join(vm, t, nil))
							}())
						}
						return nil
					}))
					store(Np, S23 /* json */, a)
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
								t, m := method(Nlist, S25 /* insert */)
								return call(vm, m, join(vm, t, length(Nlist), Nval))
							}())
						}
						return nil
					}))
					store(Np, S26 /* push */, a)
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
					store(Np, S27 /* pop */, a)
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
								t, m := method(Nlist, S25 /* insert */)
								return call(vm, m, join(vm, t, Int(0), Nval))
							}())
						}
						return nil
					}))
					store(Np, S10 /* shove */, a)
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
					store(Np, S28 /* shift */, a)
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
												t, m := method(Nlist, S27 /* pop */)
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
					store(Np, S29 /* clear */, a)
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
					store(Np, S19 /* iterate */, a)
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
												t, m := method(Nlist, S26 /* push */)
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
					store(Np, S30 /* extend */, a)
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
					store(Np, S31 /* set */, a)
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
					store(Np, S32 /* get */, a)
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
					store(Np, S20 /* max */, a)
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
					store(Np, S21 /* min */, a)
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
												t, m := method(Nparts, S26 /* push */)
												return call(vm, m, join(vm, t, func() *Args {
													t, m := method(Nitem, S23 /* json */)
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
									t, m := method(Nparts, S7 /* join */)
									return call(vm, m, join(vm, t, Text(",")))
								}()), Text("]")})), S7 /* join */)
								return call(vm, m, join(vm, t, nil))
							}())
						}
						return nil
					}))
					store(Np, S23 /* json */, a)
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
									t, m := method(Nm, S33 /* keys */)
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
					store(Np, S19 /* iterate */, a)
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
					store(Np, S31 /* set */, a)
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
					store(Np, S32 /* get */, a)
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
					store(Np, S20 /* max */, a)
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
					store(Np, S21 /* min */, a)
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
														t, m := method(Nk, S22 /* text */)
														return call(vm, m, join(vm, t, nil))
													}()), S24 /* quote */)
													return call(vm, m, join(vm, t, nil))
												}())
												Nkq = a
												return a
											}()
											if eq(Nv, Nnil) {
												{
													vm.da(func() *Args {
														t, m := method(Nparts, S26 /* push */)
														return call(vm, m, join(vm, t, concat(Nkq, Text(": null"))))
													}())
													return join(vm, nil)
												}
											}
											if !eq(one(vm, find(Nv, S23 /* json */)), Nnil) {
												{
													vm.da(func() *Args {
														t, m := method(Nparts, S26 /* push */)
														return call(vm, m, join(vm, t, concat(concat(Nkq, Text(":")), one(vm, func() *Args {
															t, m := method(Nv, S23 /* json */)
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
								t, m := method(Nparts, S7 /* join */)
								return call(vm, m, join(vm, t, Text(",")))
							}())), Text("}")))
						}
						return nil
					}))
					store(Np, S23 /* json */, a)
					return a
				}()
			}
			return nil
		}), join(vm, nil)))
		vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
			vm.da(aa)
			{
				var Ntick Any
				noop(Ntick)
				var Nti Any
				noop(Nti)
				func() Any {
					a := one(vm, call(vm, find(Ntime, S34 /* ticker */), join(vm, Int(1000000))))
					Nti = a
					return a
				}()
				func() Any { a := one(vm, call(vm, Ngetprototype, join(vm, Nti))); Ntick = a; return a }()
				vm.da(func() *Args {
					t, m := method(Nti, S35 /* stop */)
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
							var Nv Any
							noop(Nv)
							var Nn Any
							noop(Nn)
							func() Any { a := Int(0); Ni = a; return a }()
							return join(vm, Func(func(vm *VM, aa *Args) *Args {
								vm.da(aa)
								{
									func() Any {
										a := one(vm, func() *Args {
											t, m := method(Nt, S9 /* read */)
											return call(vm, m, join(vm, t, nil))
										}())
										Nv = a
										return a
									}()
									if !eq(Nv, Nnil) {
										{
											vm.da(func() *Args { aa := join(vm, Ni, add(Ni, Int(1))); Nn = aa.get(0); Ni = aa.get(1); return aa }())
											return join(vm, Nn, Nv)
										}
									}
									return join(vm, Nnil)
								}
								return nil
							}))
						}
						return nil
					}))
					store(Ntick, S19 /* iterate */, a)
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
						S9 /* read */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
							Nq := aa.get(0)
							noop(Nq)
							vm.da(aa)
							{
								var Njob Any
								noop(Njob)
								vm.da(func() *Args {
									t, m := method(one(vm, find(Nq, S36 /* lock */)), S5 /* write */)
									return call(vm, m, join(vm, t, nil))
								}())
								func() Any {
									a := one(vm, func() *Args {
										t, m := method(one(vm, find(Nq, S37 /* jobs */)), S28 /* shift */)
										return call(vm, m, join(vm, t, nil))
									}())
									Njob = a
									return a
								}()
								vm.da(func() *Args {
									t, m := method(one(vm, find(Nq, S36 /* lock */)), S9 /* read */)
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
									t, m := method(one(vm, find(Nq, S36 /* lock */)), S5 /* write */)
									return call(vm, m, join(vm, t, nil))
								}())
								vm.da(func() *Args {
									t, m := method(one(vm, find(Nq, S37 /* jobs */)), S26 /* push */)
									return call(vm, m, join(vm, t, Nfn))
								}())
								vm.da(func() *Args {
									t, m := method(one(vm, find(Nq, S36 /* lock */)), S9 /* read */)
									return call(vm, m, join(vm, t, nil))
								}())
							}
							return nil
						})),
						S19 /* iterate */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
							Nq := aa.get(0)
							noop(Nq)
							vm.da(aa)
							{
								var Njobs Any
								noop(Njobs)
								vm.da(func() *Args {
									t, m := method(one(vm, find(Nq, S36 /* lock */)), S5 /* write */)
									return call(vm, m, join(vm, t, nil))
								}())
								func() Any { a := one(vm, find(Nq, S37 /* jobs */)); Njobs = a; return a }()
								func() Any { a := one(vm, NewList([]Any{})); store(Nq, S37 /* jobs */, a); return a }()
								vm.da(func() *Args {
									t, m := method(one(vm, find(Nq, S36 /* lock */)), S9 /* read */)
									return call(vm, m, join(vm, t, nil))
								}())
								return join(vm, Func(func(vm *VM, aa *Args) *Args {
									vm.da(aa)
									{
										return join(vm, func() *Args {
											t, m := method(Njobs, S28 /* shift */)
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
								S36 /* lock */ : one(vm, call(vm, find(Nsync, S3 /* channel */), join(vm, Int(1)))),
								S37 /* jobs */ : one(vm, NewList([]Any{}))}), NprotoQueue)))
						}
						return nil
					}))
					store(Nsync, S38 /* queue */, a)
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
					a := one(vm, call(vm, Ngetprototype, join(vm, call(vm, find(Ntime, S11 /* now */), join(vm, nil)))))
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
									t, m := method(Nt, S22 /* text */)
									return call(vm, m, join(vm, t, nil))
								}()), S23 /* json */)
								return call(vm, m, join(vm, t, nil))
							}())
						}
						return nil
					}))
					store(Np, S23 /* json */, a)
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
										t, m := method(Ns, S39 /* readrune */)
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
												t, m := method(Nline, S26 /* push */)
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
								t, m := method(Nline, S7 /* join */)
								return call(vm, m, join(vm, t, nil))
							}())
						}
						return nil
					}))
					store(Np, S40 /* readline */, a)
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
											t, m := method(Ns, S40 /* readline */)
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
					store(Np, S19 /* iterate */, a)
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
								aa := join(vm, call(vm, find(Nio, S41 /* open */), join(vm, Npath, Text("r"))))
								Nis = aa.get(0)
								Nfile = aa.get(1)
								return aa
							}())) {
								{
									if truth(one(vm, func() *Args {
										aa := join(vm, func() *Args {
											t, m := method(Nfile, S42 /* readall */)
											return call(vm, m, join(vm, t, nil))
										}())
										Nis = aa.get(0)
										Ncontent = aa.get(1)
										return aa
									}())) {
										{
											vm.da(func() *Args {
												t, m := method(Nfile, S43 /* close */)
												return call(vm, m, join(vm, t, nil))
											}())
											return join(vm, Nis, func() *Args {
												t, m := method(Ncontent, S22 /* text */)
												return call(vm, m, join(vm, t, nil))
											}())
										}
									}
								}
							}
							return join(vm, Nis)
						}
						return nil
					}))
					store(Np, S44 /* slurp */, a)
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
					a := one(vm, call(vm, Ngetprototype, join(vm, call(vm, find(Nsync, S3 /* channel */), join(vm, Int(1))))))
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
										t, m := method(Nchan, S9 /* read */)
										return call(vm, m, join(vm, t, nil))
									}())
								}
								return nil
							}))
						}
						return nil
					}))
					store(Np, S19 /* iterate */, a)
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
				S45 /* c */ : one(vm, NewMap(MapData{
					S46 /* d */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
						vm.da(aa)
						{
							return join(vm, Text("hello world"))
						}
						return nil
					}))})),
				S47 /* a */ : Int(1)}))
			Nt = a
			return a
		}()
		vm.da(call(vm, Nprint, join(vm, call(vm, find(one(vm, find(Nt, S45 /* c */)), S46 /* d */), join(vm, nil)))))
		func() Any { a := Int(42); store(Nt, S47 /* a */, a); return a }()
		vm.da(call(vm, Nprint, join(vm, Nt)))
		vm.da(call(vm, Nprint, join(vm, Text(""), func() *Args {
			t, m := method(Nt, S33 /* keys */)
			return call(vm, m, join(vm, t, nil))
		}())))
		vm.da(call(vm, Nprint, join(vm, add(one(vm, mul(Int(2), Int(2))), Int(3)))))
		func() Any {
			a := one(vm, NewMap(MapData{
				S48 /* g */ : one(vm, Func(func(vm *VM, aa *Args) *Args {
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
						t, m := method(Nself, S48 /* g */)
						return call(vm, m, join(vm, t, nil))
					}())
				}
				return nil
			}))
			store(Nt, S49 /* m */, a)
			return a
		}()
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(Nt, S49 /* m */)
			return call(vm, m, join(vm, t, nil))
		}())))
		func() Any { a := Text("goodbye world"); Ns = a; return a }()
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(Ns, S18 /* len */)
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
			t, m := method(Na, S31 /* set */)
			return call(vm, m, join(vm, t, Text("1"), Int(1)))
		}())
		vm.da(call(vm, Nprint, join(vm, Na)))
		func() Any {
			a := one(vm, NewMap(MapData{}))
			Nb = a
			return a
		}()
		vm.da(func() *Args {
			t, m := method(Na, S31 /* set */)
			return call(vm, m, join(vm, t, Nb, Int(2)))
		}())
		vm.da(call(vm, Nprint, join(vm, Na)))
		vm.da(func() *Args {
			t, m := method(Nb, S31 /* set */)
			return call(vm, m, join(vm, t, Text("2"), Int(2)))
		}())
		vm.da(call(vm, Nprint, join(vm, Na)))
		func() Any { a := one(vm, NewList([]Any{Int(1), Int(2), Int(3)})); Nl = a; return a }()
		vm.da(call(vm, Nprint, join(vm, Nl)))
		vm.da(func() *Args {
			t, m := method(Nl, S26 /* push */)
			return call(vm, m, join(vm, t, Int(4)))
		}())
		vm.da(call(vm, Nprint, join(vm, Nl)))
		vm.da(call(vm, Nprint, join(vm, call(vm, Ngetprototype, join(vm, Nl)))))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(Nl, S27 /* pop */)
			return call(vm, m, join(vm, t, nil))
		}())))
		vm.da(call(vm, Nprint, join(vm, Nl)))
		vm.da(call(vm, Nprint, join(vm, concat(Text("a"), Text("b")))))
		func() Any { a := Text("hi"); Nlen = a; return a }()
		vm.da(call(vm, Nprint, join(vm, Text("yo"), func() *Args {
			t, m := method(Nl, S18 /* len */)
			return call(vm, m, join(vm, t, nil))
		}())))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(Text("a,b,c"), S50 /* split */)
			return call(vm, m, join(vm, t, Text(",")))
		}())))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(one(vm, func() *Args {
				t, m := method(Text("a,b,c"), S50 /* split */)
				return call(vm, m, join(vm, t, Text(",")))
			}()), S7 /* join */)
			return call(vm, m, join(vm, t, Text(":")))
		}())))
		func() Any { a := one(vm, call(vm, find(Nsync, S3 /* channel */), join(vm, Int(10)))); Nc = a; return a }()
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
			t, m := method(Nc, S9 /* read */)
			return call(vm, m, join(vm, t, nil))
		}())))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(Nc, S9 /* read */)
			return call(vm, m, join(vm, t, nil))
		}())))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(Nc, S9 /* read */)
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
		func() Any { a := one(vm, call(vm, find(Nsync, S15 /* group */), join(vm, nil))); Ng = a; return a }()
		vm.da(func() *Args {
			t, m := method(Ng, S16 /* run */)
			return call(vm, m, join(vm, t, Nhi))
		}())
		vm.da(func() *Args {
			t, m := method(Ng, S16 /* run */)
			return call(vm, m, join(vm, t, Nhi))
		}())
		vm.da(func() *Args {
			t, m := method(Ng, S16 /* run */)
			return call(vm, m, join(vm, t, Nhi))
		}())
		vm.da(func() *Args {
			t, m := method(Ng, S51 /* wait */)
			return call(vm, m, join(vm, t, nil))
		}())
		vm.da(call(vm, Nprint, join(vm, Text("done"))))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(Nb, S32 /* get */)
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
				S52 /* tom */ :   Int(1),
				S53 /* dick */ :  Int(2),
				S54 /* harry */ : Int(43)})))
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
		func() Any { a := one(vm, call(vm, find(Nsync, S38 /* queue */), join(vm, nil))); Nblink = a; return a }()
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
				S47 /* a */ : Int(1),
				S55 /* b */ : one(vm, NewMap(MapData{
					S45 /* c */ : Int(4)}))}))
			Nm = a
			return a
		}()
		vm.da(call(vm, Nprint, join(vm, field(field(Nm, Text("b")), Text("c")))))
		func() Any { a := Int(5); store(field(Nm, Text("b")), Text("c"), a); return a }()
		vm.da(call(vm, Nprint, join(vm, field(field(Nm, Text("b")), Text("c")))))
		vm.da(call(vm, Nprint, join(vm, Text("length"), length(Nl), length(Nm))))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(Int(0), S20 /* max */)
			return call(vm, m, join(vm, t, Int(2)))
		}())))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(one(vm, NewList([]Any{Int(2), Int(4), Int(6), Int(8), Int(3)})), S56 /* sort */)
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
			t, m := method(Text("abc"), S57 /* match */)
			return call(vm, m, join(vm, t, Text("[aeiou]")))
		}())))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(Text("abc"), S57 /* match */)
			return call(vm, m, join(vm, t, Text("[aeiou]")))
		}())))
		vm.da(call(vm, Func(func(vm *VM, aa *Args) *Args {
			vm.da(aa)
			{
				vm.da(call(vm, Nprint, join(vm, Text("hi"))))
			}
			return nil
		}), join(vm, nil)))
		vm.da(call(vm, Nprint, join(vm, find(one(vm, call(vm, Ngetprototype, join(vm, Int(0)))), S58 /* huge */))))
		vm.da(call(vm, Nprint, join(vm, find(one(vm, call(vm, Ngetprototype, join(vm, Dec(1)))), S58 /* huge */))))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(one(vm, NewList([]Any{})), S30 /* extend */)
			return call(vm, m, join(vm, t, Int(3)))
		}())))
		vm.da(call(vm, Nprint, join(vm, mul(add(Int(1), Int(2)), Int(3)))))
		vm.da(call(vm, Nprint, join(vm, mul(Int(3), add(Int(1), Int(2))))))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(one(vm, NewList([]Any{Int(1), Int(2), Int(3)})), S25 /* insert */)
			return call(vm, m, join(vm, t, Int(1), Int(7)))
		}())))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(one(vm, NewList([]Any{Int(1), Int(2), Int(3)})), S25 /* insert */)
			return call(vm, m, join(vm, t, Int(0), Int(7)))
		}())))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(one(vm, NewList([]Any{Int(1), Int(2), Int(3)})), S25 /* insert */)
			return call(vm, m, join(vm, t, Int(4), Int(7)))
		}())))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(one(vm, NewList([]Any{Int(1), Int(2), Int(3)})), S2 /* remove */)
			return call(vm, m, join(vm, t, Int(0)))
		}())))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(one(vm, NewList([]Any{Int(1), Int(2), Int(3)})), S2 /* remove */)
			return call(vm, m, join(vm, t, Int(1)))
		}())))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(one(vm, NewList([]Any{Int(1), Int(2), Int(3)})), S2 /* remove */)
			return call(vm, m, join(vm, t, Int(3)))
		}())))
		func() Any { a := one(vm, NewList([]Any{Int(1), Int(2), Int(3)})); Nl = a; return a }()
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(Nl, S2 /* remove */)
			return call(vm, m, join(vm, t, Int(2)))
		}(), Nl)))
		func() Any {
			a := one(vm, NewMap(MapData{
				S55 /* b */ : one(vm, NewList([]Any{Int(1), Int(2), Int(3)}))}))
			Na = a
			return a
		}()
		vm.da(call(vm, Nprint, join(vm, length(Na))))
		if truth(one(vm, func() Any {
			var a Any
			a = one(vm, func() Any {
				a := one(vm, NewMap(MapData{
					S55 /* b */ : Int(1)}))
				Na = a
				return a
			}())
			if truth(a) {
				var b Any
				b = Bool(eq(one(vm, find(Na, S55 /* b */)), Int(1)))
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
		vm.da(call(vm, Nprint, join(vm, Text("and"), lshift(Int(1), Int(2)))))
		vm.da(call(vm, Nprint, join(vm, Text("hex"), Int(255))))
		vm.da(call(vm, Nprint, join(vm, nil)))
		vm.da(call(vm, Nprint, join(vm, nil)))
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(one(vm, call(vm, find(Ntime, S11 /* now */), join(vm, nil))), S59 /* year */)
			return call(vm, m, join(vm, t, nil))
		}())))
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
		vm.da(call(vm, Nprint, join(vm, func() *Args {
			t, m := method(one(vm, NewMap(MapData{
				S47 /* a */ : one(vm, NewMap(MapData{
					S55 /* b */ : one(vm, NewList([]Any{Int(1), Int(2), Int(3)}))}))})), S23 /* json */)
			return call(vm, m, join(vm, t, nil))
		}())))
		defer catch(vm, Func(func(vm *VM, aa *Args) *Args {
			Ns := aa.get(0)
			noop(Ns)
			vm.da(aa)
			{
				vm.da(call(vm, Nlog, join(vm, Text("caught"), Ns)))
			}
			return nil
		}))
		vm.da(call(vm, Nprint, join(vm, try(vm, call(vm, Func(func(vm *VM, aa *Args) *Args {
			vm.da(aa)
			{
				return join(vm, call(vm, Nstatus, join(vm, Nnil)), Text("hello"))
			}
			return nil
		}), join(vm, nil))))))
		vm.da(call(vm, Nprint, join(vm, try(vm, call(vm, Func(func(vm *VM, aa *Args) *Args {
			vm.da(aa)
			{
				return join(vm, call(vm, Nstatus, join(vm, Text("wtf"))), Text("world"))
			}
			return nil
		}), join(vm, nil))))))
	}
}
