test:
  go test go-microservices/authentication
  go test go-microservices/user

lint:
  staticcheck ./...
  errcheck -ignoretests ./...