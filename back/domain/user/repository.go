package user

import (
	"context"
)

// IUserRepository defines the interface for user data operations.
type IUserRepository interface {
	FindByID(ctx context.Context, id uint) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, userEntity *User) error
	Update(ctx context.Context, userEntity *User) error
	Delete(ctx context.Context, id uint) error
	FindByHouseholdID(ctx context.Context, householdID uint) ([]*User, error)
}
