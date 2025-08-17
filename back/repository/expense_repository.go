package repository

import (
	"github.com/yanatoritakuma/budget/back/model"
	"gorm.io/gorm"
)

type IExpenseRepository interface {
	CreateExpense(expense *model.Expense) error
	GetExpense(householdID uint, year int, month int, category *string) ([]model.Expense, error)
}

type expenseRepository struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) IExpenseRepository {
	return &expenseRepository{db}
}

func (er *expenseRepository) CreateExpense(expense *model.Expense) error {
	if err := er.db.Create(expense).Error; err != nil {
		return err
	}
	return nil
}

func (er *expenseRepository) GetExpense(householdID uint, year int, month int, category *string) ([]model.Expense, error) {
	var expenses []model.Expense
	// Join with users table to filter by household_id
	query := er.db.Joins("JOIN users ON users.id = expenses.user_id").
		Where("users.household_id = ?", householdID).
		Where("EXTRACT(YEAR FROM date) = ? AND EXTRACT(MONTH FROM date) = ?", year, month)

	if category != nil && *category != "" {
		query = query.Where("category = ?", *category)
	}

	// Preload the Payer information (which is a User)
	if err := query.Preload("Payer").Find(&expenses).Error; err != nil {
		return nil, err
	}

	return expenses, nil
}
