VERSION := $(shell cat ./VERSION)

.PHONY: run build ver

run: ver
	npm run dev

build: ver
	npm run build

ver:
	cd script && node inject-version.js $(VERSION)