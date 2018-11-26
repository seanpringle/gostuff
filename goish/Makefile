all:
	rm -f test/*.go
	cp base.go.tpl test/base.go
	go build && cat lib.goi test2.goi | ./goish | gofmt > test/main.go && cd test && go run *.go

bug:
	rm -f test/*.go
	cp base.go.tpl test/base.go
	go build && cat test.goi | ./goish | gofmt > test/main.go && cd test && go run *.go

web:
	rm -f test/*.go
	cp base.go.tpl test/base.go
	cp http.go.tpl test/http.go
	cp db.go.tpl test/db.go
	go build && cat lib.goi web.goi | ./goish | gofmt > test/main.go && cd test && go run *.go

db:
	rm -f test/*.go
	cp base.go.tpl test/base.go
	cp db.go.tpl test/db.go
	go build && cat lib.goi db.goi | ./goish | gofmt > test/main.go && cd test && go run *.go

fib:
	cat lib.goi fib.goi | ./goish | gofmt > test/main.go

sim:
	cat lib.goi sim.goi | ./goish | gofmt > test/main.go && cd test && go run *.go

benchmark-build:
	go build bench.go

benchmark:
	./bench
