package expense

import "context"

// ExpenseRepository defines the interface for expense data operations.
type ExpenseRepository interface {
	CreateExpense(ctx context.Context, expense *Expense) error
	GetExpense(ctx context.Context, householdID uint, year int, month int, category *string) ([]*Expense, error)
	UpdateExpense(ctx context.Context, expense *Expense) error
	DeleteExpense(ctx context.Context, expenseId ExpenseID) error
}
