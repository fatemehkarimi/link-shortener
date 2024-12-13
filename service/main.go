package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var now = time.Now

const (
	host   = "localhost"
	port   = 5432
	dbname = "link_shortener"
)

func getDBCredentials() string {
	user := os.Getenv("db_user")
	password := os.Getenv("db_password")
	credentials := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", user, password, host, port, dbname)

	return credentials
}

func insertIntoDatabase(db *sql.DB, orignal_url, short_code string, create_date, expires_at time.Time) error {
	query := `INSERT INTO urls(original_url, short_code, create_date, expires_at)
	VALUES ($1, $2, $3, $4);`

	_, err := db.Exec(query, orignal_url, short_code, create_date, expires_at)
	return err
}

func getLinkByHashFromDb(db *sql.DB, hash string) (*string, error) {
	query := `SELECT original_url FROM urls WHERE short_code=$1`
	var originalURL string
	err := db.QueryRow(query, hash).Scan(&originalURL)
	if err != nil {
		return nil, err
	}
	return &originalURL, nil
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

	requestData := &RequestCreateLink{}
	err := json.NewDecoder(r.Body).Decode(requestData)

	if err != nil {
		fmt.Println("error = ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	linkHash := generateLinkHash(requestData.URL)
	create_at := now()
	expires_at := create_at.AddDate(0, 0, 7)

	err = insertIntoDatabase(h.DB, requestData.URL, linkHash, create_at, expires_at)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := ResponseCreateLink{
		Hash:       linkHash,
		Create_at:  create_at,
		Expires_at: expires_at,
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) getURLByHash(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	requestData := &RequestGetURLByHash{}
	err := json.NewDecoder(r.Body).Decode(requestData)

	if err != nil {
		fmt.Println("error = ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	originalURL, err := getLinkByHashFromDb(h.DB, requestData.Hash)
	if err != nil {
		http.Error(w, "url not found", http.StatusNotFound)
		return
	}

	if originalURL != nil {
		http.Redirect(w, r, *originalURL, http.StatusFound)
		return
	}

	http.Error(w, "not found", http.StatusNotFound)
}

func generateLinkHash(str string) string {
	strWithDate := str + now().String()
	hash := sha256.Sum256([]byte(strWithDate))
	result := fmt.Sprintf("%x", hash)
	return result[:12]
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
	http.HandleFunc("/v1/get-link-by-hash", handler.getURLByHash)

	fmt.Println("Running the server")
	if err := http.ListenAndServe(":8080", nil); err != http.ErrServerClosed {
		panic(err)
	}
}
