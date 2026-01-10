package repository

import (
	"github.com/yanatoritakuma/budget/back/domain/expense"
	"github.com/yanatoritakuma/budget/back/model"
	"gorm.io/gorm"
)

var _ expense.IExpenseRepository = (*ExpenseRepositoryImpl)(nil)

type ExpenseRepositoryImpl struct {
	db *gorm.DB
}

func NewExpenseRepositoryImpl(db *gorm.DB) expense.IExpenseRepository {
	return &ExpenseRepositoryImpl{db}
}

func (er *ExpenseRepositoryImpl) CreateExpense(e *expense.Expense) error {
	expenseModel := toModelExpense(e)
	if err := er.db.Create(expenseModel).Error; err != nil {
		return err
	}
	e.ID = expenseModel.ID
	e.CreatedAt = expenseModel.CreatedAt
	e.UpdatedAt = expenseModel.UpdatedAt
	return nil
}

func (er *ExpenseRepositoryImpl) GetExpense(householdID uint, year int, month int, category *string) ([]*expense.Expense, error) {
	var expenseModels []model.Expense
	query := er.db.Table("expenses").
		Joins(`JOIN "user" ON "user".id = expenses.user_id`).
		Where(`"user".household_id = ?`, householdID).
		Where("EXTRACT(YEAR FROM date) = ? AND EXTRACT(MONTH FROM date) = ?", year, month)

	if category != nil && *category != "" {
		query = query.Where("category = ?", *category)
	}

	if err := query.Find(&expenseModels).Error; err != nil {
		return nil, err
	}

	var expenses []*expense.Expense
	for _, em := range expenseModels {
		expenses = append(expenses, toDomainExpense(&em))
	}

	return expenses, nil
}

func (er *ExpenseRepositoryImpl) UpdateExpense(e *expense.Expense, expenseId uint) error {
	expenseModel := toModelExpense(e)
	if err := er.db.Model(&model.Expense{}).Where("id = ?", expenseId).Updates(expenseModel).Error; err != nil {
		return err
	}
	return nil
}

func (er *ExpenseRepositoryImpl) DeleteExpense(expenseId uint) error {
	if err := er.db.Where("id = ?", expenseId).Delete(&model.Expense{}).Error; err != nil {
		return err
	}
	return nil
}

func toDomainExpense(em *model.Expense) *expense.Expense {
	if em == nil {
		return nil
	}
	return &expense.Expense{
		ID:        em.ID,
		Amount:    em.Amount,
		StoreName: em.StoreName,
		Date:      em.Date,
		Category:  em.Category,
		Memo:      em.Memo,
		CreatedAt: em.CreatedAt,
		UpdatedAt: em.UpdatedAt,
		UserID:    em.UserID,
		PayerID:   em.PayerID,
	}
}

func toModelExpense(e *expense.Expense) *model.Expense {
	if e == nil {
		return nil
	}
	return &model.Expense{
		ID:        e.ID,
		Amount:    e.Amount,
		StoreName: e.StoreName,
		Date:      e.Date,
		Category:  e.Category,
		Memo:      e.Memo,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		UserID:    e.UserID,
		PayerID:   e.PayerID,
	}
}
