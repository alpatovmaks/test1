package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"expenses-app/internal/models"
	"expenses-app/internal/repository"
)

// ListExpenses отображает список всех трат.
func (h *Handlers) ListExpenses(w http.ResponseWriter, r *http.Request) {
	expenses, err := h.repo.GetAll()
	if err != nil {
		log.Printf("list expenses: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	h.render(w, "expenses.html", map[string]any{
		"Title":    "Список трат",
		"Expenses": expenses,
	})
}

// NewExpenseForm показывает пустую форму добавления.
func (h *Handlers) NewExpenseForm(w http.ResponseWriter, r *http.Request) {
	h.render(w, "expense_form.html", map[string]any{
		"Title": "Новая трата",
		"Form":  models.ExpenseForm{Date: time.Now().Format("2006-01-02")},
	})
}

// CreateExpense обрабатывает POST формы добавления.
func (h *Handlers) CreateExpense(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	form := models.ExpenseForm{
		Amount:      strings.TrimSpace(r.FormValue("amount")),
		Description: strings.TrimSpace(r.FormValue("description")),
		Date:        strings.TrimSpace(r.FormValue("date")),
	}

	expense, errMsg := formToExpense(form)
	if errMsg != "" {
		w.WriteHeader(http.StatusBadRequest)
		h.render(w, "expense_form.html", map[string]any{
			"Title": "Новая трата",
			"Form":  form,
			"Error": errMsg,
		})
		return
	}

	if _, err := h.repo.Create(expense); err != nil {
		log.Printf("create expense: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	// Post/Redirect/Get
	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}

// EditExpenseForm показывает форму редактирования.
func (h *Handlers) EditExpenseForm(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	expense, err := h.repo.GetByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		log.Printf("get expense: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.render(w, "expense_edit.html", map[string]any{
		"Title":   "Редактирование траты",
		"ID":      expense.ID,
		"Form": models.ExpenseForm{
			Amount:      strconv.FormatFloat(expense.Amount, 'f', 2, 64),
			Description: expense.Description,
			Date:        expense.Date.Format("2006-01-02"),
		},
	})
}

// UpdateExpense обрабатывает POST формы редактирования.
func (h *Handlers) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	form := models.ExpenseForm{
		Amount:      strings.TrimSpace(r.FormValue("amount")),
		Description: strings.TrimSpace(r.FormValue("description")),
		Date:        strings.TrimSpace(r.FormValue("date")),
	}
	expense, errMsg := formToExpense(form)
	if errMsg != "" {
		w.WriteHeader(http.StatusBadRequest)
		h.render(w, "expense_edit.html", map[string]any{
			"Title": "Редактирование траты",
			"ID":    id,
			"Form":  form,
			"Error": errMsg,
		})
		return
	}
	expense.ID = id

	if err := h.repo.Update(expense); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		log.Printf("update expense: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}

// DeleteExpense удаляет трату по id.
func (h *Handlers) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if err := h.repo.Delete(id); err != nil && !errors.Is(err, repository.ErrNotFound) {
		log.Printf("delete expense: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}

// --- Вспомогательные функции ---

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}

// formToExpense валидирует и преобразует форму в модель.
// Возвращает либо заполненную модель, либо текст ошибки.
func formToExpense(f models.ExpenseForm) (*models.Expense, string) {
	if f.Amount == "" || f.Description == "" || f.Date == "" {
		return nil, "Все поля обязательны для заполнения."
	}
	amount, err := strconv.ParseFloat(f.Amount, 64)
	if err != nil {
		return nil, "Сумма должна быть числом."
	}
	if amount <= 0 {
		return nil, "Сумма должна быть положительной."
	}
	date, err := time.Parse("2006-01-02", f.Date)
	if err != nil {
		return nil, "Дата должна быть в формате ГГГГ-ММ-ДД."
	}
	return &models.Expense{
		Amount:      amount,
		Description: f.Description,
		Date:        date,
	}, ""
}