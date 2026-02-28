package repository

import (
	"context"

	"github.com/yanatoritakuma/budget/back/domain/budget"
	"github.com/yanatoritakuma/budget/back/model"
	"gorm.io/gorm"
)

type budgetRepository struct {
	db *gorm.DB
}

func NewBudgetRepository(db *gorm.DB) budget.IBudgetRepository {
	return &budgetRepository{db}
}

func (r *budgetRepository) FindByHouseholdIDAndYearMonth(ctx context.Context, householdID uint, yearMonth string) (*budget.Budget, error) {
	var m budget.Budget
	if err := r.db.WithContext(ctx).Table("budgets").Where("household_id = ? AND year_month = ?", householdID, yearMonth).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *budgetRepository) Update(ctx context.Context, b *budget.Budget) error {
	m := model.Budget{
		ID:          b.ID,
		HouseholdID: b.HouseholdID,
		YearMonth:   b.YearMonth,
		Amount:      b.Amount,
		NotifiedAt:  b.NotifiedAt,
		UpdatedAt:   b.UpdatedAt,
	}
	return r.db.WithContext(ctx).Save(&m).Error
}
