.PHONY: gen test

gen:
	go generate ./...

test: gen
	go test -v ./...
