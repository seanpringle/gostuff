all:
	go build && cat lib.goi test2.goi | ./goish | gofmt > test/main.go && cd test && go run main.go

fib:
	cat lib.goi fib.goi | ./goish | gofmt > test/main.go

benchmark-build:
	go build bench.go

benchmark:
	./bench
