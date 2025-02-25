SRC := $(wildcard src/*.go)
BINARY := sshke
VERSION := $(shell cat VERSION)
REL_DIR := release
REL_LINUX_BIN := $(BINARY)-linux
REL_MACOS_BIN := $(BINARY)-macos

all: build

build:
	go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY) $(SRC)
format:
	go fmt $(SRC)
clean:
	rm -f $(BINARY)
release: check-rel-dir
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(REL_DIR)/$(REL_LINUX_BIN) $(SRC) 
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(REL_DIR)/$(REL_MACOS_BIN) $(SRC)


check-rel-dir:
	@if [ ! -d "$(REL_DIR)" ]; then \
		mkdir -p $(REL_DIR); \
		echo "Directory $(REL_DIR) created."; \
	fi

.PHONY: release
