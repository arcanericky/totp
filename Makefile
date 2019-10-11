VERSION=1.0.0-incubation
VERSION_INJECT=github.com/arcanericky/totp/cmd.versionText
SRCS=*.go totp/*.go cmd/*.go

EXECUTABLE=bin/totp

LINUX=$(EXECUTABLE)-linux
DARWIN=$(EXECUTABLE)-darwin
WINDOWS=$(EXECUTABLE)-windows

LINUX_AMD64=$(LINUX)-amd64
DARWIN_AMD64=$(DARWIN)-amd64
WINDOWS_AMD64=$(WINDOWS)-amd64.exe

LINUX_386=$(LINUX)-386
WINDOWS_386=$(WINDOWS)-386.exe

all: linux-amd64 windows-amd64 darwin-amd64

linux-amd64: $(LINUX_AMD64)

windows-amd64: $(WINDOWS_AMD64)

darwin-amd64: $(DARWIN_AMD64)

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic . ./cmd
	go tool cover -html=coverage.txt -o coverage.html

$(WINDOWS_AMD64): $(SRCS)
	GOOS=windows GOARCH=amd64 go build -o $@ -ldflags "-X $(VERSION_INJECT)=$(shell sh scripts/get-version.sh)" github.com/arcanericky/totp/totp

$(LINUX_AMD64): $(SRCS)
	GOOS=linux GOARCH=amd64 go build -o $@ -ldflags "-X $(VERSION_INJECT)=$(shell sh scripts/get-version.sh)" github.com/arcanericky/totp/totp

$(DARWIN_AMD64): $(SRCS)
	GOOS=darwin GOARCH=amd64 go build -o $@ -ldflags "-X $(VERSION_INJECT)=$(shell sh scripts/get-version.sh)" github.com/arcanericky/totp/totp

clean:
	rm -rf bin
