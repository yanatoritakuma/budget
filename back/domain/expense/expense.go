package expense

import (
	"time"

	"github.com/yanatoritakuma/budget/back/domain/user"
)

type Expense struct {
	ID        ExpenseID
	Amount    Amount
	StoreName StoreName
	Date      time.Time
	Category  Category
	Memo      Memo
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    UserID
	PayerID   PayerID
}

// NewExpense creates a new Expense domain entity.
func NewExpense(amount int, storeName string, date time.Time, category string, memo string, userID uint, payerID uint) (*Expense, error) {
	voAmount, err := NewAmount(amount)
	if err != nil {
		return nil, err
	}

	voStoreName, err := NewStoreName(storeName)
	if err != nil {
		return nil, err
	}

	voCategory, err := NewCategory(category)
	if err != nil {
		return nil, err
	}

	voMemo, err := NewMemo(memo)
	if err != nil {
		return nil, err
	}

	return &Expense{
		Amount:    voAmount,
		StoreName: voStoreName,
		Date:      date,
		Category:  voCategory,
		Memo:      voMemo,
		UserID:    UserID(user.UserID(userID)),
		PayerID:   PayerID(user.UserID(payerID)),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
