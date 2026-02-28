package budget

import "context"

type IBudgetRepository interface {
	FindByHouseholdIDAndYearMonth(ctx context.Context, householdID uint, yearMonth string) (*Budget, error)
	Update(ctx context.Context, budget *Budget) error
}
