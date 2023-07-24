SHELL := /bin/bash


build:
	go build -o txparser cmd/txparser/main.go

run: build
	./txparser

test-db:
	go test -v ./internal/state/inmemorydb/inmemorydb_test.go

test-txparser:
	go test -v internal/txparser/txparser_test.go