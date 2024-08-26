VERSION := $(shell cat ./VERSION)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: run build

run:
	go run $(LDFLAGS) main.go ;

build:
	go build $(LDFLAGS) ;