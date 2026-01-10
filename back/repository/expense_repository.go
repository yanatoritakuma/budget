package repository

import (
	"context"

	"github.com/yanatoritakuma/budget/back/domain/expense"
	"github.com/yanatoritakuma/budget/back/model"
	"gorm.io/gorm"
)

var _ expense.ExpenseRepository = (*ExpenseRepositoryImpl)(nil)

type ExpenseRepositoryImpl struct {
	db *gorm.DB
}

func NewExpenseRepositoryImpl(db *gorm.DB) expense.ExpenseRepository {
	return &ExpenseRepositoryImpl{db}
}

func (er *ExpenseRepositoryImpl) CreateExpense(ctx context.Context, e *expense.Expense) error {
	expenseModel := toModelExpense(e)
	if err := er.db.WithContext(ctx).Create(expenseModel).Error; err != nil {
		return err
	}
	e.ID = expense.ExpenseID(expenseModel.ID)
	e.CreatedAt = expenseModel.CreatedAt
	e.UpdatedAt = expenseModel.UpdatedAt
	return nil
}

func (er *ExpenseRepositoryImpl) GetExpense(ctx context.Context, householdID uint, year int, month int, category *string) ([]*expense.Expense, error) {
	var expenseModels []model.Expense
	query := er.db.WithContext(ctx).Table("expenses").
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
	for i := range expenseModels {
		domainExpense, err := toDomainExpense(&expenseModels[i])
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, domainExpense)
	}

	return expenses, nil
}

func (er *ExpenseRepositoryImpl) UpdateExpense(ctx context.Context, e *expense.Expense) error {
	expenseModel := toModelExpense(e)
	return er.db.WithContext(ctx).Model(&model.Expense{}).Where("id = ?", e.ID.Value()).Updates(expenseModel).Error
}

func (er *ExpenseRepositoryImpl) DeleteExpense(ctx context.Context, expenseId expense.ExpenseID) error {
	return er.db.WithContext(ctx).Where("id = ?", expenseId.Value()).Delete(&model.Expense{}).Error
}

func toDomainExpense(em *model.Expense) (*expense.Expense, error) {
	if em == nil {
		return nil, nil
	}

	amount, err := expense.NewAmount(em.Amount)
	if err != nil {
		return nil, err
	}
	storeName, err := expense.NewStoreName(em.StoreName)
	if err != nil {
		return nil, err
	}
	category, err := expense.NewCategory(em.Category)
	if err != nil {
		return nil, err
	}
	memo, err := expense.NewMemo(em.Memo)
	if err != nil {
		return nil, err
	}

	return &expense.Expense{
		ID:        expense.ExpenseID(em.ID),
		Amount:    amount,
		StoreName: storeName,
		Date:      em.Date,
		Category:  category,
		Memo:      memo,
		CreatedAt: em.CreatedAt,
		UpdatedAt: em.UpdatedAt,
		UserID:    expense.UserID(em.UserID),
		PayerID:   expense.PayerID(em.PayerID),
	}, nil
}

func toModelExpense(e *expense.Expense) *model.Expense {
	if e == nil {
		return nil
	}
	return &model.Expense{
		ID:        e.ID.Value(),
		Amount:    e.Amount.Value(),
		StoreName: e.StoreName.Value(),
		Date:      e.Date,
		Category:  e.Category.Value(),
		Memo:      e.Memo.Value(),
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		UserID:    uint(e.UserID),
		PayerID:   uint(e.PayerID),
	}
}
