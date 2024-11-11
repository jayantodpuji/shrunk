package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	var err error
	connStr := "user=postgres dbname=shrunk host=localhost port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("database connection failed", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error closing the database connection: %v", err)
		}
	}()

	if err = db.Ping(); err != nil {
		log.Fatalf("unable to reach the database: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		type requestBody struct {
			URL string `json:"url"`
		}

		var req requestBody
		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	})

	mux.HandleFunc("/:slug", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})

	server := &http.Server{Addr: ":3002", Handler: mux}
	server.ListenAndServe()
}
