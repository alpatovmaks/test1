
package main

import (
	"log"
	"net/http"

	"expenses/handlers"
)

func main() {
	// Статика
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Роуты
	http.HandleFunc("/", handlers.ListHandler)
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.AddHandler(w, r)
			return
		}
		handlers.AddFormHandler(w, r)
	})

	log.Println("Сервер запущен: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}