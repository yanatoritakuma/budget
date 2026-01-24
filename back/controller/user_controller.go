package controller

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yanatoritakuma/budget/back/internal/api"
	"github.com/yanatoritakuma/budget/back/usecase"
	"github.com/yanatoritakuma/budget/back/utils"
)

type UserController interface {
	SignUp(c *gin.Context)
	LogIn(c *gin.Context)
	LogOut(c *gin.Context)
	CsrfToken(c *gin.Context)
	GetLoggedInUser(c *gin.Context)
	GetHouseholdUsers(c *gin.Context)
	JoinHousehold(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	ValidateCSRFToken(sessionID, token string) bool
}

type userController struct {
	uu usecase.UserUsecase
}

func NewUserController(uu usecase.UserUsecase) UserController {
	return &userController{uu}
}

func (uc *userController) SignUp(c *gin.Context) {
	user := api.SignUpRequest{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userRes, err := uc.uu.SignUp(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, userRes)
}

func (uc *userController) LogIn(c *gin.Context) {
	user := api.SignUpRequest{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokenString, err := uc.uu.Login(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		MaxAge:   int(time.Hour * 24 / time.Second),
		Path:     "/",
		Domain:   os.Getenv("API_DOMAIN"),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		HttpOnly: true,
	})
	c.Status(http.StatusOK)
}

func (uc *userController) LogOut(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Domain:   os.Getenv("API_DOMAIN"),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		HttpOnly: true,
	})
	c.Status(http.StatusOK)
}

func (uc *userController) GetLoggedInUser(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userRes, err := uc.uu.GetLoggedInUser(cookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, userRes)
}

func (uc *userController) CsrfToken(c *gin.Context) {
	// セッションIDの取得
	sessionID, err := c.Cookie("token")
	if err != nil {
		sessionID = "default" // フォールバック値
	}

	// 既存のトークンを取得または新しいトークンを生成
	token, err := uc.uu.GetOrGenerateCSRFToken(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to handle CSRF token"})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "csrf_token",
		Value:    token,
		MaxAge:   int(time.Hour / time.Second),
		Path:     "/",
		Domain:   utils.ExtractHostname(os.Getenv("FE_URL")),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		HttpOnly: false,
	})

	c.JSON(http.StatusOK, gin.H{
		"csrf_token": token,
	})
}

// ValidateCSRFToken はCSRFトークンを検証します
func (uc *userController) ValidateCSRFToken(sessionID, token string) bool {
	return uc.uu.ValidateCSRFToken(sessionID, token)
}

func (uc *userController) UpdateUser(c *gin.Context) {
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	claims := userClaims.(jwt.MapClaims)
	userId := uint(claims["user_id"].(float64))

	var req api.UserUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRes, err := uc.uu.UpdateUser(userId, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userRes)
}

func (uc *userController) DeleteUser(c *gin.Context) {
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	claims := userClaims.(jwt.MapClaims)
	userId := uint(claims["user_id"].(float64))

	err := uc.uu.DeleteUser(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (uc *userController) GetHouseholdUsers(c *gin.Context) {
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	claims := userClaims.(jwt.MapClaims)
	userId := uint(claims["user_id"].(float64))

	users, err := uc.uu.GetHouseholdUsers(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

type JoinHouseholdRequest struct {
	InviteCode string `json:"invite_code"`
}

func (uc *userController) JoinHousehold(c *gin.Context) {
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	claims := userClaims.(jwt.MapClaims)
	userId := uint(claims["user_id"].(float64))

	var req JoinHouseholdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := uc.uu.JoinHousehold(userId, req.InviteCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
