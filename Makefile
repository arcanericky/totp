VERSION=1.0.0-incubation
VERSION_INJECT=github.com/arcanericky/totp/cmd.versionText
SRCS=*.go totp/*.go cmd/*.go
export GO111MODULE=on

all: linux-amd64 windows-amd64 darwin-amd64

linux-amd64: bin/totp

windows-amd64: bin/totp.exe

darwin-amd64: bin/totp-darwin

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt -o coverage.html

bin/totp.exe: $(SRCS)
	GOOS=windows GOARCH=amd64 go build -o $@ -ldflags "-X $(VERSION_INJECT)=$(shell sh scripts/get-version.sh)" github.com/arcanericky/totp/totp

bin/totp: $(SRCS)
	GOOS=linux GOARCH=amd64 go build -o $@ -ldflags "-X $(VERSION_INJECT)=$(shell sh scripts/get-version.sh)" github.com/arcanericky/totp/totp

bin/totp-darwin: $(SRCS)
	GOOS=darwin GOARCH=amd64 go build -o $@ -ldflags "-X $(VERSION_INJECT)=$(shell sh scripts/get-version.sh)" github.com/arcanericky/totp/totp

clean:
	rm -rf bin
