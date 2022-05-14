package user

import (
	"context"
	"go-microservices/common"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SqlUserStore struct {
	conn *pgxpool.Pool
}

func NewSqlUserStore() SqlUserStore {
	conn, err := pgxpool.Connect(context.Background(), "postgres://user:password@localhost:5432/human_resources")
	common.PanicOnError(err)
	return SqlUserStore{
		conn: conn,
	}
}

func (store *SqlUserStore) Save(user User) error {
	_, err := store.conn.Exec(context.Background(), "INSERT INTO users (id, preferred_name, email) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET preferred_name = EXCLUDED.preferred_name", user.id, user.preferredName, user.email)
	if err != nil {
		return err
	}
	return nil
}

func (store *SqlUserStore) Get(userId UserID) (User, error) {
	var user User
	err := store.conn.QueryRow(context.Background(), "SELECT id, preferred_name, email FROM users WHERE id = $1", userId).Scan(&user.id, &user.preferredName, &user.email)
	return user, convertPgxError(err)
}

func (store *SqlUserStore) EmailExists(email string) (bool, error) {
	var count int
	err := store.conn.QueryRow(context.Background(), "SELECT count(*) FROM users WHERE email = $1", email).Scan(&count)
	return count > 0, err
}

func (store *SqlUserStore) GetMany() ([]User, error) {
	users := make([]User, 0)
	rows, err := store.conn.Query(context.Background(), "SELECT id, preferred_name, email FROM users")
	if err != nil {
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		err := rows.Scan(&user.id, &user.preferredName, &user.email)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (store *SqlUserStore) GetByEmail(email string) (User, error) {
	var user User
	err := store.conn.QueryRow(context.Background(), "SELECT id, preferred_name, email FROM users WHERE email = $1", email).Scan(&user.id, &user.preferredName, &user.email)
	return user, err
}

func convertPgxError(err error) error {
	if err == pgx.ErrNoRows {
		return ErrUserNotFound
	}
	return err
}
