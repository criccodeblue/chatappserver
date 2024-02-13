package database

import (
	"chatappserver/internal/model"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewStorage() *PostgresStorage {
	var storage = PostgresStorage{}
	postgresConnectionString := os.Getenv("PG_URL")

	db, err := sql.Open("postgres", postgresConnectionString)
	if err != nil {
		log.Fatalln("Error opening connection ", err.Error())
		return nil
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln("Error pinging connection ", err.Error())
		return nil
	}

	log.Print("Connection Established")
	storage.db = db
	return &storage
}

func (ps *PostgresStorage) CloseDBConnection() {
	err := ps.db.Close()
	if err != nil {
		log.Fatalln("Error closing connection ", err.Error())
	}
}

func (ps *PostgresStorage) GetUsers() ([]model.User, error) {
	query := "SELECT id, name, email, handle FROM users"

	rows, err := ps.db.Query(query)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatalln(err.Error())
		}
	}(rows)

	if err != nil {
		log.Fatalln("Error executing query")
		return nil, err
	}

	var users []model.User

	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Handle); err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	return users, nil
}

func (ps *PostgresStorage) GetUserByHandle(handle string) (model.User, error) {
	query := "SELECT id, name, email, handle FROM users WHERE handle=$1"

	row := ps.db.QueryRow(query, handle)
	var user model.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Handle); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (ps *PostgresStorage) CreateUser(authUser *model.AuthUser) (model.User, error) {
	query := "INSERT INTO users (name, email, handle, password) VALUES ($1, $2, $3, $4) RETURNING id, name, email, handle"

	row := ps.db.QueryRow(query, authUser.Name, authUser.Email, authUser.Handle, authUser.PasswordHash)
	var user model.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Handle); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (ps *PostgresStorage) GetUserPassword(handle string) ([]byte, error) {
	var password []byte
	query := "SELECT users.password FROM users WHERE handle = $1"

	row := ps.db.QueryRow(query, handle)
	if err := row.Scan(&password); err != nil {
		return []byte{}, err
	}

	return password, nil
}
