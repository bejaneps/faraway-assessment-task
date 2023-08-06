include .env
export ENV_FILE=.env

GOENVS=GOOS=linux CGO_ENABLED=0 GOARCH=amd64

lint:
	@golangci-lint run -v ./...

hook:
	@cp tools/dist/pre-commit .git/hooks/pre-commit && chmod +x .git/hooks/pre-commit
	@cp tools/dist/pre-push .git/hooks/pre-push && chmod +x .git/hooks/pre-push

deps:
	@go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3

update-mocks:
	@mockery --all --dir internal/api --output internal/api/mocks/ --outpkg mocks

test-unit:
	@go test -json -race -tags=unit -p=1 -short -count=1 ./... 2>&1 | tee /tmp/gotest.log | gotestfmt

test-integration:
	@go test -json -race -tags=integration -p=1 -short -count=1 ./internal/pkg/repository_tests 2>&1 | tee /tmp/gotest.log | gotestfmt

build:
	@$(GOENVS) go build -ldflags="-s -w" -o dist/server ./cmd/server
	@$(GOENVS) go build -ldflags="-s -w" -o dist/client ./cmd/client
	@chmod +x dist/server dist/client

build-local:
	@go build -o bin/server ./cmd/server
	@go build -o bin/client ./cmd/client

stop:
	@docker compose --profile default down --remove-orphans
	@docker compose --profile testing down --remove-orphans --timeout=3

run: stop build
	@docker compose --profile default up --build
