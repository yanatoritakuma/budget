package usecase

import (
	"github.com/yanatoritakuma/budget/back/model"
	"github.com/yanatoritakuma/budget/back/repository"
)

type IExpenseUsecase interface {
	CreateExpense(expense model.Expense) (model.ExpenseResponse, error)
}

type expenseUsecase struct {
	er repository.IExpenseRepository
}

func NewExpenseUsecase(er repository.IExpenseRepository) IExpenseUsecase {
	return &expenseUsecase{er: er}
}

func (eu *expenseUsecase) CreateExpense(expense model.Expense) (model.ExpenseResponse, error) {
	if err := eu.er.CreateExpense(&expense); err != nil {
		return model.ExpenseResponse{}, err
	}

	resExpense := model.ExpenseResponse{
		ID:        expense.ID,
		UserID:    expense.UserID,
		Amount:    expense.Amount,
		StoreName: expense.StoreName,
		Date:      expense.Date,
		Category:  expense.Category,
		Memo:      expense.Memo,
		CreatedAt: expense.CreatedAt,
	}

	return resExpense, nil
}
