package http

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fatemehkarimi/link-shortener/service/db"
	"github.com/fatemehkarimi/link-shortener/service/internal"
)

type Handler struct {
	DB *sql.DB
}

var now = time.Now

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

	linkHash := internal.GenerateLinkHash(now, requestData.URL)
	create_at := now()
	expires_at := create_at.AddDate(0, 0, 7)

	err = db.InsertIntoDatabase(h.DB, requestData.URL, linkHash, create_at, expires_at)
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

	hash := r.URL.Query().Get("hash")
	originalURL, err := db.GetLinkByHashFromDb(h.DB, hash)
	if err != nil {
		http.Error(w, "url not found", http.StatusNotFound)
		return
	}

	if originalURL != nil {
		w.WriteHeader(http.StatusFound)
		response := ResponseGetURLByHash{
			URL: originalURL,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "not found", http.StatusNotFound)
}

func (h *Handler) Start() {
	dbCredentials := db.GetDBCredentials()
	db, err := sql.Open("postgres", dbCredentials)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	h.DB = db
	http.HandleFunc("/v1/create-link", h.createLinkHandler)
	http.HandleFunc("/v1/get-link-by-hash", h.getURLByHash)

	fmt.Println("Running the server")
	if err := http.ListenAndServe(":8080", nil); err != http.ErrServerClosed {
		panic(err)
	}
}

func NewHandler() *Handler {
	return &Handler{}
}
