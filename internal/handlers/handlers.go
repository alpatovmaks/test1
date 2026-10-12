package handlers

import (
	"html/template"
	"log"
	"net/http"

	"expenses-app/internal/repository"
)

// Handlers хранит общие зависимости всех HTTP-обработчиков.
type Handlers struct {
	repo      *repository.ExpenseRepository
	templates *template.Template
}

// New создаёт набор обработчиков и парсит шаблоны.
func New(repo *repository.ExpenseRepository) (*Handlers, error) {
	tmpl, err := template.ParseGlob("web/templates/*.html")
	if err != nil {
		return nil, err
	}
	return &Handlers{repo: repo, templates: tmpl}, nil
}

// render — маленький помощник для рендера шаблона с обработкой ошибок.
func (h *Handlers) render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("template %s: %v", name, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// Routes регистрирует все маршруты приложения.
func (h *Handlers) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", h.ListExpenses)
	mux.HandleFunc("GET /expenses", h.ListExpenses)
	mux.HandleFunc("GET /expenses/new", h.NewExpenseForm)
	mux.HandleFunc("POST /expenses", h.CreateExpense)
	mux.HandleFunc("GET /expenses/{id}/edit", h.EditExpenseForm)
	mux.HandleFunc("POST /expenses/{id}", h.UpdateExpense)
	mux.HandleFunc("POST /expenses/{id}/delete", h.DeleteExpense)

	return mux
}