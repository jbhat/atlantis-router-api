PROJECT_ROOT := $(shell pwd)
PROJECT_NAME :=	$(shell pwd | xargs basename)
PROJ_CLIENT_NAME := $(PROJECT_NAME)-client
VENDOR_PATH  := $(PROJECT_ROOT)/vendor
API_PATH := $(PROJECT_ROOT)/lib/atlantis-routerapi
ROUTER_PATH := $(PROJECT_ROOT)/lib/atlantis-router

GOPATH := $(GOPATH):$(PROJECT_ROOT):$(VENDOR_PATH):$(API_PATH):$(ROUTER_PATH)
export GOPATH

both: all allc

all: build 

clean:
	rm -rf bin 

init: clean
	mkdir bin

build: init
	go build -o bin/$(PROJECT_NAME) $(PROJECT_NAME).go

allc: buildc

buildc: 
	go build -o bin/$(PROJ_CLIENT_NAME) $(PROJ_CLIENT_NAME).go

test:
	@for p in `find $(API_PATH)/src -type f -name "*.go" |sed 's-\./src/\(.*\)/.*-\1-' |sort -u`; do \
		echo "Testing $$p..."; \
		go test $$p -cover -v || exit 1; \
	done
	@echo
	echo "ok."
fmt:
	@find $(API_PATH)/src -name \*.go -exec gofmt -l -w {} \;
	@gofmt -l -w $(PROJ_CLIENT_NAME).go
	@gofmt -l -w $(PROJECT_NAME).go
