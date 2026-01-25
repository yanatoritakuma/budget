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

	// セッションIDをユーザーごとに生成（state 自体をキーとして使用）
	// または HTTP-only Cookie にセットして、リクエスト時に取得
	sessionID := fmt.Sprintf("line-login-%s", state)
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

	// stateの検証（セッションIDとして state を使用）
	sessionID := fmt.Sprintf("line-login-%s", state)
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

	// Cookie domain を安全に設定（ホスト名のみ）
	domain := os.Getenv("API_DOMAIN")
	isSecure := os.Getenv("GO_ENV") != "dev" // 本番環境では secure=true

	// http.Cookie構造体を作成
	tokenCookie := &http.Cookie{
		Name:     "token",
		Value:    jwtToken,
		MaxAge:   60 * 60 * 12, // 12時間
		Path:     "/",
		Domain:   os.Getenv("API_DOMAIN"),
		Secure:   isSecure,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	loggedInCookie := &http.Cookie{
		Name:     "logged_in",
		Value:    "true",
		MaxAge:   60 * 60 * 12, // 12時間
		Path:     "/",
		Domain:   domain,
		Secure:   isSecure,
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
	}

	// 開発環境など非セキュアな場合はSameSiteをLaxに設定
	if !isSecure {
		tokenCookie.SameSite = http.SameSiteLaxMode
		loggedInCookie.SameSite = http.SameSiteLaxMode
	}

	// レスポンスヘッダーに直接Cookieを設定
	http.SetCookie(c.Writer, tokenCookie)
	http.SetCookie(c.Writer, loggedInCookie)

	// フロントエンドのログイン成功時のリダイレクトURL
	redirectURL := fmt.Sprintf("%s/budget", os.Getenv("FE_URL")) // FE_URLはフロントエンドのドメイン
	if os.Getenv("FE_URL") == "" {
		redirectURL = "http://localhost:3000/budget" // ローカル開発用フォールバック
	}
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
