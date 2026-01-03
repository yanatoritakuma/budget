package usecase

import (
	"context"
	"fmt"

	"github.com/yanatoritakuma/budget/back/domain/expense"
	"github.com/yanatoritakuma/budget/back/domain/user"
	"github.com/yanatoritakuma/budget/back/internal/api"
)

type IExpenseUsecase interface {
	CreateExpense(req api.ExpenseRequest) (api.ExpenseResponse, error)
	GetExpense(userID uint, year int, month int, category *string) ([]api.ExpenseResponse, error)
	UpdateExpense(req api.ExpenseRequest, expenseId uint) (api.ExpenseResponse, error)
	DeleteExpense(expenseId uint) error
}

type expenseUsecase struct {
	er expense.IExpenseRepository
	ur user.IUserRepository
}

func NewExpenseUsecase(er expense.IExpenseRepository, ur user.IUserRepository) IExpenseUsecase {
	return &expenseUsecase{er: er, ur: ur}
}

func (eu *expenseUsecase) CreateExpense(req api.ExpenseRequest) (api.ExpenseResponse, error) {
	memo := ""
	if req.Memo != nil {
		memo = *req.Memo
	}
	domainExpense, err := expense.NewExpense(
		req.Amount,
		req.StoreName,
		req.Date,
		req.Category,
		memo,
		uint(req.UserId),
		uint(req.UserId), // PayerID is the same as UserID for now
	)
	if err != nil {
		return api.ExpenseResponse{}, err
	}

	if err := eu.er.CreateExpense(domainExpense); err != nil {
		return api.ExpenseResponse{}, err
	}

	resExpense := api.ExpenseResponse{
		Id:        int(domainExpense.ID),
		UserId:    int(domainExpense.UserID),
		Amount:    domainExpense.Amount,
		StoreName: domainExpense.StoreName,
		Date:      domainExpense.Date,
		Category:  domainExpense.Category,
		Memo:      &domainExpense.Memo,
		CreatedAt: domainExpense.CreatedAt,
	}

	return resExpense, nil
}

func (eu *expenseUsecase) GetExpense(userID uint, year int, month int, category *string) ([]api.ExpenseResponse, error) {
	ctx := context.Background()

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
	for _, domainExpense := range expenses {
		var payerName string
		payer, err := eu.ur.FindByID(ctx, domainExpense.PayerID)
		if err != nil {
			payerName = "不明"
		} else if payer == nil {
			payerName = "不明"
		} else {
			payerName = payer.Name.Value()
		}

		expenseResponse := api.ExpenseResponse{
			Id:        int(domainExpense.ID),
			UserId:    int(domainExpense.UserID),
			Amount:    domainExpense.Amount,
			StoreName: domainExpense.StoreName,
			Date:      domainExpense.Date,
			Category:  domainExpense.Category,
			Memo:      &domainExpense.Memo,
			CreatedAt: domainExpense.CreatedAt,
			PayerName: &payerName,
		}
		expenseResponses = append(expenseResponses, expenseResponse)
	}

	return expenseResponses, nil
}

func (eu *expenseUsecase) UpdateExpense(req api.ExpenseRequest, expenseId uint) (api.ExpenseResponse, error) {
	ctx := context.Background()

	memo := ""
	if req.Memo != nil {
		memo = *req.Memo
	}

	domainExpense, err := expense.NewExpense(
		req.Amount,
		req.StoreName,
		req.Date,
		req.Category,
		memo,
		uint(req.UserId),
		uint(req.UserId), // PayerID is the same as UserID for now
	)
	if err != nil {
		return api.ExpenseResponse{}, err
	}
	domainExpense.ID = expenseId

	if err := eu.er.UpdateExpense(domainExpense, expenseId); err != nil {
		return api.ExpenseResponse{}, err
	}

	payer, err := eu.ur.FindByID(ctx, domainExpense.UserID)
	if err != nil {
		return api.ExpenseResponse{}, err
	}
	if payer == nil {
		return api.ExpenseResponse{}, fmt.Errorf("payer not found")
	}

	payerName := payer.Name.Value()

	resExpense := api.ExpenseResponse{
		Id:        int(expenseId),
		UserId:    int(domainExpense.UserID),
		Amount:    domainExpense.Amount,
		StoreName: domainExpense.StoreName,
		Date:      domainExpense.Date,
		Category:  domainExpense.Category,
		Memo:      &domainExpense.Memo,
		CreatedAt: domainExpense.CreatedAt,
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
