package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Todo struct {
	Title     string
	Notes     string
	Completed bool
}

func main() {
	log.Println("App starting...")
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	connStr := os.Getenv("DB_CONN")
	db, err := sql.Open("postgres", connStr)

	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Checking for Todo db...")

	createTodoTable(db)

	log.Println("App running...")

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, nil)
	}

	http.HandleFunc("/", h1)

	log.Fatal(http.ListenAndServe(":8008", nil))
}

func createTodoTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS todo (
		id SERIAL PRIMARY KEY,
		title VARCHAR(50) NOT NULL,
		notes VARCHAR(200) NOT NULL,
		createdon TIMESTAMP DEFAULT NOW(),
		updatedon TIMESTAMP DEFAULT NOW(),
		completed BOOLEAN
	)`

	_, err := db.Exec(query)

	if err != nil {
		log.Fatal(err)
	}
}
