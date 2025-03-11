VERSION := $(shell cat ./VERSION)

.PHONY: run build ver

run: ver
	bun run dev

build: ver
	bun run build

ver:
	cd script && node inject-version.js $(VERSION)