package main

import (
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	var err error
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"))

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

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		type requestBody struct {
			URL string `json:"url"`
		}

		var req requestBody
		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Try to insert with incrementing counter until we succeed
		counter := 0
		var slug string
		for {
			slug = generateSlug(req.URL, counter)

			// Try to insert the URL
			_, err = db.Exec("insert into urls (slug, original) values ($1, $2)", slug, req.URL)
			if err != nil {
				// Check if it's a unique constraint violation
				if isUniqueViolation(err) {
					counter++
					continue // Try again with incremented counter
				}
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			break // If we get here, insert succeeded
		}

		w.Write([]byte(slug))
	}).Methods(http.MethodPost)

	r.HandleFunc("/{slug}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		slug := mux.Vars(r)["slug"]
		var og string
		err = db.QueryRow("select original from urls where slug = $1", slug).Scan(&og)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("update urls set clicked = clicked + 1 where slug = $1", slug)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, og, http.StatusFound)
	}).Methods(http.MethodGet)

	server := &http.Server{Addr: ":3002", Handler: r}
	log.Printf("server started listening on %s", server.Addr)
	server.ListenAndServe()
}

func generateSlug(originalUrl string, counter int) string {
	// If counter is 0, just use the original hash
	if counter == 0 {
		hash := sha1.Sum([]byte(originalUrl))
		return base64.RawURLEncoding.EncodeToString(hash[:])[:7]
	}

	// If counter > 0, append it to the URL before hashing
	urlWithCounter := originalUrl + "#" + strconv.Itoa(counter)
	hash := sha1.Sum([]byte(urlWithCounter))
	return base64.RawURLEncoding.EncodeToString(hash[:])[:7]
}

// Helper function to check if an error is a unique constraint violation
func isUniqueViolation(err error) bool {
	// This checks for Postgres unique violation error code
	return err.Error() == `pq: duplicate key value violates unique constraint "urls_pkey"`
}
