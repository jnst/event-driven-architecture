.PHONY: start format

start:
	go run *.go

format:
	gofmt -w *.go **/*.go
	goimports -w *.go **/*.go
