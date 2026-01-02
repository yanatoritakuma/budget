package usecase

import (
	"context"
	"fmt"

	"github.com/yanatoritakuma/budget/back/domain/user" // Added
	"github.com/yanatoritakuma/budget/back/internal/api"
	"github.com/yanatoritakuma/budget/back/model"
	"github.com/yanatoritakuma/budget/back/repository"
)

type IExpenseUsecase interface {
	CreateExpense(expense model.Expense) (api.ExpenseResponse, error)
	GetExpense(userID uint, year int, month int, category *string) ([]api.ExpenseResponse, error)
	UpdateExpense(expense model.Expense, expenseId uint) (api.ExpenseResponse, error)
	DeleteExpense(expenseId uint) error
}

type expenseUsecase struct {
	er repository.IExpenseRepository
	ur user.UserRepository
}

func NewExpenseUsecase(er repository.IExpenseRepository, ur user.UserRepository) IExpenseUsecase {
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
	ctx := context.Background()

	// Get user to find their household ID
	currentUser, err := eu.ur.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}
	if currentUser == nil {
		return nil, fmt.Errorf("current user not found")
	}

	expenses, err := eu.er.GetExpense(currentUser.HouseholdID, year, month, category)
	if err != nil {
		return nil, err
	}

	var expenseResponses []api.ExpenseResponse
	for _, expense := range expenses {

		var payerName string
		payer, err := eu.ur.FindByID(ctx, expense.PayerID)
		if err != nil {
			payerName = "不明"
		} else if payer == nil {
			payerName = "不明"
		} else {
			payerName = payer.Name
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
			PayerName: &payerName,
		}
		expenseResponses = append(expenseResponses, expenseResponse)
	}

	return expenseResponses, nil
}

func (eu *expenseUsecase) UpdateExpense(expense model.Expense, expenseId uint) (api.ExpenseResponse, error) {
	ctx := context.Background()

	if err := eu.er.UpdateExpense(&expense, expenseId); err != nil {
		return api.ExpenseResponse{}, err
	}

	payer, err := eu.ur.FindByID(ctx, expense.UserID)
	if err != nil {
		return api.ExpenseResponse{}, err
	}
	if payer == nil {
		return api.ExpenseResponse{}, fmt.Errorf("payer not found")
	}

	payerName := payer.Name

	resExpense := api.ExpenseResponse{
		Id:        int(expenseId),
		UserId:    int(expense.UserID),
		Amount:    expense.Amount,
		StoreName: expense.StoreName,
		Date:      expense.Date,
		Category:  expense.Category,
		Memo:      &expense.Memo,
		CreatedAt: expense.CreatedAt,
		PayerName: &payerName,
	}
	return resExpense, nil
}

func (eu *expenseUsecase) DeleteExpense(expenseId uint) error {
	if err := eu.er.DeleteExpense(expenseId); err != nil {
		return err
	}
	return nil
}
