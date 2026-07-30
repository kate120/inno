package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type SaveRequest struct {
	User string `json:"user"`
	Book string `json:"book"`
}

var userServiceURL = getEnv("USER_SERVICE_URL", "http://user-service/api/users")
var bookServiceURL = getEnv("BOOK_SERVICE_URL", "http://book-service/api/books")

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// 🔥 НОВЫЙ ОБРАБОТЧИК: отдает HTML-страницу
func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("static", "index.html"))
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Только POST запросы", http.StatusMethodNotAllowed)
		return
	}

	var req SaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if err := sendToService(userServiceURL, map[string]string{"name": req.User}); err != nil {
		log.Println("Ошибка сохранения пользователя:", err)
		http.Error(w, "Не удалось сохранить пользователя", http.StatusInternalServerError)
		return
	}

	if err := sendToService(bookServiceURL, map[string]string{"title": req.Book}); err != nil {
		log.Println("Ошибка сохранения книги:", err)
		http.Error(w, "Не удалось сохранить книгу", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func sendToService(url string, payload map[string]string) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func main() {
	// ✅ Отдаем HTML-страницу
	http.HandleFunc("/", serveIndex)

	// ✅ Отдаем API
	http.HandleFunc("/api/save", saveHandler)

	log.Println("Backend запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
