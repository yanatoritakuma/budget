package budget

import "context"

type BudgetExceededEvent struct {
	Type          string `json:"type"`
	HouseholdID   uint   `json:"household_id"`
	YearMonth     string `json:"year_month"`
	BudgetAmount  int    `json:"budget_amount"`
	CurrentAmount int    `json:"current_amount"`
}

type INotificationService interface {
	SendBudgetExceededNotification(ctx context.Context, event BudgetExceededEvent) error
}
