GO ?= go
GOBINREL = build
GOBIN = $(CURDIR)/$(GOBINREL)
GOBUILD = $(GO) build

go-version:
	@if [ $(shell $(GO) version | cut -c 16-17) -lt 24 ]; then \
		echo "minimum required Golang version is 1.24"; \
		exit 1 ;\
	fi

build: go-version
	@cd ./cmd && $(GOBUILD) -o $(GOBIN)/compare
	@echo "Run \"$(GOBIN)/compare\" to launch realtime-compare-tool"
