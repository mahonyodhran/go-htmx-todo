package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

type Todo struct {
	Title     string
	Notes     string
	CreatedOn time.Time
	UpdatedOn time.Time
	Completed bool
}

func main() {
	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, nil)
	}

	http.HandleFunc("/", h1)

	log.Fatal(http.ListenAndServe(":8008", nil))
}
