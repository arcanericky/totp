VERSION=1.0.0-incubation
VERSION_INJECT=github.com/arcanericky/totp/cmd.versionText
SRCS=*.go cmd/**/*.go commands/*.go
MAIN=./cmd/...
EXECUTABLE=bin/totp

LINUX=$(EXECUTABLE)-linux
DARWIN=$(EXECUTABLE)-darwin
WINDOWS=$(EXECUTABLE)-windows

LINUX_AMD64=$(LINUX)-amd64
DARWIN_AMD64=$(DARWIN)-amd64
DARWIN_ARM64=$(DARWIN)-arm64
WINDOWS_AMD64=$(WINDOWS)-amd64.exe

LINUX_386=$(LINUX)-386
WINDOWS_386=$(WINDOWS)-386.exe

LINUX_ARM32=$(LINUX)-arm
LINUX_ARM64=$(LINUX)-arm64

all: linux-amd64 windows-amd64 darwin-amd64 darwin-arm64 linux-arm linux-arm64

linux-amd64: $(LINUX_AMD64)

windows-amd64: $(WINDOWS_AMD64)

darwin-amd64: $(DARWIN_AMD64)

darwin-arm64: $(DARWIN_ARM64)

linux-arm: $(LINUX_ARM32)

linux-arm64: $(LINUX_ARM64)

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic . ./commands
	go tool cover -html=coverage.txt -o coverage.html

$(LINUX_AMD64) $(DARWIN_AMD64) $(DARWIN_ARM64) $(LINUX_ARM32) $(LINUX_ARM64) $(WINDOWS_AMD64): $(SRCS)
	$(eval GOOS := $(shell echo $@ | awk -F - '{print $$2}'))
	$(eval GOARCH := $(shell echo $@ | awk -F - '{print $$3}' | sed 's/.exe//'))
	$(eval LDFLAGS := "-X $(VERSION_INJECT)=$(shell sh scripts/get-version.sh)")
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $@ -ldflags $(LDFLAGS) $(MAIN)

clean:
	rm -rf bin
