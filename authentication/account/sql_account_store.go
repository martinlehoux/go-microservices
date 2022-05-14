package account

import (
	"context"
	"go-microservices/common"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SqlAccountStore struct {
	conn *pgxpool.Pool
}

func NewSqlAccountStore() SqlAccountStore {
	conn, err := pgxpool.Connect(context.Background(), "postgres://user:password@localhost:5432/authentication")
	common.PanicOnError(err)
	return SqlAccountStore{
		conn: conn,
	}
}

func (store *SqlAccountStore) Save(ctx context.Context, account Account) error {
	_, err := store.conn.Exec(ctx, "INSERT INTO accounts (id, identifier, hashed_password) VALUES ($1, $2, $3)", account.id, account.identifier, account.hashedPassword)
	return err

}

func (store *SqlAccountStore) GetByIdentifier(ctx context.Context, identifier string) (Account, error) {
	var account Account

	err := store.conn.QueryRow(ctx, "SELECT id, identifier, hashed_password FROM accounts WHERE identifier = $1", identifier).Scan(&account.id, &account.identifier, &account.hashedPassword)

	return account, convertPgxError(err)
}

func (store *SqlAccountStore) Clear() {
	_, err := store.conn.Exec(context.Background(), "DELETE FROM accounts")
	common.PanicOnError(err)
}

func convertPgxError(err error) error {
	if err == pgx.ErrNoRows {
		return ErrAccountNotFound
	}
	return err
}
