package user

import "time"

type User struct {
	ID          UserID
	Email       Email
	Password    Password
	Name        Name
	Image       string
	Admin       bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	HouseholdID uint
}

// NewUser は新しいUserドメインエンティティを生成します。
func NewUser(email, password, name, image string, admin bool, householdID uint) (*User, error) {
	voEmail, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	voPassword, err := NewPassword(password)
	if err != nil {
		return nil, err
	}

	voName, err := NewName(name)
	if err != nil {
		return nil, err
	}

	return &User{
		Email:       voEmail,
		Password:    voPassword,
		Name:        voName,
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
func (u *User) UpdateName(newName Name) {
	u.Name = newName
	u.UpdatedAt = time.Now()
}

// ChangePassword changes the user's password.
func (u *User) ChangePassword(newPassword Password) {
	u.Password = newPassword
	u.UpdatedAt = time.Now()
}
