package authentication

import (
	"database/sql"
	"errors"
	"go-microservices/common"

	_ "github.com/mattn/go-sqlite3"
)

type SqlAccountStore struct {
	db *sql.DB
}

func NewSqlAccountStore() SqlAccountStore {
	db, err := sql.Open("sqlite3", "test.db")
	common.PanicOnError(err)
	return SqlAccountStore{
		db: db,
	}
}

func (store *SqlAccountStore) Save(account Account) error {
	query, err := store.db.Prepare("INSERT INTO accounts (id, identifier, hashed_password) VALUES (?, ?, ?)")
	if err != nil {
		return errors.New("failed to prepare insert statement")
	}
	_, err = query.Exec(account.Id, account.Identifier, account.HashedPassword)
	if err != nil {
		return errors.New("failed to execute insert statement")
	}
	return nil

}

func (store *SqlAccountStore) LoadForIdentifier(identifier string) (Account, error) {
	var account = Account{}
	query, err := store.db.Prepare("SELECT id, identifier, hashed_password FROM accounts WHERE identifier = ?")
	if err != nil {
		return account, errors.New("failed to prepare query statement")
	}
	rows, err := query.Query(identifier)
	if err != nil {
		return account, errors.New("failed to query")
	}
	if err != nil {
		return account, errors.New("failed to query rows")
	}
	if !rows.Next() {
		return account, errors.New("no row = not found")
	}
	err = rows.Scan(&account.Id, &account.Identifier, &account.HashedPassword)
	if err != nil {
		return account, errors.New("failed to parse row")
	}
	err = rows.Err()
	if err != nil {
		return account, errors.New("error occured")
	}
	err = rows.Close()
	if err != nil {
		return account, errors.New("error occured")
	}
	return account, nil
}
