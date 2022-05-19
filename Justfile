database_base_url := "postgres://user:password@localhost:5432"

all: test lint build

test: spec intg

spec:
  go test -tags spec ./...

intg:
  go test -tags intg -p 1 ./...

coverage:
  go test -coverprofile=coverage.out -tags="spec intg" ./...
  go tool cover -html=coverage.out

lint:
  go fmt ./...
  staticcheck ./...
  errcheck -ignoretests -asserts -exclude rules/errcheck  ./...
  go vet -composites=false ./...

build:
  go build

migrate SERVICE DIRECTION:
  migrate -path=./{{SERVICE}}/migrations -database="{{database_base_url}}/{{SERVICE}}?sslmode=disable" {{DIRECTION}}

create-database SERVICE:
  PGPASSWORD=password psql -h localhost -U user -w -d postgres -c "CREATE DATABASE {{SERVICE}};"

database SERVICE:
  just create-database {{SERVICE}}
  just migrate {{SERVICE}} up

generate-migration SERVICE MIGRATION:
  migrate create -ext sql -dir {{SERVICE}}/migrations -seq {{MIGRATION}}

query-database QUERY:
  PGPASSWORD=password psql -h localhost -U user -w -d postgres -c "{{QUERY}}"