package expense

// IExpenseRepository defines the interface for expense data operations.
type IExpenseRepository interface {
	CreateExpense(expense *Expense) error
	GetExpense(householdID uint, year int, month int, category *string) ([]*Expense, error)
	UpdateExpense(expense *Expense, expenseId uint) error
	DeleteExpense(expenseId uint) error
}
