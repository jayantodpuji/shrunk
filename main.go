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
	"sync"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Global connection pool
var (
	db   *sql.DB
	pool sync.Pool
)

func init() {
	// Initialize a pool of reusable request body objects
	pool.New = func() interface{} {
		return &requestBody{}
	}
}

type requestBody struct {
	URL string `json:"url"`
}

func main() {
	var err error
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"))

	// Configure connection pool
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("database connection failed", err)
	}

	// Set connection pool parameters
	db.SetMaxOpenConns(25)                 // Limit max open connections
	db.SetMaxIdleConns(25)                 // Keep connections ready
	db.SetConnMaxLifetime(5 * time.Minute) // Recycle connections
	db.SetConnMaxIdleTime(5 * time.Minute) // Close idle connections

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error closing the database connection: %v", err)
		}
	}()

	if err = db.Ping(); err != nil {
		log.Fatalf("unable to reach the database: %v", err)
	}

	// Create prepared statements
	insertStmt, err := db.Prepare("INSERT INTO urls (slug, original) VALUES ($1, $2)")
	if err != nil {
		log.Fatal("failed to prepare insert statement:", err)
	}
	defer insertStmt.Close()

	selectStmt, err := db.Prepare("SELECT original FROM urls WHERE slug = $1")
	if err != nil {
		log.Fatal("failed to prepare select statement:", err)
	}
	defer selectStmt.Close()

	updateStmt, err := db.Prepare("UPDATE urls SET clicked = clicked + 1 WHERE slug = $1")
	if err != nil {
		log.Fatal("failed to prepare update statement:", err)
	}
	defer updateStmt.Close()

	r := mux.NewRouter()

	// POST handler for creating short URLs
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Get request body object from pool
		req := pool.Get().(*requestBody)
		defer pool.Put(req)

		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Try to insert with incrementing counter until we succeed
		counter := 0
		var slug string
		for i := 0; i < 3; i++ { // Limit retries
			slug = generateSlug(req.URL, counter)

			// Try to insert the URL using prepared statement
			_, err = insertStmt.Exec(slug, req.URL)
			if err != nil {
				if isUniqueViolation(err) {
					counter++
					continue
				}
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			break
		}

		w.Write([]byte(slug))
	}).Methods(http.MethodPost)

	// GET handler for redirecting short URLs
	r.HandleFunc("/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := mux.Vars(r)["slug"]

		var og string
		err = selectStmt.QueryRow(slug).Scan(&og)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Update click count asynchronously
		go func() {
			_, err := updateStmt.Exec(slug)
			if err != nil {
				log.Printf("error updating click count: %v", err)
			}
		}()

		http.Redirect(w, r, og, http.StatusFound)
	}).Methods(http.MethodGet)

	// Configure server timeouts
	server := &http.Server{
		Addr:           ":3002",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	log.Printf("server started listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func generateSlug(originalUrl string, counter int) string {
	if counter == 0 {
		hash := sha1.Sum([]byte(originalUrl))
		return base64.RawURLEncoding.EncodeToString(hash[:])[:7]
	}

	urlWithCounter := originalUrl + "#" + strconv.Itoa(counter)
	hash := sha1.Sum([]byte(urlWithCounter))
	return base64.RawURLEncoding.EncodeToString(hash[:])[:7]
}

func isUniqueViolation(err error) bool {
	return err.Error() == `pq: duplicate key value violates unique constraint "urls_pkey"`
}
