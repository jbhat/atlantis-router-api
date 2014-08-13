PROJECT_ROOT := $(shell pwd)
PROJECT_NAME :=	$(shell pwd | xargs basename)
VENDOR_PATH  := $(PROJECT_ROOT)/vendor
API_PATH := $(PROJECT_ROOT)/lib/atlantis-routerapi
ROUTER_PATH := $(PROJECT_ROOT)/lib/atlantis-router

GOPATH := $(GOPATH):$(PROJECT_ROOT):$(VENDOR_PATH):$(API_PATH):$(ROUTER_PATH)
export GOPATH

all: build

clean:
	rm -rf bin 

init: clean
	mkdir bin

build: init
	go build -o bin/$(PROJECT_NAME) $(PROJECT_NAME).go

test:
	echo "I don't normally test my code, but when I do, I do it on production."

fmt:
	gofmt -l -w $(PROJECT_NAME).go

