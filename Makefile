start:
	docker compose build
	docker compose up -d

stop:
	docker compose down

restart:
	make stop
	make start