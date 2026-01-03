package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yanatoritakuma/budget/back/internal/api"
	"github.com/yanatoritakuma/budget/back/usecase"
)

type IExpenseController interface {
	CreateExpense(c *gin.Context)
	GetExpense(c *gin.Context)
	UpdateExpense(c *gin.Context)
	DeleteExpense(c *gin.Context)
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
	var req api.ExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不正なリクエストデータです: " + err.Error()})
		return
	}

	// ユーザーIDをセット
	req.UserId = int(userID)

	// 支出を作成
	expenseRes, err := ec.eu.CreateExpense(req)
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

func (ec *expenseController) UpdateExpense(c *gin.Context) {
	// ユーザーIDを取得（認証済みユーザーのコンテキストから）
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ユーザーが認証されていません"})
		return
	}

	// パスパラメータからIDを取得
	id := c.Param("id")
	expenseId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不正なIDフォーマットです"})
		return
	}

	// リクエストボディからexpenseデータをバインド
	var req api.ExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不正なリクエストデータです: " + err.Error()})
		return
	}

	// 支出を更新
	expenseRes, err := ec.eu.UpdateExpense(req, uint(expenseId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "支出の更新に失敗しました: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, expenseRes)
}

func (ec *expenseController) DeleteExpense(c *gin.Context) {
	// パスパラメータからIDを取得
	id := c.Param("id")
	expenseId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不正なIDフォーマットです"})
		return
	}

	// 支出を削除
	if err := ec.eu.DeleteExpense(uint(expenseId)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "支出の削除に失敗しました: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
