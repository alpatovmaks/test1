package repository

import (
	"database/sql"
	"time"

	"expenses/models"
)

type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) Create(e models.Expense) (int64, error) {
	res, err := r.db.Exec(
		`INSERT INTO expenses (amount, description, date) VALUES (?, ?, ?)`,
		e.Amount, e.Description, e.Date.Format("2006-01-02"),
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *ExpenseRepository) GetAll() ([]models.Expense, error) {
	rows, err := r.db.Query(
		`SELECT id, amount, description, date FROM expenses ORDER BY date DESC, id DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Expense
	for rows.Next() {
		var (
			e    models.Expense
			date string
		)
		if err := rows.Scan(&e.ID, &e.Amount, &e.Description, &date); err != nil {
			return nil, err
		}
		e.Date, err = time.Parse("2006-01-02", date)
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, rows.Err()
}

func (r *ExpenseRepository) GetByID(id int64) (models.Expense, error) {
	var (
		e    models.Expense
		date string
	)
	err := r.db.QueryRow(
		`SELECT id, amount, description, date FROM expenses WHERE id = ?`, id,
	).Scan(&e.ID, &e.Amount, &e.Description, &date)
	if err != nil {
		return models.Expense{}, err
	}
	e.Date, err = time.Parse("2006-01-02", date)
	return e, err
}

func (r *ExpenseRepository) Update(e models.Expense) error {
	_, err := r.db.Exec(
		`UPDATE expenses SET amount = ?, description = ?, date = ? WHERE id = ?`,
		e.Amount, e.Description, e.Date.Format("2006-01-02"), e.ID,
	)
	return err
}

func (r *ExpenseRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM expenses WHERE id = ?`, id)
	return err
}