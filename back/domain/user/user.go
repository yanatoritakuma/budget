package user

import "time"

type User struct {
	ID          uint
	Email       string
	Password    string
	Name        string
	Image       string
	Admin       bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	HouseholdID uint
}

// NewUser creates a new User domain entity.
func NewUser(email, password, name, image string, admin bool, householdID uint) (*User, error) {
	// TODO: Add validation and business rules here, e.g., password hashing
	return &User{
		Email:       email,
		Password:    password, // Should be hashed before saving
		Name:        name,
		Image:       image,
		Admin:       admin,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		HouseholdID: householdID,
	}, nil
}

// HasAdminPrivileges checks if the user has admin rights.
func (u *User) HasAdminPrivileges() bool {
	return u.Admin
}

// UpdateName updates the user's name.
func (u *User) UpdateName(newName string) {
	u.Name = newName
	u.UpdatedAt = time.Now()
}

// ChangePassword changes the user's password.
func (u *User) ChangePassword(newPassword string) error {
	// TODO: Add password hashing and security logic here
	u.Password = newPassword
	u.UpdatedAt = time.Now()
	return nil
}
