package usecase

import (
	"context"
	"fmt"

	"github.com/yanatoritakuma/budget/back/domain/budget"
	"github.com/yanatoritakuma/budget/back/domain/expense"
	"github.com/yanatoritakuma/budget/back/domain/user"
	"github.com/yanatoritakuma/budget/back/internal/api"
)

type ExpenseUsecase interface {
	CreateExpense(ctx context.Context, req api.ExpenseRequest) (api.ExpenseResponse, error)
	GetExpense(ctx context.Context, userID uint, year int, month int, category *string) ([]api.ExpenseResponse, error)
	UpdateExpense(ctx context.Context, req api.ExpenseRequest, expenseId uint) (api.ExpenseResponse, error)
	DeleteExpense(ctx context.Context, expenseId uint) error
}

type expenseUsecase struct {
	er  expense.ExpenseRepository
	ur  user.UserRepository
	br  budget.IBudgetRepository
	uow UnitOfWork
	ns  budget.INotificationService
}

func NewExpenseUsecase(
	er expense.ExpenseRepository,
	ur user.UserRepository,
	br budget.IBudgetRepository,
	uow UnitOfWork,
	ns budget.INotificationService,
) ExpenseUsecase {
	return &expenseUsecase{
		er:  er,
		ur:  ur,
		br:  br,
		uow: uow,
		ns:  ns,
	}
}

func (eu *expenseUsecase) CreateExpense(ctx context.Context, req api.ExpenseRequest) (api.ExpenseResponse, error) {
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

	var resExpense api.ExpenseResponse

	err = eu.uow.Transaction(func(repos Repositories) error {
		// 支出の登録
		if err := repos.Expense.CreateExpense(ctx, domainExpense); err != nil {
			return err
		}

		// ユーザー情報の取得（世帯ID取得のため）
		currentUser, err := repos.User.FindByID(ctx, uint(req.UserId))
		if err != nil {
			return err
		}
		if currentUser == nil {
			return fmt.Errorf("user not found")
		}

		// 予算超過チェック
		year := domainExpense.Date.Year()
		month := int(domainExpense.Date.Month())
		yearMonth := fmt.Sprintf("%d/%02d", year, month)

		b, err := repos.Budget.FindByHouseholdIDAndYearMonth(ctx, currentUser.HouseholdID, yearMonth)
		if err != nil {
			return err
		}

		// 予算設定があり、かつ未通知の場合のみチェック
		if b != nil && !b.IsNotified() {
			total, err := repos.Expense.GetTotalAmountByMonth(ctx, currentUser.HouseholdID, year, month)
			if err != nil {
				return err
			}

			if total > b.Amount {
				// SQS通知送信
				event := budget.BudgetExceededEvent{
					Type:          "BUDGET_EXCEEDED",
					HouseholdID:   currentUser.HouseholdID,
					YearMonth:     yearMonth,
					BudgetAmount:  b.Amount,
					CurrentAmount: total,
				}

				// 通知送信（NotificationServiceが初期化されている場合のみ）
				if eu.ns != nil {
					if err := eu.ns.SendBudgetExceededNotification(ctx, event); err != nil {
						// 通知失敗をログに記録し、支出登録は継続する（通知はベストエフォート）
						fmt.Printf("Warning: failed to send notification: %v\n", err)
					} else {
						// 通知成功時のみ通知済みフラグを更新
						b.SetNotified()
						if err := repos.Budget.Update(ctx, b); err != nil {
							return err
						}
					}
				}
			}
		}

		resExpense = api.ExpenseResponse{
			Id:        int(domainExpense.ID.Value()),
			UserId:    int(domainExpense.UserID),
			Amount:    domainExpense.Amount.Value(),
			StoreName: domainExpense.StoreName.Value(),
			Date:      domainExpense.Date,
			Category:  domainExpense.Category.Value(),
			Memo:      &memo,
			CreatedAt: domainExpense.CreatedAt,
		}

		return nil
	})

	if err != nil {
		return api.ExpenseResponse{}, err
	}

	return resExpense, nil
}

func (eu *expenseUsecase) GetExpense(ctx context.Context, userID uint, year int, month int, category *string) ([]api.ExpenseResponse, error) {
	currentUser, err := eu.ur.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}
	if currentUser == nil {
		return nil, fmt.Errorf("current user not found")
	}

	expenses, err := eu.er.GetExpense(ctx, currentUser.HouseholdID, year, month, category)
	if err != nil {
		return nil, err
	}

	var expenseResponses []api.ExpenseResponse
	for _, domainExpense := range expenses {
		var payerName string
		payer, err := eu.ur.FindByID(ctx, uint(domainExpense.PayerID))
		if err != nil {
			payerName = "不明"
		} else if payer == nil {
			payerName = "不明"
		} else {
			payerName = payer.Name.Value()
		}

		memo := domainExpense.Memo.Value()
		expenseResponse := api.ExpenseResponse{
			Id:        int(domainExpense.ID.Value()),
			UserId:    int(domainExpense.UserID),
			Amount:    domainExpense.Amount.Value(),
			StoreName: domainExpense.StoreName.Value(),
			Date:      domainExpense.Date,
			Category:  domainExpense.Category.Value(),
			Memo:      &memo,
			CreatedAt: domainExpense.CreatedAt,
			PayerName: &payerName,
		}
		expenseResponses = append(expenseResponses, expenseResponse)
	}

	return expenseResponses, nil
}

func (eu *expenseUsecase) UpdateExpense(ctx context.Context, req api.ExpenseRequest, expenseId uint) (api.ExpenseResponse, error) {
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
	domainExpense.ID = expense.ExpenseID(expenseId)

	if err := eu.er.UpdateExpense(ctx, domainExpense); err != nil {
		return api.ExpenseResponse{}, err
	}

	payer, err := eu.ur.FindByID(ctx, uint(domainExpense.UserID))
	if err != nil {
		return api.ExpenseResponse{}, err
	}
	if payer == nil {
		return api.ExpenseResponse{}, fmt.Errorf("payer not found")
	}

	payerName := payer.Name.Value()
	resMemo := domainExpense.Memo.Value()

	resExpense := api.ExpenseResponse{
		Id:        int(domainExpense.ID.Value()),
		UserId:    int(domainExpense.UserID),
		Amount:    domainExpense.Amount.Value(),
		StoreName: domainExpense.StoreName.Value(),
		Date:      domainExpense.Date,
		Category:  domainExpense.Category.Value(),
		Memo:      &resMemo,
		CreatedAt: domainExpense.CreatedAt,
		PayerName: &payerName,
	}
	return resExpense, nil
}

func (eu *expenseUsecase) DeleteExpense(ctx context.Context, expenseId uint) error {
	if err := eu.er.DeleteExpense(ctx, expense.ExpenseID(expenseId)); err != nil {
		return err
	}
	return nil
}
