package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yanatoritakuma/budget/back/usecase"
)

type IHouseholdController interface {
	GenerateInviteCode(c *gin.Context)
}

type householdController struct {
	hu usecase.IHouseholdUsecase
}

func NewHouseholdController(hu usecase.IHouseholdUsecase) IHouseholdController {
	return &householdController{hu}
}

func (hc *householdController) GenerateInviteCode(c *gin.Context) {
	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	claims := userClaims.(jwt.MapClaims)
	userId := uint(claims["user_id"].(float64))

	inviteCode, err := hc.hu.GenerateInviteCode(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"invite_code": inviteCode})
}