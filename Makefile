build:
	docker compose build

up-app:
	docker compose up -d app

up: up-app

down:
	docker compose down

restart: down up

app-shell:
	docker compose exec app sh
