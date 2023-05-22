
build:
# Build for M1 OSX
	@env GOOS=darwin GOARCH=arm64 go build -o tome-m1-darwin ./protocol/v1/librarian/cmd
# Build for Linux
	@env GOOS=linux GOARCH=amd64 go build -o tome-linux-amd64 ./protocol/v1/librarian/cmd
# Build for Windows
	@env GOOS=windows GOARCH=amd64 go build -o tome.exe ./protocol/v1/librarian/cmd

local-build:
	@go run ./protocol/v1/librarian/cmd/main.go