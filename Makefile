configure:
	ln -s "${HOME}/.ssh" .ssh
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migration-create: ## usage: 'make migration-create name="{{migration-name}}"'
	migrate create -dir migrations -ext .sql $(name)

migration-down: ## usage: 'make migration-down count={{count}}'
	migrate -path ./migrations -database $(DATABASE_URL) down $(count)

run-services:
	docker-compose up db redis

dev:
	docker-compose build
	docker-compose up &	npx nodemon --watch pkg --ext ".go" --exec docker-compose restart api

generate:
	go generate ./...

test:
	go test -v ./...

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out