package usecase

import (
	"github.com/yanatoritakuma/budget/back/internal/api"
	"github.com/yanatoritakuma/budget/back/model"
	"github.com/yanatoritakuma/budget/back/repository"
)

type IExpenseUsecase interface {
	CreateExpense(expense model.Expense) (api.ExpenseResponse, error)
	GetExpense(userID uint, year int, month int, category *string) ([]api.ExpenseResponse, error)
}

type expenseUsecase struct {
	er repository.IExpenseRepository
	ur repository.IUserRepository
}

func NewExpenseUsecase(er repository.IExpenseRepository, ur repository.IUserRepository) IExpenseUsecase {
	return &expenseUsecase{er: er, ur: ur}
}

func (eu *expenseUsecase) CreateExpense(expense model.Expense) (api.ExpenseResponse, error) {
	if err := eu.er.CreateExpense(&expense); err != nil {
		return api.ExpenseResponse{}, err
	}

	resExpense := api.ExpenseResponse{
		Id:        int(expense.ID),
		UserId:    int(expense.UserID),
		Amount:    expense.Amount,
		StoreName: expense.StoreName,
		Date:      expense.Date,
		Category:  expense.Category,
		Memo:      &expense.Memo,
		CreatedAt: expense.CreatedAt,
	}

	return resExpense, nil
}

func (eu *expenseUsecase) GetExpense(userID uint, year int, month int, category *string) ([]api.ExpenseResponse, error) {
	// Get user to find their household ID
	var user model.User
	if err := eu.ur.GetUserByID(&user, userID); err != nil {
		return nil, err
	}

	expenses, err := eu.er.GetExpense(user.HouseholdID, year, month, category)
	if err != nil {
		return nil, err
	}

	var expenseResponses []api.ExpenseResponse
	for _, expense := range expenses {
		var payerName *string
		if expense.PayerID != nil {
			// Ensure Payer is not nil before accessing Name
			if expense.Payer.ID != 0 {
				payerName = &expense.Payer.Name
			}
		}

		expenseResponse := api.ExpenseResponse{
			Id:        int(expense.ID),
			UserId:    int(expense.UserID),
			Amount:    expense.Amount,
			StoreName: expense.StoreName,
			Date:      expense.Date,
			Category:  expense.Category,
			Memo:      &expense.Memo,
			CreatedAt: expense.CreatedAt,
			PayerName: payerName,
		}
		expenseResponses = append(expenseResponses, expenseResponse)
	}

	return expenseResponses, nil
}
