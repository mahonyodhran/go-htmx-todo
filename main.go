package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Todo struct {
	Title     string
	Note      string
	Completed bool
}

func main() {
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

	createTodoTable(db)

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		todos := getAllTodos(db)

		tmpl.Execute(w, todos)
	}

	h2 := func(w http.ResponseWriter, r *http.Request) {
		title := r.PostFormValue("title")
		note := r.PostFormValue("note")
		todo := Todo{title, note, false}

		insertTodo(db, todo)

		todos := getAllTodos(db)

		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, todos)
	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/add-todo/", h2)

	log.Println("App starting...")
	log.Fatal(http.ListenAndServe(":8008", nil))
}

func createTodoTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS todo (
		id SERIAL PRIMARY KEY,
		title VARCHAR(50) NOT NULL,
		note VARCHAR(200) NOT NULL,
		created TIMESTAMP DEFAULT NOW(),
		updated TIMESTAMP DEFAULT NOW(),
		completed BOOLEAN DEFAULT FALSE
	)`

	_, err := db.Exec(query)

	if err != nil {
		log.Fatal(err)
	}
}

func insertTodo(db *sql.DB, todo Todo) {
	fmt.Println("INSIDE INSSERTTODO")
	query := `INSERT INTO todo (title, note)
			VALUES ($1, $2) RETURNING id`

	var pk int
	err := db.QueryRow(query, todo.Title, todo.Note).Scan(&pk)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully added Todo")
	fmt.Printf("ID: %d\n", pk)
	fmt.Printf("Title: %s\n", todo.Title)
	fmt.Printf("Title: %s\n", todo.Note)
}

func getTodoById(db *sql.DB, pk int) Todo {
	var title string
	var note string
	var completed bool

	query := `SELECT * FROM todo WHERE id = $1`
	err := db.QueryRow(query, pk).Scan(&title, &note, &completed)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatalf("No records found for ID: %d", pk)
		}
		log.Fatal(err)
	}

	return Todo{title, note, completed}
}

func getAllTodos(db *sql.DB) []Todo {
	data := []Todo{}
	rows, err := db.Query("SELECT title, note, completed FROM todo")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var title string
	var note string
	var completed bool

	for rows.Next() {
		err := rows.Scan(&title, &note, &completed)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, Todo{title, note, completed})
	}

	return data
}
