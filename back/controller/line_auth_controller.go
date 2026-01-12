package controller

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yanatoritakuma/budget/back/model"
	"github.com/yanatoritakuma/budget/back/usecase"
)

// LineLoginController はLINEログインコントローラのインターフェースです。
type LineLoginController interface {
	Login(c *gin.Context)
	Callback(c *gin.Context)
}

// LineLoginControllerImpl はLineLoginControllerの実装です。
type LineLoginControllerImpl struct {
	lineLoginUsecase usecase.LineLoginUsecase
	tokenStore       *model.TokenStore // CSRF stateの保存用
}

// NewLineLoginController はLineLoginControllerの新しいインスタンスを生成します。
func NewLineLoginController(lineLoginUsecase usecase.LineLoginUsecase) LineLoginController {
	return &LineLoginControllerImpl{
		lineLoginUsecase: lineLoginUsecase,
		tokenStore:       model.NewTokenStore(), // 専用のトークンストア
	}
}

// Login はLINE認証開始のためのURLを返します。
func (ctrl *LineLoginControllerImpl) Login(c *gin.Context) {
	// CSRF対策のためのstateを生成し、セッションに保存
	state, err := usecase.GenerateState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate state"})
		return
	}
	// TODO: セッションIDはどこから取得するか検討（Cookieからなど）
	// 現状は固定のSessionIDを使用し、本番環境ではセキュアなものにする
	sessionID := "line-login-state" // 仮のセッションID
	ctrl.tokenStore.SaveToken(sessionID, model.CSRFToken{
		Token:     state,
		ExpiresAt: time.Now().Add(5 * time.Minute), // 5分間有効
	})

	authURL, err := ctrl.lineLoginUsecase.GetLineAuthURL(c.Request.Context(), state)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get LINE auth URL: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"auth_url": authURL})
}

// Callback はLINE認証後のコールバックを処理します。
func (ctrl *LineLoginControllerImpl) Callback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing code or state in callback"})
		return
	}

	// stateの検証
	sessionID := "line-login-state" // Loginメソッドで保存した仮のセッションID
	if !ctrl.tokenStore.ValidateToken(sessionID, state) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid state"})
		return
	}
	ctrl.tokenStore.DeleteToken(sessionID) // stateは一度使用したら削除

	// LINEログインの処理を行い、JWTを取得
	jwtToken, err := ctrl.lineLoginUsecase.LineLoginCallback(c.Request.Context(), code, state)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("LINE login failed: %v", err)})
		return
	}

	// JWTをCookieにセット
	c.SetCookie("token", jwtToken, 60*60*12, "/", os.Getenv("FE_URL"), true, true) // 12時間有効
	c.SetCookie("logged_in", "true", 60*60*12, "/", os.Getenv("FE_URL"), true, false)

	// フロントエンドのログイン成功時のリダイレクトURL
	redirectURL := fmt.Sprintf("%s/budget", os.Getenv("FE_URL")) // FE_URLはフロントエンドのドメイン
	if os.Getenv("FE_URL") == "" {
		redirectURL = "http://localhost:3000/budget" // ローカル開発用フォールバック
	}
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
