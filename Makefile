test-up:
	docker compose up -d redis

test:
	make down
	make test-up
	@go test -v ./...
	make down

up:
	docker compose up -d --build

down:
	docker compose down

