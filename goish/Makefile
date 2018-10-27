all:
	go build && cat lib.goi test2.goi | ./goish | gofmt > test/main.go && cd test && go run main.go

bug:
	go build && cat test.goi | ./goish | gofmt > test/main.go && cd test && go run main.go

web:
	go build && cat lib.goi web.goi | ./goish | gofmt > test/main.go && cat http.go.tpl > test/http.go && cd test && go run *.go

fib:
	cat lib.goi fib.goi | ./goish | gofmt > test/main.go

benchmark-build:
	go build bench.go

benchmark:
	./bench
