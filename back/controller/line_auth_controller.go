package controller

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yanatoritakuma/budget/back/internal/api"
	"github.com/yanatoritakuma/budget/back/model"
	"github.com/yanatoritakuma/budget/back/usecase"
)

// LineLoginController はLINEログインコントローラのインターフェースです。
type LineLoginController interface {
	Login(c *gin.Context)
	Callback(c *gin.Context)
	LinkAccount(c *gin.Context)
	CreateAccount(c *gin.Context)
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

	// stateの検証
	sessionID := fmt.Sprintf("line-login-%s", state)
	if !ctrl.tokenStore.ValidateToken(sessionID, state) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid state"})
		return
	}
	ctrl.tokenStore.DeleteToken(sessionID)

	// LINEログインの処理
	token, claims, err := ctrl.lineLoginUsecase.LineLoginCallback(c.Request.Context(), code, state)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("LINE login failed: %v", err)})
		return
	}

	domain := os.Getenv("API_DOMAIN")
	isSecure := os.Getenv("GO_ENV") != "dev"

	// ユーザーが存在する場合（ログイン成功）
	if token != "" {
		tokenCookie := &http.Cookie{
			Name:     "token",
			Value:    token,
			MaxAge:   60 * 60 * 12, // 12時間
			Path:     "/",
			Domain:   domain,
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

		if !isSecure {
			tokenCookie.SameSite = http.SameSiteLaxMode
			loggedInCookie.SameSite = http.SameSiteLaxMode
		}

		http.SetCookie(c.Writer, tokenCookie)
		http.SetCookie(c.Writer, loggedInCookie)
		c.JSON(http.StatusOK, gin.H{"status": "logged_in", "message": "LINEログインに成功しました"})
		return
	}

	// ユーザーが存在しない場合（未登録）
	if claims != nil {
		// プレ認証トークンの生成
		preAuthToken, err := ctrl.lineLoginUsecase.GeneratePreAuthToken(claims.Subject, claims.Name, claims.Picture)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate pre-auth token"})
			return
		}

		preAuthCookie := &http.Cookie{
			Name:     "line_pre_auth",
			Value:    preAuthToken,
			MaxAge:   60 * 30, // 30分
			Path:     "/",
			Domain:   domain,
			Secure:   isSecure,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		}
		if !isSecure {
			preAuthCookie.SameSite = http.SameSiteLaxMode
		}
		http.SetCookie(c.Writer, preAuthCookie)

		c.JSON(http.StatusOK, gin.H{
			"status":       "unregistered",
			"line_name":    claims.Name,
			"line_picture": claims.Picture,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected login state"})
}

// LinkAccount は既存アカウントとLINEアカウントを紐付けます。
func (ctrl *LineLoginControllerImpl) LinkAccount(c *gin.Context) {
	preAuthToken, err := c.Cookie("line_pre_auth")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No pending LINE login found"})
		return
	}

	var req api.LinkAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	token, err := ctrl.lineLoginUsecase.LinkLineAccount(c.Request.Context(), preAuthToken, string(req.Email), req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 成功時のCookie設定
	ctrl.setLoginCookies(c, token)
	c.JSON(http.StatusOK, gin.H{"message": "Account linked successfully"})
}

// CreateAccount はLINEアカウントから新規ユーザーを作成します。
func (ctrl *LineLoginControllerImpl) CreateAccount(c *gin.Context) {
	preAuthToken, err := c.Cookie("line_pre_auth")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No pending LINE login found"})
		return
	}

	token, err := ctrl.lineLoginUsecase.CreateUserFromLine(c.Request.Context(), preAuthToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 成功時のCookie設定
	ctrl.setLoginCookies(c, token)
	c.JSON(http.StatusCreated, gin.H{"message": "Account created successfully"})
}

// setLoginCookies はログイン成功時の共通Cookie設定を行います。
func (ctrl *LineLoginControllerImpl) setLoginCookies(c *gin.Context, token string) {
	domain := os.Getenv("API_DOMAIN")
	isSecure := os.Getenv("GO_ENV") != "dev"

	// Pre-Auth Cookieの削除
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "line_pre_auth",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Domain:   domain,
		Secure:   isSecure,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	tokenCookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		MaxAge:   60 * 60 * 12,
		Path:     "/",
		Domain:   domain,
		Secure:   isSecure,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	loggedInCookie := &http.Cookie{
		Name:     "logged_in",
		Value:    "true",
		MaxAge:   60 * 60 * 12,
		Path:     "/",
		Domain:   domain,
		Secure:   isSecure,
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
	}

	if !isSecure {
		tokenCookie.SameSite = http.SameSiteLaxMode
		loggedInCookie.SameSite = http.SameSiteLaxMode
	}

	http.SetCookie(c.Writer, tokenCookie)
	http.SetCookie(c.Writer, loggedInCookie)
}
