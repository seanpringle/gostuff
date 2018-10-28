package main

import "net/http"
import "fmt"
import "github.com/gorilla/websocket"
import "golang.org/x/net/html"

var libHttp *Map
var protoHttpReq *Map
var protoWebSock *Map

type HttpReq struct {
	r    *http.Request
	w    http.ResponseWriter
	vals *Map
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
			Text("serve"): Func(func(vm *VM, aa *Args) *Args {
				addr := totext(aa.get(0))
				routes := aa.get(1).(*Map)
				vm.da(aa)
				for pattern, handler := range routes.data {
					http.HandleFunc(totext(pattern), func(w http.ResponseWriter, r *http.Request) {
						r.ParseForm()
						vm := &VM{}
						call(vm, handler, join(vm, HttpReq{r: r, w: w}))
					})
				}
				return join(vm, NewStatus(http.ListenAndServe(addr, nil)))
			}),
		})
		Nhttp = libHttp

		protoText.Set(Text("escape"), Func(func(vm *VM, aa *Args) *Args {
			str := totext(aa.get(0))
			vm.da(aa)
			return join(vm, Text(html.EscapeString(str)))
		}))

		protoHttpReq = NewMap(MapData{

			Text("get"): Func(func(vm *VM, aa *Args) *Args {
				req := aa.get(0).(HttpReq)
				key := totext(aa.get(1))
				vm.da(aa)
				if v, ok := req.r.Form[key]; ok {
					return join(vm, NewStatus(nil), Blob(v[0]))
				}
				return join(vm, NewStatus(fmt.Errorf("key not found: %s", key)))
			}),

			Text("write"): Func(func(vm *VM, aa *Args) *Args {
				req := aa.get(0).(HttpReq)
				data := totext(aa.get(1))
				vm.da(aa)
				length, err := req.w.Write([]byte(data))
				return join(vm, NewStatus(err), Int(length))
			}),

			Text("websocket"): Func(func(vm *VM, aa *Args) *Args {
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

			Text("text"):   Int(websocket.TextMessage),
			Text("binary"): Int(websocket.BinaryMessage),

			Text("read"): Func(func(vm *VM, aa *Args) *Args {
				ws := aa.get(0).(WebSock)
				vm.da(aa)
				mt, payload, err := ws.c.ReadMessage()
				return join(vm, NewStatus(err), Int(mt), Blob(payload))
			}),

			Text("write"): Func(func(vm *VM, aa *Args) *Args {
				ws := aa.get(0).(WebSock)
				mt := aa.get(1).(IntIsh).Int()
				data := aa.get(2).(BlobIsh).Blob()
				vm.da(aa)
				err := ws.c.WriteMessage(int(mt), []byte(data))
				return join(vm, NewStatus(err))
			}),

			Text("close"): Func(func(vm *VM, aa *Args) *Args {
				ws := aa.get(0).(WebSock)
				vm.da(aa)
				err := ws.c.Close()
				return join(vm, NewStatus(err))
			}),
		})
		protoWebSock.meta = protoDef

	})
}
