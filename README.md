# go-microservices

## Dependencies

- **github.com/google/uuid** for generating unique IDs
- **github.com/jackc/pgx/v4** for connection to the database Postgres
- **go.mongodb.org/mongo-driver** for connection to the database MongoDB
- **golang.org/x/crypto** for password hashing
- **github.com/golang-migrate/migrate** for managing database migrations
- **github.com/go-chi/chi** for http routing
- **github.com/sirupsen/logrus** for logging

## Dev Dependencies

- **github.com/stretchr/testify** for easier assertion

## Maintenance

- `TODO`
- `Deprecated`

## Domain Driven Design & Clean Architecture

### Stores

- A store manages persitence for aggregates
- It can be cleared using `store.Clear()`
- All methods should accept a `ctx context.Context` to handle context closure
- `GetByXXX` methods return a single result

### Errors

- An error has a text in lower case

### Dto Validation

- A function `validateXxxDto` may return errors
