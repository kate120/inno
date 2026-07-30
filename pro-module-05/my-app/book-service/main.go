package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Book struct {
	Title string `json:"title"`
}

var db *sql.DB

func initDB() {
	var err error
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(mysql-books:3306)/books_db?charset=utf8&parseTime=True&loc=Local"
	}

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS books (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("Ошибка создания таблицы:", err)
	}

	log.Println("Подключение к MySQL (books) установлено")
}

func saveBookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Only POST allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if book.Title == "" {
		http.Error(w, `{"error":"Title is required"}`, http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO books (title) VALUES (?)", book.Title)
	if err != nil {
		log.Println("Ошибка вставки:", err)
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    id,
		"title": book.Title,
	})
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/api/books", saveBookHandler)

	log.Println("Book Service запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
