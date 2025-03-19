dev:
	air
generate:
	bun ./api-schema/index.ts && workpad-codegen --config=gen.yaml ./api.json

test-db-up:
	DB_URL=./test.db make migration-up

migration-create:
	echo "Enter migration name: "; \
	read NAME; \
	goose -dir ./migrations create $$NAME sql

migration-up:
	goose -dir ./migrations postgres $(DB_URL) up

migration-down:
	goose -dir ./migrations postgres $(DB_URL) down

docker-build:
	docker build --rm -t workpad .

docker-run:
	docker run --name workpad -p 8080:8080 -d --restart=unless-stopped workpad

test:
	go test -v ./...

db-up:
	docker run --name laserpg -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=laserpg -p 5432:5432 -d postgres
