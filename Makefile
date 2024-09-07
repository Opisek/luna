.PHONY: dev, down, restart, build, up, purge

dev:
	docker compose up -d --build

down:
	docker compose down

restart:
	docker compose restart

build:
	docker compose build --no-cache

up:
	docker compose up -d

purge:
	docker compose down
	sudo rm -rf /srv/luna/postgres