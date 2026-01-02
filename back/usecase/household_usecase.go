package usecase

import (
	"context"
	"fmt"

	"github.com/yanatoritakuma/budget/back/domain/user"
	"github.com/yanatoritakuma/budget/back/model"
	"github.com/yanatoritakuma/budget/back/repository"
	"github.com/yanatoritakuma/budget/back/utils"
)

type IHouseholdUsecase interface {
	GenerateInviteCode(userID uint) (string, error)
}

type householdUsecase struct {
	hr repository.IHouseholdRepository
	ur user.UserRepository
}

func NewHouseholdUsecase(hr repository.IHouseholdRepository, ur user.UserRepository) IHouseholdUsecase {
	return &householdUsecase{hr, ur}
}

func (hu *householdUsecase) GenerateInviteCode(userID uint) (string, error) {
	ctx := context.Background()

	// Get user to find their household
	domainUser, err := hu.ur.FindByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("could not find user: %w", err)
	}
	if domainUser == nil {
		return "", fmt.Errorf("user not found")
	}

	// Get the household
	var household model.Household
	if err := hu.hr.GetHousehold(&household, domainUser.HouseholdID); err != nil {
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
