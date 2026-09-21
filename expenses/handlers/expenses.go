package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"expenses/models"
)

var listTmpl = template.Must(template.ParseFiles(
	"templates/layout.html",
	"templates/list.html",
))

var addTmpl = template.Must(template.ParseFiles(
	"templates/layout.html",
	"templates/add.html",
))

func ListHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	var total float64
	for _, e := range models.Expenses {
		total += e.Amount
	}

	data := map[string]any{
		"Title":    "Список трат",
		"Expenses": models.Expenses,
		"Total":    total,
	}

	if err := listTmpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AddFormHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title": "Добавить трату",
	}
	if err := addTmpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		renderAddError(w, "Не удалось разобрать форму")
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil || amount <= 0 {
		renderAddError(w, "Некорректная сумма")
		return
	}

	date, err := time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		renderAddError(w, "Некорректная дата")
		return
	}

	models.Expenses = append(models.Expenses, models.Expense{
		Amount:      amount,
		Description: r.FormValue("description"),
		Date:        date,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func renderAddError(w http.ResponseWriter, msg string) {
	data := map[string]any{
		"Title": "Добавить трату",
		"Error": msg,
	}
	addTmpl.ExecuteTemplate(w, "layout", data)
}