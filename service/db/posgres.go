package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

const (
	host   = "localhost"
	port   = 5432
	dbname = "link_shortener"
)

func GetDBCredentials() string {
	user := os.Getenv("db_user")
	password := os.Getenv("db_password")
	credentials := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", user, password, host, port, dbname)

	return credentials
}

func InsertIntoDatabase(db *sql.DB, orignal_url, short_code string, create_date, expires_at time.Time) error {
	query := `INSERT INTO urls(original_url, short_code, create_date, expires_at)
	VALUES ($1, $2, $3, $4);`

	_, err := db.Exec(query, orignal_url, short_code, create_date, expires_at)
	return err
}

func GetLinkByHashFromDb(db *sql.DB, hash string) (*string, error) {
	query := `SELECT original_url FROM urls WHERE short_code=$1`
	var originalURL string
	err := db.QueryRow(query, hash).Scan(&originalURL)
	if err != nil {
		return nil, err
	}
	return &originalURL, nil
}
