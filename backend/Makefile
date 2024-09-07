VERSION := $(shell cat ./VERSION)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: run build

run:
	cd src && go run $(LDFLAGS) main.go ;

build:
	cd src && go build $(LDFLAGS) ;