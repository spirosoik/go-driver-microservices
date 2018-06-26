 # Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
PACKAGE  = gateway
DATE ?= $(shell date +%FT%T%z)
VERSION = v$(shell git rev-parse --short HEAD)
LD_FLAGS = "-X 'cmd.version=$(VERSION) -X $(PACKAGE)/cmd.BuildDate=$(DATE)'"

NSQ_ADDR=nsqd:4150
HTTP_ADDR=:8080

all: clean build

build:
	@docker build -t $(PACKAGE) --build-arg NSQ_ADDR=$(NSQ_ADDR) --build-arg  HTTP_ADDR=$(HTTP_ADDR) .

clean:
	$(GOCLEAN)
	@rm -rf bin

deps:
	@dep ensure

test: deps
	$(GOTEST) -v ./...