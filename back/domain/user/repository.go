package user

import "context"

// UserRepository is the interface for persisting User domain entities.
type UserRepository interface {
	FindByID(ctx context.Context, id uint) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByHouseholdID(ctx context.Context, householdID uint) ([]*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
}
