package usecase

import (
	"fmt"

	"github.com/yanatoritakuma/budget/back/model"
	"github.com/yanatoritakuma/budget/back/repository"
	"github.com/yanatoritakuma/budget/back/utils"
)

type IHouseholdUsecase interface {
	GenerateInviteCode(userID uint) (string, error)
}

type householdUsecase struct {
	hr repository.IHouseholdRepository
	ur repository.IUserRepository
}

func NewHouseholdUsecase(hr repository.IHouseholdRepository, ur repository.IUserRepository) IHouseholdUsecase {
	return &householdUsecase{hr, ur}
}

func (hu *householdUsecase) GenerateInviteCode(userID uint) (string, error) {
	// Get user to find their household
	var user model.User
	if err := hu.ur.GetUserByID(&user, userID); err != nil {
		return "", fmt.Errorf("could not find user: %w", err)
	}

	// Get the household
	var household model.Household
	if err := hu.hr.GetHousehold(&household, user.HouseholdID); err != nil {
		return "", fmt.Errorf("could not find household: %w", err)
	}

	// Generate a unique invite code
	inviteCode := utils.GenerateRandomString(16)

	// Save the code to the household
	household.InviteCode = inviteCode
	if err := hu.hr.UpdateHousehold(&household); err != nil {
		return "", fmt.Errorf("could not save invite code: %w", err)
	}

	return inviteCode, nil
}