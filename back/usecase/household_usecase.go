package usecase

import (
	"context"
	"fmt"

	"github.com/yanatoritakuma/budget/back/domain/household" // Added
	"github.com/yanatoritakuma/budget/back/domain/user"
	"github.com/yanatoritakuma/budget/back/utils"
)

type HouseholdUsecase interface {
	GenerateInviteCode(userID uint) (string, error)
}

type householdUsecase struct {
	hr household.HouseholdRepository
	ur user.UserRepository
}

func NewHouseholdUsecase(hr household.HouseholdRepository, ur user.UserRepository) HouseholdUsecase { // Changed hr type
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
	inviteCodeStr := utils.GenerateRandomString(household.InviteCodeLength)
	inviteCode, err := household.NewInviteCode(inviteCodeStr)
	if err != nil {
		// This should theoretically not happen if GenerateRandomString is correct
		return "", fmt.Errorf("failed to create a valid invite code: %w", err)
	}

	// Save the code to the household
	domainHousehold.GenerateNewInviteCode(inviteCode)          // Use domain method
	if err := hu.hr.Update(ctx, domainHousehold); err != nil { // Use Update
		return "", fmt.Errorf("could not save invite code: %w", err)
	}

	return inviteCode.Value(), nil
}
