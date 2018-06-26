GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
PACKAGE = driver-location
DATE ?= $(shell date +%FT%T%z)
VERSION = v$(shell git rev-parse --short HEAD)
LD_FLAGS = "-X 'cmd.version=$(VERSION) -X $(PACKAGE)/cmd.BuildDate=$(DATE)'"

# Addresses
HTTP_ADDR=:8081
NSQ_LOOKUP_ADDR=nsqlookupd:4161
REDIS_ADDR=redis:6379

all: clean build

build:
	@docker build -t $(PACKAGE) --build-arg NSQ_LOOKUP_ADDR=$(NSQ_LOOKUP_ADDR) --build-arg  HTTP_ADDR=$(HTTP_ADDR) --build-arg NSQ_LOOKUP_ADDR=$(NSQ_LOOKUP_ADDR) --build-arg  REDIS_ADDR=$(REDIS_ADDR) .

clean:
	$(GOCLEAN)
	@rm -rf bin

deps:
	@dep ensure

test: deps
	$(GOTEST) -v ./...