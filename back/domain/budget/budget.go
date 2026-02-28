package budget

import (
	"time"
)

type Budget struct {
	ID          uint
	HouseholdID uint
	YearMonth   string
	Amount      int
	NotifiedAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewBudget(householdID uint, yearMonth string, amount int) *Budget {
	return &Budget{
		HouseholdID: householdID,
		YearMonth:   yearMonth,
		Amount:      amount,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// IsNotified checks if the notification has already been sent.
func (b *Budget) IsNotified() bool {
	return b.NotifiedAt != nil
}

// SetNotified sets the notification time to the current time.
func (b *Budget) SetNotified() {
	now := time.Now()
	b.NotifiedAt = &now
	b.UpdatedAt = now
}
