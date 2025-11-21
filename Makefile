tidy ::
	@go mod tidy && go mod vendor

seed ::
	@set -a; \
	. .env; \
	set +a; \
	DB_HOST=localhost go run cmd/seed/main.go

run ::
	@set -a; \
	. .env; \
	set +a; \
	DB_HOST=localhost go run cmd/server/main.go

test ::
	@go test -v -count=1 -race ./... -coverprofile=coverage.out -covermode=atomic

lint ::
	@gofumpt -l -w .
	@go vet ./...
	@staticcheck ./...

docker-up ::
	docker compose up -d

docker-down ::
	docker compose down
