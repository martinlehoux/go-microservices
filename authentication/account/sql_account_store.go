package account

import (
	"context"
	"go-microservices/common"

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

func (store *SqlAccountStore) Save(account Account) error {
	_, err := store.conn.Exec(context.Background(), "INSERT INTO accounts (id, identifier, hashed_password) VALUES ($1, $2, $3)", account.id, account.identifier, account.hashedPassword)
	if err != nil {
		return err
	}
	return nil

}

func (store *SqlAccountStore) LoadForIdentifier(identifier string) (Account, error) {
	var account = Account{}

	err := store.conn.QueryRow(context.Background(), "SELECT id, identifier, hashed_password FROM accounts WHERE identifier = $1", identifier).Scan(&account.id, &account.identifier, &account.hashedPassword)

	return account, err
}

func (store *SqlAccountStore) truncate() error {
	_, err := store.conn.Exec(context.Background(), "DELETE FROM accounts")
	return err
}
