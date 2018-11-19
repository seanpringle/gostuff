package main

import "net/http"
import "fmt"
import "github.com/gorilla/websocket"
import "golang.org/x/net/html"
import "context"

var libHttp *Map
var protoHttpSrv *Map
var protoHttpReq *Map
var protoWebSock *Map

type HttpReq struct {
	r    *http.Request
	w    http.ResponseWriter
	d    chan struct{}
	vals *Map
}

func (r HttpReq) Type() string {
	return "http-request"
}

func (r HttpReq) String() string {
	return r.Type()
}

func (r HttpReq) Lib() Searchable {
	return protoHttpReq
}

func (r HttpReq) Bool() Bool {
	return true
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

func (r WebSock) Lib() Searchable {
	return protoWebSock
}

func (r WebSock) Bool() Bool {
	return true
}

type HttpSrv struct {
	s *http.Server
	c *Chan
}

func (r HttpSrv) Type() string {
	return "http-server"
}

func (r HttpSrv) String() string {
	return r.Type()
}

func (r HttpSrv) Lib() Searchable {
	return protoHttpSrv
}

func (r HttpSrv) Bool() Bool {
	return true
}

func (r HttpSrv) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	done := make(chan struct{})
	r.c.c <- HttpReq{r: re, w: w, d: done}
	<-done
}

var Nhttp Any

