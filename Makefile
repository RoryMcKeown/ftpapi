#!/usr/bin/env make 

run: setup
	go run github.com/rorymckeown/ftpapi
	
setup:
	go get ./...
	
test: setup
	go test ./... -coverprofile=/tmp/coverage.out
	
cover: test
	go tool cover -html=/tmp/coverage.out
		
clean:
	go clean	
		 