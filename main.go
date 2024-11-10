package main

import (
	"html/template"
	"log"
	"net/http"
)

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	indexTemplate, err := template.ParseFiles("static/index.html")
	if err != nil {
		http.Error(w, "unable to load template", http.StatusInternalServerError)
		return
	}

	indexTemplate.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", indexHandler)

	err := http.ListenAndServe(":3002", nil)
	if err != nil {
		log.Fatal("server failed to start:", err)
	}
}