func init() {
	onInit = append(onInit, func() {

		libHttp = NewMap(MapData{

			Text("Continue"):                      toInt(Int(http.StatusContinue)),
			Text("SwitchingProtocols"):            toInt(Int(http.StatusSwitchingProtocols)),
			Text("Processing"):                    toInt(Int(http.StatusProcessing)),
			Text("OK"):                            toInt(Int(http.StatusOK)),
			Text("Created"):                       toInt(Int(http.StatusCreated)),
			Text("Accepted"):                      toInt(Int(http.StatusAccepted)),
			Text("NonAuthoritativeInfo"):          toInt(Int(http.StatusNonAuthoritativeInfo)),
			Text("NoContent"):                     toInt(Int(http.StatusNoContent)),
			Text("ResetContent"):                  toInt(Int(http.StatusResetContent)),
			Text("PartialContent"):                toInt(Int(http.StatusPartialContent)),
			Text("MultiStatus"):                   toInt(Int(http.StatusMultiStatus)),
			Text("AlreadyReported"):               toInt(Int(http.StatusAlreadyReported)),
			Text("IMUsed"):                        toInt(Int(http.StatusIMUsed)),
			Text("MultipleChoices"):               toInt(Int(http.StatusMultipleChoices)),
			Text("MovedPermanently"):              toInt(Int(http.StatusMovedPermanently)),
			Text("Found"):                         toInt(Int(http.StatusFound)),
			Text("SeeOther"):                      toInt(Int(http.StatusSeeOther)),
			Text("NotModified"):                   toInt(Int(http.StatusNotModified)),
			Text("UseProxy"):                      toInt(Int(http.StatusUseProxy)),
			Text("TemporaryRedirect"):             toInt(Int(http.StatusTemporaryRedirect)),
			Text("PermanentRedirect"):             toInt(Int(http.StatusPermanentRedirect)),
			Text("BadRequest"):                    toInt(Int(http.StatusBadRequest)),
			Text("Unauthorized"):                  toInt(Int(http.StatusUnauthorized)),
			Text("PaymentRequired"):               toInt(Int(http.StatusPaymentRequired)),
			Text("Forbidden"):                     toInt(Int(http.StatusForbidden)),
			Text("NotFound"):                      toInt(Int(http.StatusNotFound)),
			Text("MethodNotAllowed"):              toInt(Int(http.StatusMethodNotAllowed)),
			Text("NotAcceptable"):                 toInt(Int(http.StatusNotAcceptable)),
			Text("ProxyAuthRequired"):             toInt(Int(http.StatusProxyAuthRequired)),
			Text("RequestTimeout"):                toInt(Int(http.StatusRequestTimeout)),
			Text("Conflict"):                      toInt(Int(http.StatusConflict)),
			Text("Gone"):                          toInt(Int(http.StatusGone)),
			Text("LengthRequired"):                toInt(Int(http.StatusLengthRequired)),
			Text("PreconditionFailed"):            toInt(Int(http.StatusPreconditionFailed)),
			Text("RequestEntityTooLarge"):         toInt(Int(http.StatusRequestEntityTooLarge)),
			Text("RequestURITooLong"):             toInt(Int(http.StatusRequestURITooLong)),
			Text("UnsupportedMediaType"):          toInt(Int(http.StatusUnsupportedMediaType)),
			Text("RequestedRangeNotSatisfiable"):  toInt(Int(http.StatusRequestedRangeNotSatisfiable)),
			Text("ExpectationFailed"):             toInt(Int(http.StatusExpectationFailed)),
			Text("Teapot"):                        toInt(Int(http.StatusTeapot)),
			Text("MisdirectedRequest"):            toInt(Int(http.StatusMisdirectedRequest)),
			Text("UnprocessableEntity"):           toInt(Int(http.StatusUnprocessableEntity)),
			Text("Locked"):                        toInt(Int(http.StatusLocked)),
			Text("FailedDependency"):              toInt(Int(http.StatusFailedDependency)),
			Text("UpgradeRequired"):               toInt(Int(http.StatusUpgradeRequired)),
			Text("PreconditionRequired"):          toInt(Int(http.StatusPreconditionRequired)),
			Text("TooManyRequests"):               toInt(Int(http.StatusTooManyRequests)),
			Text("RequestHeaderFieldsTooLarge"):   toInt(Int(http.StatusRequestHeaderFieldsTooLarge)),
			Text("UnavailableForLegalReasons"):    toInt(Int(http.StatusUnavailableForLegalReasons)),
			Text("InternalServerError"):           toInt(Int(http.StatusInternalServerError)),
			Text("NotImplemented"):                toInt(Int(http.StatusNotImplemented)),
			Text("BadGateway"):                    toInt(Int(http.StatusBadGateway)),
			Text("ServiceUnavailable"):            toInt(Int(http.StatusServiceUnavailable)),
			Text("GatewayTimeout"):                toInt(Int(http.StatusGatewayTimeout)),
			Text("HTTPVersionNotSupported"):       toInt(Int(http.StatusHTTPVersionNotSupported)),
			Text("VariantAlsoNegotiates"):         toInt(Int(http.StatusVariantAlsoNegotiates)),
			Text("InsufficientStorage"):           toInt(Int(http.StatusInsufficientStorage)),
			Text("LoopDetected"):                  toInt(Int(http.StatusLoopDetected)),
			Text("NotExtended"):                   toInt(Int(http.StatusNotExtended)),
			Text("NetworkAuthenticationRequired"): toInt(Int(http.StatusNetworkAuthenticationRequired)),

			Text("server"): Func(func(vm *VM, aa *Args) *Args {
				addr := totext(aa.get(0))
				vm.da(aa)

				srv := HttpSrv{
					s: &http.Server{Addr: addr},
					c: NewChan(8),
				}
				srv.s.Handler = srv

				return join(vm, srv)
			}),
		})
		Nhttp = libHttp

		protoHttpSrv = NewMap(MapData{
			Text("serve"): Func(func(vm *VM, aa *Args) *Args {
				s := aa.get(0).(HttpSrv)
				go func() {
					err := s.s.ListenAndServe()
					if err != nil && err != http.ErrServerClosed {
						panic(err)
					}
					close(s.c.c)
				}()
				return join(vm, s.c)
			}),
			Text("stop"): Func(func(vm *VM, aa *Args) *Args {
				s := aa.get(0).(HttpSrv)
				err := s.s.Shutdown(context.Background())
				return join(vm, NewStatus(err))
			}),
		})
		protoHttpSrv.meta = protoDef

		protoText.Set(Text("htmlescape"), Func(func(vm *VM, aa *Args) *Args {
			str := totext(aa.get(0))
			vm.da(aa)
			return join(vm, Text(html.EscapeString(str)))
		}))

		protoHttpReq = NewMap(MapData{

			Text("static"): Func(func(vm *VM, aa *Args) *Args {
				r := aa.get(0).(HttpReq)
				p := totext(aa.get(1))
				vm.da(aa)
				http.FileServer(http.Dir(p)).ServeHTTP(r.w, r.r)
				r.d <- struct{}{}
				return join(vm, NewStatus(nil))
			}),

			Text("path"): Func(func(vm *VM, aa *Args) *Args {
				r := aa.get(0).(HttpReq)
				return join(vm, Text(r.r.URL.Path))
			}),

			Text("get"): Func(func(vm *VM, aa *Args) *Args {
				req := aa.get(0).(HttpReq)
				key := totext(aa.get(1))
				vm.da(aa)
				v := req.r.FormValue(key)
				if v == "" {
					return join(vm, NewStatus(fmt.Errorf("key not found: %s", key)))
				}
				return join(vm, NewStatus(nil), Text(v))
			}),

			Text("code"): Func(func(vm *VM, aa *Args) *Args {
				req := aa.get(0).(HttpReq)
				code := int(aa.get(1).(IntIsh).Int())
				vm.da(aa)
				req.w.WriteHeader(code)
				return join(vm, req)
			}),

			Text("header"): Func(func(vm *VM, aa *Args) *Args {
				req := aa.get(0).(HttpReq)
				key := totext(aa.get(1))
				val := totext(aa.get(2))
				vm.da(aa)
				req.w.Header().Add(key, val)
				return join(vm, req)
			}),

			Text("write"): Func(func(vm *VM, aa *Args) *Args {
				req := aa.get(0).(HttpReq)
				data := totext(aa.get(1))
				vm.da(aa)
				length, err := req.w.Write([]byte(data))
				req.d <- struct{}{}
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
