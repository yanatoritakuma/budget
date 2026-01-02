package usecase

import (
	"context"
	"fmt"

	"github.com/yanatoritakuma/budget/back/domain/household" // Added
	"github.com/yanatoritakuma/budget/back/domain/user"
	"github.com/yanatoritakuma/budget/back/utils"
)

type IHouseholdUsecase interface {
	GenerateInviteCode(userID uint) (string, error)
}

type householdUsecase struct {
	hr household.IHouseholdRepository // Changed type
	ur user.IUserRepository
}

func NewHouseholdUsecase(hr household.IHouseholdRepository, ur user.IUserRepository) IHouseholdUsecase { // Changed hr type
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
	domainHousehold, err := hu.hr.FindByID(ctx, domainUser.HouseholdID) // Use FindByID
	if err != nil {
		return "", fmt.Errorf("could not find household: %w", err)
	}
	if domainHousehold == nil {
		return "", fmt.Errorf("household not found")
	}

	// Generate a unique invite code
	inviteCode := utils.GenerateRandomString(16)

	// Save the code to the household
	domainHousehold.GenerateNewInviteCode(inviteCode)          // Use domain method
	if err := hu.hr.Update(ctx, domainHousehold); err != nil { // Use Update
		return "", fmt.Errorf("could not save invite code: %w", err)
	}

	return inviteCode, nil
}
