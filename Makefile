
build:
	@go build -o tome ./protocol/v1/librarian/cmd

local-build:
	@go run ./protocol/v1/librarian/cmd/main.go