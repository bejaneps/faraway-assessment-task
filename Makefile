include .env
export ENV_FILE=.env

TEST_PKGS=$(shell go list -f '{{if .TestGoFiles}}{{.ImportPath}}{{end}}' ./...)
GOENVS=GOOS=linux CGO_ENABLED=0 GOARCH=amd64

lint:
	@golangci-lint run -v ./...

hook:
	@cp tools/dist/pre-commit .git/hooks/pre-commit && chmod +x .git/hooks/pre-commit
	@cp tools/dist/pre-push .git/hooks/pre-push && chmod +x .git/hooks/pre-push

deps:
	@go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3

install-mocks:
	@go install github.com/vektra/mockery/v2@v2.23.4

update-mocks:
	@mockery --all --dir internal/pkg/transport --output internal/pkg/transport/mocks/ --outpkg mocks
	@mockery --all --dir internal/pkg/db --output internal/pkg/db/mocks/ --outpkg mocks
	@mockery --all --dir internal/pkg/cache --output internal/pkg/cache/mocks/ --outpkg mocks
	@mockery --all --dir internal/repository/server --output internal/repository/server/mocks/ --outpkg mocks
	@mockery --all --dir internal/service/server --output internal/service/server/mocks/ --outpkg mocks

test-unit:
	LOG_ENABLED=false \
	  go test ${TEST_PKGS} -json -race -p=1 -short -count=1 2>&1 | tee /tmp/gotest.log | gotestfmt

test-integration: stop clean setup
	LOG_ENABLED=false \
	CACHE_URL=127.0.0.1:6379 \
	DB_DSN=postgresql://postgres:postgres@127.0.0.1:5432/faraway-db?sslmode=disable \
	  go test ${TEST_PKGS} -json -p=1 -race -count=1 2>&1 | tee /tmp/gotest.log | gotestfmt

test-cover: stop clean setup
	LOG_ENABLED=false \
	CACHE_URL=127.0.0.1:6379 \
	DB_DSN=postgresql://postgres:postgres@127.0.0.1:5432/faraway-db?sslmode=disable \
	  go test ${TEST_PKGS} -p=1 -coverprofile=cover.out && go tool cover -html=cover.out

docker-test-integration: stop clean
	@docker compose --profile database up -d
	@docker compose --profile cache up -d
	until docker compose exec -T postgresdb pg_isready -h postgresdb; do sleep 1; done
	@docker compose --profile migration build migrations
	@docker compose --profile migration run migrations
	@docker compose --profile integration-test up --build --abort-on-container-exit --exit-code-from integration-test
	@docker compose --profile database down
	@docker compose --profile cache down
	@docker compose --profile integration-test down

build:
	@$(GOENVS) go build -ldflags="-s -w" -o dist/server ./cmd/server
	@$(GOENVS) go build -ldflags="-s -w" -o dist/client ./cmd/client
	@chmod +x dist/server dist/client

build-local:
	@go build -o bin/server ./cmd/server
	@go build -o bin/client ./cmd/client

#TODO: wait for redis as well
setup: build
	@docker compose --profile database up -d
	@docker compose --profile cache up -d
	until docker compose exec -T postgresdb pg_isready -h postgresdb; do sleep 1; done
	@docker compose --profile migration build migrations
	@docker compose --profile migration run migrations

stop:
	@docker compose down --remove-orphans

run: stop clean setup
	@docker compose --profile default up --build --abort-on-container-exit

clean:
	@docker volume rm -f faraway-assessment-task_db-data
	@docker volume rm -f faraway-assessment-task_redis
