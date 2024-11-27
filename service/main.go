package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const (
	host   = "158.255.74.123"
	port   = 5432
	dbname = "link_shortener"
)

type Link struct {
	URL string `json:"URL"`
}

func getDBCredentials() string {
	user := os.Getenv("db_user")
	password := os.Getenv("db_password")
	credentials := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	return credentials
}

func insertIntoDatabase(db *sql.DB, orignal_url, short_code string, create_date, expires_at time.Time) error {
	query := `INSERT INTO urls(original_url, short_code, create_date, expires_at)
	VALUES ($1, $2, $3, $4);`

	_, err := db.Exec(query, orignal_url, short_code, create_date, expires_at)
	return err
}

type Handler struct {
	DB *sql.DB
}

func (h *Handler) createLinkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	link := &Link{}
	err := json.NewDecoder(r.Body).Decode(link)

	if err != nil {
		fmt.Println("error = ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("create link", link)
	err = insertIntoDatabase(h.DB, link.URL, "hello fatemeh", time.Now(), time.Now())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func main() {
	dbCredentials := getDBCredentials()
	db, err := sql.Open("postgres", dbCredentials)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	handler := &Handler{db}
	http.HandleFunc("/v1/create-link", handler.createLinkHandler)

	fmt.Println("Running the server")
	if err := http.ListenAndServe(":8080", nil); err != http.ErrServerClosed {
		panic(err)
	}
}
