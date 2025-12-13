package repository

import (
	"github.com/yanatoritakuma/budget/back/model"
	"gorm.io/gorm"
)

type IExpenseRepository interface {
	CreateExpense(expense *model.Expense) error
	GetExpense(householdID uint, year int, month int, category *string) ([]model.Expense, error)
	UpdateExpense(expense *model.Expense, expenseId uint) error
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
	// Join with user table to filter by household_id
	query := er.db.Joins(`JOIN "user" ON "user".id = expenses.user_id`).
		Where(`"user".household_id = ?`, householdID).
		Where("EXTRACT(YEAR FROM date) = ? AND EXTRACT(MONTH FROM date) = ?", year, month)

	if category != nil && *category != "" {
		query = query.Where("category = ?", *category)
	}

	if err := query.Find(&expenses).Error; err != nil {
		return nil, err
	}

	return expenses, nil
}

func (er *expenseRepository) UpdateExpense(expense *model.Expense, expenseId uint) error {
	if err := er.db.Model(&model.Expense{}).Where("id = ?", expenseId).Updates(expense).Error; err != nil {
		return err
	}
	return nil
}
