package household

import "time"

// Household is the domain entity for a household.
type Household struct {
	ID         HouseholdID
	Name       Name
	InviteCode InviteCode
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewHousehold creates a new Household domain entity.
func NewHousehold(name string, inviteCode string) (*Household, error) {
	voName, err := NewName(name)
	if err != nil {
		return nil, err
	}
	voInviteCode, err := NewInviteCode(inviteCode)
	if err != nil {
		return nil, err
	}

	return &Household{
		Name:       voName,
		InviteCode: voInviteCode,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}

// UpdateName updates the household's name.
func (h *Household) UpdateName(newName Name) {
	h.Name = newName
	h.UpdatedAt = time.Now()
}

// GenerateNewInviteCode generates a new invite code for the household.
func (h *Household) GenerateNewInviteCode(newCode InviteCode) {
	h.InviteCode = newCode
	h.UpdatedAt = time.Now()
}
