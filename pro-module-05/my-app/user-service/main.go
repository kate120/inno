package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Name string `json:"name"`
}

var db *sql.DB

func initDB() {
	var err error
	// Получаем строку подключения из переменной окружения
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(mysql-users:3306)/users_db?charset=utf8&parseTime=True&loc=Local"
	}

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	// Создаём таблицу, если её нет
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("Ошибка создания таблицы:", err)
	}

	log.Println("Подключение к MySQL (users) установлено")
}

func saveUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Only POST allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if user.Name == "" {
		http.Error(w, `{"error":"Name is required"}`, http.StatusBadRequest)
		return
	}

	// Сохраняем в БД
	result, err := db.Exec("INSERT INTO users (name) VALUES (?)", user.Name)
	if err != nil {
		log.Println("Ошибка вставки:", err)
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":   id,
		"name": user.Name,
	})
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/api/users", saveUserHandler)

	log.Println("User Service запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
