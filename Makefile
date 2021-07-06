.PHONY: start format

start:
	go run *.go

fmt:
	gofmt -w *.go **/*.go
	goimports -w *.go **/*.go
