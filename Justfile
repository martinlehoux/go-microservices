database_url := "postgres://user:password@localhost:5432/authentication?sslmode=disable"

all: test lint build

test:
  go test ./...

lint:
  staticcheck ./...
  errcheck -ignoretests ./...
  go vet ./...

build:
  go build

migrate DIRECTION:
  migrate -path=./authentication/migrations -database={{database_url}} {{DIRECTION}}
