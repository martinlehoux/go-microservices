database_base_url := "postgres://user:password@localhost:5432"

all: test lint build

test: spec intg

spec:
  go test -tags spec ./...

intg:
  just db-create test_authentication
  just db-create test_human_resources
  just db-migrate authentication test_authentication up
  just db-migrate human_resources test_human_resources up
  go test -tags intg -p 1 ./...
  just db-delete test_authentication
  just db-delete test_human_resources

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

db-migrate SERVICE DATABASE DIRECTION:
  migrate -path=./{{SERVICE}}/migrations -database="{{database_base_url}}/{{DATABASE}}?sslmode=disable" {{DIRECTION}}

db-create DATABASE:
  PGPASSWORD=password psql -h localhost -U user -w -d postgres -c "CREATE DATABASE {{DATABASE}};"

db-delete DATABASE:
  PGPASSWORD=password psql -h localhost -U user -w -d postgres -c "DROP DATABASE {{DATABASE}};"

database SERVICE:
  just db-create {{SERVICE}}
  just db-migrate {{SERVICE}} {{SERVICE}} up

generate-migration SERVICE MIGRATION:
  migrate create -ext sql -dir {{SERVICE}}/migrations -seq {{MIGRATION}}

query-database QUERY:
  PGPASSWORD=password psql -h localhost -U user -w -d postgres -c "{{QUERY}}"