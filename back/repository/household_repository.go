package repository

import (
	"context"
	"time"

	"github.com/yanatoritakuma/budget/back/domain/household"
	"github.com/yanatoritakuma/budget/back/model"

	"gorm.io/gorm"
)

var _ household.IHouseholdRepository = (*HouseholdRepositoryImpl)(nil)

// HouseholdRepositoryImpl implements domain.HouseholdRepository using GORM.
type HouseholdRepositoryImpl struct {
	db *gorm.DB
}

// NewHouseholdRepositoryImpl creates a new HouseholdRepositoryImpl.
func NewHouseholdRepositoryImpl(db *gorm.DB) *HouseholdRepositoryImpl {
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
	return toDomainHousehold(&householdModel), nil
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
	return toDomainHousehold(&householdModel), nil
}

// Create creates a new household.
func (repo *HouseholdRepositoryImpl) Create(ctx context.Context, householdEntity *household.Household) error {
	householdModel := toModelHousehold(householdEntity)
	if err := repo.db.WithContext(ctx).Create(householdModel).Error; err != nil {
		return err
	}
	householdEntity.ID = householdModel.ID
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

func toDomainHousehold(h *model.Household) *household.Household {
	if h == nil {
		return nil
	}
	return &household.Household{
		ID:         h.ID,
		Name:       h.Name,
		InviteCode: h.InviteCode,
		CreatedAt:  h.CreatedAt,
		UpdatedAt:  h.UpdatedAt,
	}
}

func toModelHousehold(h *household.Household) *model.Household {
	if h == nil {
		return nil
	}
	return &model.Household{
		ID:         h.ID,
		Name:       h.Name,
		InviteCode: h.InviteCode,
		CreatedAt:  h.CreatedAt,
		UpdatedAt:  h.UpdatedAt,
	}
}
