
package main

import (
	"log"
	"net/http"

	"expenses/handlers"
	"expenses/repository"
	"expenses/storage"
)

func main() {
	db, err := storage.Open("expenses.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewExpenseRepository(db)
	h := handlers.New(repo)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", h.ListHandler)
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.AddHandler(w, r)
			return
		}
		h.AddFormHandler(w, r)
	})
	http.HandleFunc("/edit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.EditHandler(w, r)
			return
		}
		h.EditFormHandler(w, r)
	})
	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.DeleteHandler(w, r)
			return
		}
		h.DeleteConfirmHandler(w, r)
	})

	log.Println("Сервер запущен: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}