all:
	go build && cat test2.lua | ./goish | gofmt > test/main.go && cd test && go run main.go

benchmark-build:
	go build bench.go

benchmark:
	./bench
