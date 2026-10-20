package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"expenses/models"
	"expenses/repository"
)

var listTmpl = template.Must(template.ParseFiles(
	"templates/layout.html",
	"templates/list.html",
))

var formTmpl = template.Must(template.ParseFiles(
	"templates/layout.html",
	"templates/form.html",
))

var deleteTmpl = template.Must(template.ParseFiles(
	"templates/layout.html",
	"templates/delete.html",
))

type Handler struct {
	repo *repository.ExpenseRepository
}

func New(repo *repository.ExpenseRepository) *Handler {
	return &Handler{repo: repo}
}

// ListHandler — GET /
func (h *Handler) ListHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	expenses, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var total float64
	for _, e := range expenses {
		total += e.Amount
	}
	data := map[string]any{
		"Title":    "Список трат",
		"Expenses": expenses,
		"Total":    total,
	}
	render(w, listTmpl, data)
}

// AddFormHandler — GET /add
func (h *Handler) AddFormHandler(w http.ResponseWriter, r *http.Request) {
	render(w, formTmpl, map[string]any{
		"Title": "Добавить трату",
	})
}

// AddHandler — POST /add
func (h *Handler) AddHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderFormError(w, nil, "Не удалось разобрать форму")
		return
	}
	e, msg := parseExpenseForm(r)
	if msg != "" {
		h.renderFormError(w, e, msg)
		return
	}
	if _, err := h.repo.Create(*e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// EditFormHandler — GET /edit?id=1
func (h *Handler) EditFormHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	e, err := h.repo.GetByID(id)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render(w, formTmpl, map[string]any{
		"Title":   "Редактировать трату",
		"Expense": e,
		"Edit":    true,
	})
}

// EditHandler — POST /edit?id=1
func (h *Handler) EditHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	e, msg := parseExpenseForm(r)
	if msg != "" {
		h.renderFormError(w, e, msg)
		return
	}
	e.ID = id
	if err := h.repo.Update(*e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// DeleteConfirmHandler — GET /delete?id=1
func (h *Handler) DeleteConfirmHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	e, err := h.repo.GetByID(id)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render(w, deleteTmpl, map[string]any{
		"Title":   "Удалить трату",
		"Expense": e,
	})
}

// DeleteHandler — POST /delete?id=1
func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if err := h.repo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// --- helpers ---

func parseExpenseForm(r *http.Request) (*models.Expense, string) {
	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil || amount <= 0 {
		return &models.Expense{
			Amount:      amount,
			Description: r.FormValue("description"),
		}, "Некорректная сумма"
	}
	date, err := time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		return &models.Expense{
			Amount:      amount,
			Description: r.FormValue("description"),
		}, "Некорректная дата"
	}
	return &models.Expense{
		Amount:      amount,
		Description: r.FormValue("description"),
		Date:        date,
	}, ""
}

func (h *Handler) renderFormError(w http.ResponseWriter, e *models.Expense, msg string) {
	data := map[string]any{
		"Title": "Добавить трату",
		"Error": msg,
	}
	if e != nil {
		data["Expense"] = e
	}
	render(w, formTmpl, data)
}

func render(w http.ResponseWriter, tmpl *template.Template, data any) {
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}