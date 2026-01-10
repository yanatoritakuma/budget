package household

import "context"

// HouseholdRepository is the interface for persisting Household domain entities.
type HouseholdRepository interface {
	FindByID(ctx context.Context, id uint) (*Household, error)
	FindByInviteCode(ctx context.Context, inviteCode string) (*Household, error)
	Create(ctx context.Context, household *Household) error
	Update(ctx context.Context, household *Household) error
}
