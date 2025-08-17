package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yanatoritakuma/budget/back/model"
	"github.com/yanatoritakuma/budget/back/usecase"
)

type IExpenseController interface {
	CreateExpense(c *gin.Context)
	GetExpense(c *gin.Context)
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

func (ec *expenseController) GetExpense(c *gin.Context) {
	// ユーザーIDを取得（認証済みユーザーのコンテキストから）
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ユーザーが認証されていません"})
		return
	}

	// クエリパラメータから年と月を取得
	year := c.Query("year")
	month := c.Query("month")
	category := c.Query("category")

	// 年月の必須チェック
	if year == "" || month == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "年と月は必須パラメータです"})
		return
	}

	// 文字列を数値に変換
	yearInt := 0
	monthInt := 0
	var err error

	yearInt, err = strconv.Atoi(year)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不正な年のフォーマットです"})
		return
	}

	monthInt, err = strconv.Atoi(month)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不正な月のフォーマットです"})
		return
	}

	// 月の範囲チェック
	if monthInt < 1 || monthInt > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "月は1から12の間で指定してください"})
		return
	}

	var categoryPtr *string
	if category != "" {
		categoryPtr = &category
	}

	// 支出データを取得
	expenses, err := ec.eu.GetExpense(userID, yearInt, monthInt, categoryPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "支出データの取得に失敗しました: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, expenses)
}
