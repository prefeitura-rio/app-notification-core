run:
    go run cmd/server/main.go

build:
    go build -o bin/server cmd/server/main.go

dev:
    air

test:
    go test ./... -v

test-coverage:
    go test ./... -coverprofile=coverage.out
    go tool cover -html=coverage.out

install:
    go mod download
    go mod tidy

install-tools:
    go install github.com/swaggo/swag/cmd/swag@latest
    go install github.com/cosmtrek/air@latest
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

docker-up:
    docker-compose up -d

docker-down:
    docker-compose down

docker-logs:
    docker-compose logs -f

migrate-up:
    go run cmd/server/main.go migrate up

migrate-down:
    go run cmd/server/main.go migrate down

swagger:
    @which swag > /dev/null 2>&1 || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
    @if command -v swag > /dev/null 2>&1; then \
        swag init -g cmd/server/main.go -o docs; \
    else \
        $(go env GOPATH)/bin/swag init -g cmd/server/main.go -o docs || ~/go/bin/swag init -g cmd/server/main.go -o docs; \
    fi

fmt:
    go fmt ./...

lint:
    golangci-lint run

clean:
    rm -rf bin/
    rm -rf coverage.out
    rm -rf docs/

frontend-install:
    cd frontend && npm install

frontend-dev:
    cd frontend && npm run dev

frontend-build:
    cd frontend && npm run build

frontend-start:
    cd frontend && npm start

dev-all:
    @echo "Starting backend and frontend..."
    just docker-up
    sleep 2
    just run &
    just frontend-dev

help:
    @just --list
