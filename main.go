package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID        int
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
		insertTodo(db, Todo{0, title, note, false})
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	h3 := func(w http.ResponseWriter, r *http.Request) {
		id := getIDbyURL(r.URL.String())
		completeTodo(db, id)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/add-todo/", h2)
	http.HandleFunc("/complete-todo/", h3)

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
	query := `INSERT INTO todo (title, note)
			VALUES ($1, $2) RETURNING id`

	var pk int
	err := db.QueryRow(query, todo.Title, todo.Note).Scan(&pk)

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Created: %d - Title: %s", pk, todo.Title)
}

func getTodoById(db *sql.DB, pk int) Todo {
	var id int
	var title string
	var note string
	var completed bool

	query := `SELECT * FROM todo WHERE id = $1`
	err := db.QueryRow(query, pk).Scan(&id, &title, &note, &completed)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatalf("No records found for ID: %d", pk)
		}
		log.Fatal(err)
	}

	return Todo{id, title, note, completed}
}

func getAllTodos(db *sql.DB) []Todo {
	data := []Todo{}
	rows, err := db.Query("SELECT id, title, note, completed FROM todo")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var id int
	var title string
	var note string
	var completed bool

	for rows.Next() {
		err := rows.Scan(&id, &title, &note, &completed)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, Todo{id, title, note, completed})
	}

	return data
}

func completeTodo(db *sql.DB, id string) error {
	query := "UPDATE todo SET completed = 't' , updated = CURRENT_TIMESTAMP WHERE id = $1"
	_, err := db.Exec(query, id)
	log.Printf("Completed: %s", id)
	return err
}

func getIDbyURL(url string) string {
	urlStr := strings.TrimSuffix(url, "/")
	parts := strings.Split(urlStr, "/")
	pk := parts[len(parts)-1]

	return pk
}
