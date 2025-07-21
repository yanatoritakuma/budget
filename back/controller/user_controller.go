package controller

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/yanatoritakuma/budget/back/model"
	"github.com/yanatoritakuma/budget/back/usecase"
)

type IUserController interface {
	SignUp(c *gin.Context)
	LogIn(c *gin.Context)
	LogOut(c *gin.Context)
	CsrfToken(c *gin.Context)
	GetLoggedInUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type userController struct {
	uu usecase.IUserUsecase
}

func NewUserController(uu usecase.IUserUsecase) IUserController {
	return &userController{uu}
}

func (uc *userController) SignUp(c *gin.Context) {
	user := model.User{}
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
	user := model.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokenString, err := uc.uu.Login(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie(
		"token",
		tokenString,
		int(time.Hour*24/time.Second), // MaxAge in seconds
		"/",
		os.Getenv("API_DOMAIN"),
		true, // secure
		true, // httpOnly
	)
	c.Status(http.StatusOK)
}

func (uc *userController) LogOut(c *gin.Context) {
	c.SetCookie(
		"token",
		"",
		-1, // MaxAge = -1 means delete cookie
		"/",
		os.Getenv("API_DOMAIN"),
		true,
		true,
	)
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
	token := c.GetHeader("X-CSRF-Token")
	c.JSON(http.StatusOK, gin.H{
		"csrf_token": token,
	})
}

func (uc *userController) UpdateUser(c *gin.Context) {
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	claims := userClaims.(jwt.MapClaims)
	userId := uint(claims["user_id"].(float64))

	user := model.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRes, err := uc.uu.UpdateUser(user, userId)
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
