package repository

import (
	"context"
	"time"

	"github.com/yanatoritakuma/budget/back/domain/household"
	"github.com/yanatoritakuma/budget/back/model"

	"gorm.io/gorm"
)

var _ household.HouseholdRepository = (*HouseholdRepositoryImpl)(nil)

// HouseholdRepositoryImpl implements domain.HouseholdRepository using GORM.
type HouseholdRepositoryImpl struct {
	db *gorm.DB
}

// NewHouseholdRepositoryImpl creates a new HouseholdRepositoryImpl.
func NewHouseholdRepositoryImpl(db *gorm.DB) household.HouseholdRepository {
	return &HouseholdRepositoryImpl{db: db}
}

// FindByID finds a household by ID.
func (repo *HouseholdRepositoryImpl) FindByID(ctx context.Context, id uint) (*household.Household, error) {
	var householdModel model.Household
	if err := repo.db.WithContext(ctx).First(&householdModel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Household not found
		}
		return nil, err
	}
	return toDomainHousehold(&householdModel)
}

// FindByInviteCode finds a household by invite code.
func (repo *HouseholdRepositoryImpl) FindByInviteCode(ctx context.Context, inviteCode string) (*household.Household, error) {
	var householdModel model.Household
	if err := repo.db.WithContext(ctx).Where("invite_code = ?", inviteCode).First(&householdModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Household not found
		}
		return nil, err
	}
	return toDomainHousehold(&householdModel)
}

// Create creates a new household.
func (repo *HouseholdRepositoryImpl) Create(ctx context.Context, householdEntity *household.Household) error {
	householdModel := toModelHousehold(householdEntity)
	if err := repo.db.WithContext(ctx).Create(householdModel).Error; err != nil {
		return err
	}
	householdEntity.ID = household.HouseholdID(householdModel.ID)
	householdEntity.CreatedAt = householdModel.CreatedAt
	householdEntity.UpdatedAt = householdModel.UpdatedAt
	return nil
}

// Update updates an existing household.
func (repo *HouseholdRepositoryImpl) Update(ctx context.Context, householdEntity *household.Household) error {
	householdModel := toModelHousehold(householdEntity)
	householdModel.UpdatedAt = time.Now() // Ensure updated_at is current
	return repo.db.WithContext(ctx).Save(householdModel).Error
}

func toDomainHousehold(h *model.Household) (*household.Household, error) {
	if h == nil {
		return nil, nil
	}

	name, err := household.NewName(h.Name)
	if err != nil {
		return nil, err
	}
	inviteCode, err := household.NewInviteCode(h.InviteCode)
	if err != nil {
		return nil, err
	}

	return &household.Household{
		ID:         household.HouseholdID(h.ID),
		Name:       name,
		InviteCode: inviteCode,
		CreatedAt:  h.CreatedAt,
		UpdatedAt:  h.UpdatedAt,
	}, nil
}

func toModelHousehold(h *household.Household) *model.Household {
	if h == nil {
		return nil
	}
	return &model.Household{
		ID:         h.ID.Value(),
		Name:       h.Name.Value(),
		InviteCode: h.InviteCode.Value(),
		CreatedAt:  h.CreatedAt,
		UpdatedAt:  h.UpdatedAt,
	}
}
