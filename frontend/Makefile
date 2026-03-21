VERSION := $(shell cat ./VERSION)

.PHONY: run build ver

run: install ver
	bun run dev

build: install ver
	bun run build

ver: install
	cd script && bun run ./inject-version.js $(VERSION)

install: node_modules package.json bun.lock
	bun install