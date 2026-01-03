package expense

import "time"

type Expense struct {
	ID        uint
	Amount    int
	StoreName string
	Date      time.Time
	Category  string
	Memo      string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uint
	PayerID   uint
}

// NewExpense creates a new Expense domain entity.
func NewExpense(amount int, storeName string, date time.Time, category string, memo string, userID uint, payerID uint) (*Expense, error) {
	// TODO: Add validation and business rules here.
	return &Expense{
		Amount:    amount,
		StoreName: storeName,
		Date:      date,
		Category:  category,
		Memo:      memo,
		UserID:    userID,
		PayerID:   payerID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
