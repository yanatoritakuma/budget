package household

import "time"

// Household is the domain entity for a household.
type Household struct {
	ID         uint
	Name       string
	InviteCode string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewHousehold creates a new Household domain entity.
func NewHousehold(name string, inviteCode string) (*Household, error) {
	// TODO: Add validation and business rules here
	return &Household{
		Name:       name,
		InviteCode: inviteCode,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}

// UpdateName updates the household's name.
func (h *Household) UpdateName(newName string) {
	h.Name = newName
	h.UpdatedAt = time.Now()
}

// GenerateNewInviteCode generates a new invite code for the household.
func (h *Household) GenerateNewInviteCode(newCode string) {
	h.InviteCode = newCode
	h.UpdatedAt = time.Now()
}
