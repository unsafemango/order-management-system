.PHONY: up down logs migrate

up:
	docker compose up --build -d

logs:
	docker compose logs -f --tail=200

down:
	docker compose down -v