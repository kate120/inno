package main

import (
	"html/template"
	"net/http"
)

func main() {
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
		Message: "Hello, Kate",
	}

	tmpl.Execute(w, data)
}
