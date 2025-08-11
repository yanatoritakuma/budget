package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yanatoritakuma/budget/back/model"
	"github.com/yanatoritakuma/budget/back/usecase"
)

type IExpenseController interface {
	CreateExpense(c *gin.Context)
}

type expenseController struct {
	eu usecase.IExpenseUsecase
}

func NewExpenseController(eu usecase.IExpenseUsecase) IExpenseController {
	return &expenseController{eu}
}

func (ec *expenseController) CreateExpense(c *gin.Context) {
	// ユーザーIDを取得（認証済みユーザーのコンテキストから）
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ユーザーが認証されていません"})
		return
	}

	// リクエストボディからexpenseデータをバインド
	expense := model.Expense{}
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不正なリクエストデータです: " + err.Error()})
		return
	}

	// ユーザーIDをセット
	expense.UserID = userID

	// 支出を作成
	expenseRes, err := ec.eu.CreateExpense(expense)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "支出の作成に失敗しました: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, expenseRes)
}
