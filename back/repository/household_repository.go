package repository

import (
	"github.com/yanatoritakuma/budget/back/model"
	"gorm.io/gorm"
)

type IHouseholdRepository interface {
	CreateHousehold(household *model.Household) error
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