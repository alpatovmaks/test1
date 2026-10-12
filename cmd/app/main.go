package main

import (
	"log"
	"net/http"

	"expenses-app/internal/handlers"
	"expenses-app/internal/repository"
)

func main() {
	db, err := repository.NewDB("expenses.db")
	if err != nil {
		log.Fatalf("db init: %v", err)
	}
	defer db.Close()

	repo := repository.NewExpenseRepository(db)

	h, err := handlers.New(repo)
	if err != nil {
		log.Fatalf("handlers init: %v", err)
	}

	addr := ":8080"
	log.Printf("server started on %s", addr)
	if err := http.ListenAndServe(addr, h.Routes()); err != nil {
		log.Fatalf("server: %v", err)
	}
}