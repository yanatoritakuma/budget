package repository

import (
	"fmt"

	"github.com/yanatoritakuma/budget/back/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IUserRepository interface {
	GetUserByEmail(user *model.User, email string) error
	CreateUser(user *model.User) error
	GetUserByID(user *model.User, id uint) error
	UpdateUser(user *model.User, id uint) error
	DeleteUser(id uint) error
	GetUsersByHouseholdID(users *[]model.User, householdID uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db}
}

func (ur *userRepository) GetUserByEmail(user *model.User, email string) error {
	if err := ur.db.Where("email=?", email).First(user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) CreateUser(user *model.User) error {
	if err := ur.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) GetUserByID(user *model.User, id uint) error {
	if err := ur.db.Where("id=?", id).First(user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) UpdateUser(user *model.User, id uint) error {
	result := ur.db.Model(user).Clauses(clause.Returning{}).Where("id=?", id).Updates(map[string]interface{}{
		"email":        user.Email,
		"name":         user.Name,
		"image":        user.Image,
		"household_id": user.HouseholdID,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected < 1 {
		return fmt.Errorf("object does not exist")
	}
	return nil
}

func (ur *userRepository) DeleteUser(id uint) error {
	result := ur.db.Where("id=?", id).Delete(&model.User{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected < 1 {
		return fmt.Errorf("object does not exist")
	}
	return nil
}

func (ur *userRepository) GetUsersByHouseholdID(users *[]model.User, householdID uint) error {
	if err := ur.db.Where("household_id = ?", householdID).Find(users).Error; err != nil {
		return err
	}
	return nil
}
