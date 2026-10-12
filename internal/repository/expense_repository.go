
package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"expenses-app/internal/models"
)

var ErrNotFound = errors.New("expense not found")

type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) Create(e *models.Expense) (int64, error) {
	res, err := r.db.Exec(
		`INSERT INTO expenses (amount, description, date) VALUES (?, ?, ?)`,
		e.Amount, e.Description, e.Date.Format("2006-01-02"),
	)
	if err != nil {
		return 0, fmt.Errorf("insert expense: %w", err)
	}
	return res.LastInsertId()
}

func (r *ExpenseRepository) GetAll() ([]models.Expense, error) {
	rows, err := r.db.Query(
		`SELECT id, amount, description, date FROM expenses ORDER BY date DESC, id DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("query expenses: %w", err)
	}
	defer rows.Close()

	return scanExpenses(rows)
}

func (r *ExpenseRepository) GetByID(id int64) (*models.Expense, error) {
	row := r.db.QueryRow(
		`SELECT id, amount, description, date FROM expenses WHERE id = ?`, id,
	)
	var e models.Expense
	var dateStr string
	err := row.Scan(&e.ID, &e.Amount, &e.Description, &dateStr)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan expense: %w", err)
	}
	e.Date, err = time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("parse date: %w", err)
	}
	return &e, nil
}

func (r *ExpenseRepository) Update(e *models.Expense) error {
	res, err := r.db.Exec(
		`UPDATE expenses SET amount = ?, description = ?, date = ? WHERE id = ?`,
		e.Amount, e.Description, e.Date.Format("2006-01-02"), e.ID,
	)
	if err != nil {
		return fmt.Errorf("update expense: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *ExpenseRepository) Delete(id int64) error {
	res, err := r.db.Exec(`DELETE FROM expenses WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete expense: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func scanExpenses(rows *sql.Rows) ([]models.Expense, error) {
	var result []models.Expense
	for rows.Next() {
		var e models.Expense
		var dateStr string
		if err := rows.Scan(&e.ID, &e.Amount, &e.Description, &dateStr); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		d, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, fmt.Errorf("parse date: %w", err)
		}
		e.Date = d
		result = append(result, e)
	}
	return result, rows.Err()
}