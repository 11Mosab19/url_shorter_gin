package database

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	DB *sql.DB
}

func ReadConfig() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	return "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + name + "?sslmode=disable", nil
}

func (db *Database) ConnectToDatabase() error {
	stringConnection, _ := ReadConfig()
	connection, err := sql.Open("pgx", stringConnection)
	if err != nil {
		return err
	}
	db.DB = connection
	if err := db.DB.Ping(); err != nil {
		return err
	}
	return nil
}

func (db *Database) CloseConnection() error {
	err := db.DB.Close()
	return err
}
