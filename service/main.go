package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Link struct {
	URL string `json:"URL"`
}

func main() {
	http.HandleFunc("/v1/create-link", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		link := &Link{}
		err := json.NewDecoder(r.Body).Decode(link)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println("create link", link)
		w.WriteHeader(http.StatusCreated)
	})

	fmt.Println("Running the server")
	if err := http.ListenAndServe(":8080", nil); err != http.ErrServerClosed {
		panic(err)
	}
}
