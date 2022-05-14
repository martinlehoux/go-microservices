# go-microservices

## Dependencies

- **github.com/google/uuid** for generating unique IDs
- **github.com/jackc/pgx/v4** for connection to the database Postgres
- **go.mongodb.org/mongo-driver** for connection to the database MongoDB
- **golang.org/x/crypto** for password hashing
- **github.com/golang-migrate/migrate** for managing database migrations

## Dev Dependencies

- **github.com/stretchr/testify** for easier assertion

## Domain Driven Design & Clean Architecture

### Stores

- A store manages persitence for aggregates
- It can be cleared using `store.Clear()`
- All methods should accept a `ctx context.Context` to handle context closure
- `GetByXXX` methods return a single result
