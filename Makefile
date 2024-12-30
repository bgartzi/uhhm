ifndef VERBOSE
.SILENT:
endif

APPNAME := uhhm
TARGETPATH := ./cmd/uhhm
BUILDDIR := $(CURDIR)/bin

ifeq ($(PREFIX),)
PREFIX := /usr/local
endif
ifeq ($(BINDIR),)
BINDIR := $(PREFIX)/bin
endif


.PHONY: build
build:
	mkdir -p $(BUILDDIR)
	go build -o $(BUILDDIR)/$(APPNAME) $(TARGETPATH)

.PHONY: tidy
tidy:
	go mod tidy
	go mod verify

.PHONY: format
format:
	gofmt -s -w $(CURDIR)

.PHONY: clean
clean:
	go clean
	rm -rf $(BUILDDIR)

.PHONY: install
install:
	@sudo install $(BUILDDIR)/$(APPNAME) $(BINDIR)/$(APPNAME)
	@sudo install $(shell find $(shell go env GOMODCACHE)/github.com/urfave/cli/v*/autocomplete -name bash_autocomplete | sort -r | head -1) /etc/bash_completion.d/$(APPNAME)

.PHONY: uninstall
uninstall:
	@sudo rm $(BINDIR)/$(APPNAME) /etc/bash_completion.d/$(APPNAME)
