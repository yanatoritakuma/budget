package usecase

import (
	"github.com/yanatoritakuma/budget/back/model"
	"github.com/yanatoritakuma/budget/back/repository"
)

type IExpenseUsecase interface {
	CreateExpense(expense model.Expense) (model.ExpenseResponse, error)
	GetExpense(year int, month int, category *string) ([]model.ExpenseResponse, error)
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

func (eu *expenseUsecase) GetExpense(year int, month int, category *string) ([]model.ExpenseResponse, error) {
	expenses, err := eu.er.GetExpense(year, month, category)
	if err != nil {
		return nil, err
	}

	var expenseResponses []model.ExpenseResponse
	for _, expense := range expenses {
		var payerName *string
		if expense.PayerID != nil {
			payerName = &expense.Payer.Name
		}

		expenseResponse := model.ExpenseResponse{
			ID:        expense.ID,
			UserID:    expense.UserID,
			Amount:    expense.Amount,
			StoreName: expense.StoreName,
			Date:      expense.Date,
			Category:  expense.Category,
			Memo:      expense.Memo,
			CreatedAt: expense.CreatedAt,
			PayerName: payerName,
		}
		expenseResponses = append(expenseResponses, expenseResponse)
	}

	return expenseResponses, nil
}
