#docs
swagger:
	swag init --parseDependency --parseInternal -g pkg/api/api.go

# migration commands
migration-create: ## usage: 'make migration-create name="{{migration-name}}"'
	migrate create -dir migrations -ext .sql $(name)

migration-down: ## usage: 'make migration-down count={{count}}'
	migrate -path ./migrations -database $(DATABASE_URL) down $(count)

# local machine commands
go-test:
	go test $(shell go list ./... | grep -E '(/api/v1/|/app/|/repository/)')

go-test-cover:
	go test -coverprofile=coverage.out $(shell go list ./... | grep -E '(/api/v1/|/app/|/repository/)')
	#cat coverage.out | grep -v "_mock" >> coverage.out
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

generate:
	docker-compose run api-exec sh -c "make go-generate"

test:
	docker-compose run api-exec sh -c "make go-test"

test-cover:
	docker-compose run api-exec sh -c "make go-test-cover"

run:
	docker-compose up db redis api

run-watch:
	docker-compose up & npx nodemon --watch pkg --ext ".go" --exec docker-compose restart api

configure: docker-build
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/swaggo/swag/cmd/swag@v1.6.7