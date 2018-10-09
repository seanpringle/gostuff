all:
	go build && cat test.goi | ./goish | gofmt > test/main.go && cd test && go run main.go
