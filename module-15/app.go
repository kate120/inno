package main

import (
	"html/template"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Добавляем эндпоинт /metrics для Prometheus
	http.Handle("/metrics", promhttp.Handler())

	// Ваш основной обработчик
	http.HandleFunc("/", homePage)

	http.ListenAndServe(":5000", nil)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("home_page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Title   string
		Message string
	}{
		Title:   "Welcome",
		Message: "Hello, Innowise! ",
	}

	tmpl.Execute(w, data)
}
