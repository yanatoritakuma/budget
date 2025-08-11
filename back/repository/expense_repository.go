package repository

import (
	"github.com/yanatoritakuma/budget/back/model"
	"gorm.io/gorm"
)

type IExpenseRepository interface {
	CreateExpense(expense *model.Expense) error
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
