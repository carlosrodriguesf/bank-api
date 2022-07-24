configure:
	ln -s "${HOME}/.ssh" .ssh
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest


# migration commands
migration-create: ## usage: 'make migration-create name="{{migration-name}}"'
	migrate create -dir migrations -ext .sql $(name)

migration-down: ## usage: 'make migration-down count={{count}}'
	migrate -path ./migrations -database $(DATABASE_URL) down $(count)

# local machine commands
go-test:
	go test ./...

go-test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

go-run:
	go run pkg/main.go

go-generate:
	go generate ./...

# docker commands
docker-build:
	docker-compose build

run-services:
	docker-compose up db redis

generate: docker-build
	docker-compose run api sh -c "make go-generate"

test: docker-build
	docker-compose run api sh -c "make go-test"

test-cover: docker-build
	docker-compose run api sh -c "make go-test-cover"

run: docker-build
	docker-compose up

run-watch: docker-build
	docker-compose up & npx nodemon --watch pkg --ext ".go" --exec docker-compose restart api