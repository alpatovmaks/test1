package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Главная страница — приветствие
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "Добро пожаловать на главную страницу!")
	})

	// /about — описание проекта
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Это простой HTTP-сервер на Go с тремя маршрутами.")
	})

	// /ping — возвращает 200 OK и слово pong
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "pong")
	})

	fmt.Println("Сервер запущен на http ://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
	}
}