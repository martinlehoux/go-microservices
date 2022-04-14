database_base_url := "postgres://user:password@localhost:5432"

all: test lint build

test:
  go test -p 1 ./...

lint:
  staticcheck ./...
  errcheck -ignoretests ./...
  go vet ./...

build:
  go build

migrate SERVICE DIRECTION:
  migrate -path=./{{SERVICE}}/migrations -database="{{database_base_url}}/{{SERVICE}}?sslmode=disable" {{DIRECTION}}

create-database SERVICE:
  PGPASSWORD=password psql -h localhost -U user -w -d postgres -c "CREATE DATABASE {{SERVICE}};"
