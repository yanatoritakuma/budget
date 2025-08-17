package repository

import (
	"github.com/yanatoritakuma/budget/back/model"
	"gorm.io/gorm"
)

type IHouseholdRepository interface {
	CreateHousehold(household *model.Household) error
	GetHousehold(household *model.Household, householdID uint) error
	UpdateHousehold(household *model.Household) error
	GetHouseholdByInviteCode(household *model.Household, inviteCode string) error
}

type householdRepository struct {
	db *gorm.DB
}

func NewHouseholdRepository(db *gorm.DB) IHouseholdRepository {
	return &householdRepository{db}
}

func (hr *householdRepository) CreateHousehold(household *model.Household) error {
	if err := hr.db.Create(household).Error; err != nil {
		return err
	}
	return nil
}

func (hr *householdRepository) GetHousehold(household *model.Household, householdID uint) error {
	if err := hr.db.First(household, householdID).Error; err != nil {
		return err
	}
	return nil
}

func (hr *householdRepository) UpdateHousehold(household *model.Household) error {
	if err := hr.db.Save(household).Error; err != nil {
		return err
	}
	return nil
}

func (hr *householdRepository) GetHouseholdByInviteCode(household *model.Household, inviteCode string) error {
	if err := hr.db.Where("invite_code = ?", inviteCode).First(household).Error; err != nil {
		return err
	}
	return nil
}