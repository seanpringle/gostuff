package main

import "net/http"
import "github.com/gorilla/websocket"

var libHttp *Map
var protoHttpReq *Map
var protoWebSock *Map

type HttpReq struct {
	r *http.Request
	w http.ResponseWriter
}

func (r HttpReq) Type() string {
	return "http-request"
}

func (r HttpReq) String() string {
	return r.Type()
}

func (r HttpReq) Lib() *Map {
	return protoHttpReq
}

type WebSock struct {
	c *websocket.Conn
}

func (r WebSock) Type() string {
	return "web-socket"
}

func (r WebSock) String() string {
	return r.Type()
}

func (r WebSock) Lib() *Map {
	return protoWebSock
}

var Nhttp Any

func init() {
	onInit = append(onInit, func() {

		libHttp = NewMap(MapData{
			Str("serve"): Func(func(vm *VM, aa *Args) *Args {
				addr := tostring(aa.get(0))
				routes := aa.get(1).(*Map)
				vm.da(aa)
				for pattern, handler := range routes.data {
					http.HandleFunc(tostring(pattern), func(w http.ResponseWriter, r *http.Request) {
						vm := &VM{}
						call(vm, handler, join(vm, HttpReq{r, w}))
					})
				}
				return join(vm, NewStatus(http.ListenAndServe(addr, nil)))
			}),
		})
		Nhttp = libHttp

		protoHttpReq = NewMap(MapData{

			Str("write"): Func(func(vm *VM, aa *Args) *Args {
				req := aa.get(0).(HttpReq)
				data := tostring(aa.get(1))
				vm.da(aa)
				length, err := req.w.Write([]byte(data))
				return join(vm, NewStatus(err), Int(length))
			}),

			Str("websocket"): Func(func(vm *VM, aa *Args) *Args {
				req := aa.get(0).(HttpReq)
				var wsUpgrader = websocket.Upgrader{
					ReadBufferSize:  1024,
					WriteBufferSize: 1024,
					CheckOrigin: func(r *http.Request) bool {
						return true
					},
				}
				conn, err := wsUpgrader.Upgrade(req.w, req.r, nil)
				return join(vm, NewStatus(err), WebSock{conn})
			}),
		})
		protoHttpReq.meta = protoDef

		protoWebSock = NewMap(MapData{

			Str("text"):   Int(websocket.TextMessage),
			Str("binary"): Int(websocket.BinaryMessage),

			Str("read"): Func(func(vm *VM, aa *Args) *Args {
				ws := aa.get(0).(WebSock)
				vm.da(aa)
				mt, payload, err := ws.c.ReadMessage()
				return join(vm, NewStatus(err), Int(mt), Byte(payload))
			}),

			Str("write"): Func(func(vm *VM, aa *Args) *Args {
				ws := aa.get(0).(WebSock)
				mt := aa.get(1).(IntIsh).Int()
				data := aa.get(2).(ByteIsh).Byte()
				vm.da(aa)
				err := ws.c.WriteMessage(int(mt), []byte(data))
				return join(vm, NewStatus(err))
			}),

			Str("close"): Func(func(vm *VM, aa *Args) *Args {
				ws := aa.get(0).(WebSock)
				vm.da(aa)
				err := ws.c.Close()
				return join(vm, NewStatus(err))
			}),
		})
		protoWebSock.meta = protoDef

	})
}
