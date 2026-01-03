package repository

import (
	"context"
	"time"

	"github.com/yanatoritakuma/budget/back/domain/user"
	"github.com/yanatoritakuma/budget/back/model"

	"gorm.io/gorm"
)

var _ user.IUserRepository = (*UserRepositoryImpl)(nil)

// UserRepositoryImpl implements domain.UserRepository using GORM.
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepositoryImpl creates a new UserRepositoryImpl.
func NewUserRepositoryImpl(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

// FindByID finds a user by ID.
func (repo *UserRepositoryImpl) FindByID(ctx context.Context, id uint) (*user.User, error) {
	var userModel model.User
	if err := repo.db.WithContext(ctx).First(&userModel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // User not found
		}
		return nil, err
	}
	return toDomainUser(&userModel)
}

// FindByEmail finds a user by email.
func (repo *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var userModel model.User
	if err := repo.db.WithContext(ctx).Where("email = ?", email).First(&userModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // User not found
		}
		return nil, err
	}
	return toDomainUser(&userModel)
}

// Create creates a new user.
func (repo *UserRepositoryImpl) Create(ctx context.Context, userEntity *user.User) error {
	userModel := toModelUser(userEntity)
	if err := repo.db.WithContext(ctx).Create(userModel).Error; err != nil {
		return err
	}
	// Update the domain entity with the generated ID and timestamps
	userEntity.ID = user.UserID(userModel.ID)
	userEntity.CreatedAt = userModel.CreatedAt
	userEntity.UpdatedAt = userModel.UpdatedAt
	return nil
}

// Update updates an existing user.
func (repo *UserRepositoryImpl) Update(ctx context.Context, userEntity *user.User) error {
	userModel := toModelUser(userEntity)
	userModel.UpdatedAt = time.Now() // Ensure updated_at is current
	return repo.db.WithContext(ctx).Save(userModel).Error
}

// Delete deletes a user by ID.
func (repo *UserRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return repo.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

// FindByHouseholdID finds users by household ID.
func (repo *UserRepositoryImpl) FindByHouseholdID(ctx context.Context, householdID uint) ([]*user.User, error) {
	var userModels []model.User
	if err := repo.db.WithContext(ctx).Where("household_id = ?", householdID).Find(&userModels).Error; err != nil {
		return nil, err
	}

	var domainUsers []*user.User
	for i := range userModels {
		domainUser, err := toDomainUser(&userModels[i])
		if err != nil {
			// In a real app, you might want to log this error but continue,
			// or handle it more gracefully depending on business requirements.
			return nil, err
		}
		domainUsers = append(domainUsers, domainUser)
	}
	return domainUsers, nil
}

// toDomainUser converts a model.User to a domain.User
func toDomainUser(userModel *model.User) (*user.User, error) {
	if userModel == nil {
		return nil, nil
	}

	email, err := user.NewEmail(userModel.Email)
	if err != nil {
		return nil, err
	}
	password, err := user.NewPassword(userModel.Password)
	if err != nil {
		return nil, err
	}
	name, err := user.NewName(userModel.Name)
	if err != nil {
		return nil, err
	}

	return &user.User{
		ID:          user.UserID(userModel.ID),
		Email:       email,
		Password:    password,
		Name:        name,
		Image:       userModel.Image,
		Admin:       userModel.Admin,
		CreatedAt:   userModel.CreatedAt,
		UpdatedAt:   userModel.UpdatedAt,
		HouseholdID: userModel.HouseholdID,
	}, nil
}

// toModelUser converts a domain.User to a model.User
func toModelUser(userEntity *user.User) *model.User {
	if userEntity == nil {
		return nil
	}
	return &model.User{
		ID:          userEntity.ID.Value(),
		Email:       userEntity.Email.Value(),
		Password:    userEntity.Password.Value(),
		Name:        userEntity.Name.Value(),
		Image:       userEntity.Image,
		Admin:       userEntity.Admin,
		CreatedAt:   userEntity.CreatedAt,
		UpdatedAt:   userEntity.UpdatedAt,
		HouseholdID: userEntity.HouseholdID,
	}
}
