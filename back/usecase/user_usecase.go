package usecase

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/yanatoritakuma/budget/back/internal/api"
	"github.com/yanatoritakuma/budget/back/model"
	"github.com/yanatoritakuma/budget/back/repository"
	"github.com/yanatoritakuma/budget/back/utils"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	SignUp(user api.SignUpRequest) (api.UserResponse, error)
	Login(user api.SignUpRequest) (string, error)
	GetLoggedInUser(tokenString string) (*api.UserResponse, error)
	UpdateUser(user model.User, id uint) (api.UserResponse, error)
	DeleteUser(id uint) error
	GetHouseholdUsers(userID uint) ([]api.UserResponse, error)
	JoinHousehold(userID uint, inviteCode string) error
	GetOrGenerateCSRFToken(sessionID string) (string, error)
	ValidateCSRFToken(sessionID, token string) bool
}

type userUsecase struct {
	ur         repository.IUserRepository
	hr         repository.IHouseholdRepository
	tokenStore *model.TokenStore
}

func NewUserUsecase(ur repository.IUserRepository, hr repository.IHouseholdRepository) IUserUsecase {
	return &userUsecase{
		ur:         ur,
		hr:         hr,
		tokenStore: model.NewTokenStore(),
	}
}

func (uu *userUsecase) SignUp(user api.SignUpRequest) (api.UserResponse, error) {
	// Create a new household for the user
	newHousehold := model.Household{
		Name:       fmt.Sprintf("%s's Household", user.Name),
		InviteCode: utils.GenerateRandomString(16), // Generate an initial invite code
	}
	if err := uu.hr.CreateHousehold(&newHousehold); err != nil {
		return api.UserResponse{}, err
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		// Here, we should probably delete the household that was just created to avoid orphaned data.
		// For simplicity in this step, we'll omit that. In a production system, this should be a transaction.
		return api.UserResponse{}, err
	}

	// Create the new user with the household ID
	newUser := model.User{
		Email:       string(user.Email),
		Password:    string(hash),
		Name:        user.Name,
		HouseholdID: newHousehold.ID,
	}
	if user.Image != nil {
		newUser.Image = *user.Image
	}
	if err := uu.ur.CreateUser(&newUser); err != nil {
		return api.UserResponse{}, err
	}

	id := int(newUser.ID)
	emailStr := newUser.Email
	email := openapi_types.Email(emailStr)
	name := newUser.Name
	image := newUser.Image
	admin := newUser.Admin
	createdAt := newUser.CreatedAt

	resUser := api.UserResponse{
		Id:        id,
		Email:     email,
		Name:      name,
		Image:     &image,
		Admin:     admin,
		CreatedAt: createdAt,
	}
	return resUser, nil
}

func (uu *userUsecase) Login(user api.SignUpRequest) (string, error) {
	storedUser := model.User{}
	if err := uu.ur.GetUserByEmail(&storedUser, string(user.Email)); err != nil {
		return "", err
	}
	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.ID,
		"exp":     time.Now().Add(time.Hour * 12).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (uu *userUsecase) GetLoggedInUser(tokenString string) (*api.UserResponse, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid user ID in JWT token")
		}
		user := model.User{}
		err = uu.ur.GetUserByID(&user, uint(userID))
		if err != nil {
			return nil, err
		}

		id := int(user.ID)
		emailStr := user.Email
		email := openapi_types.Email(emailStr)
		name := user.Name
		image := user.Image
		admin := user.Admin
		createdAt := user.CreatedAt

		return &api.UserResponse{
			Id:        id,
			Email:     email,
			Name:      name,
			Image:     &image,
			Admin:     admin,
			CreatedAt: createdAt,
		}, nil
	} else {
		return nil, fmt.Errorf("invalid JWT token")
	}
}

func (uu *userUsecase) UpdateUser(user model.User, id uint) (api.UserResponse, error) {
	if err := uu.ur.UpdateUser(&user, id); err != nil {
		return api.UserResponse{}, err
	}

	uid := int(user.ID)
	emailStr := user.Email
	email := openapi_types.Email(emailStr)
	name := user.Name
	image := user.Image
	admin := user.Admin
	createdAt := user.CreatedAt

	resUser := api.UserResponse{
		Id:        uid,
		Email:     email,
		Name:      name,
		Image:     &image,
		Admin:     admin,
		CreatedAt: createdAt,
	}
	return resUser, nil
}

func (uu *userUsecase) DeleteUser(id uint) error {
	if err := uu.ur.DeleteUser(id); err != nil {
		return err
	}
	return nil
}

func (uu *userUsecase) GetHouseholdUsers(userID uint) ([]api.UserResponse, error) {
	// Get the current user to find their household ID
	var currentUser model.User
	if err := uu.ur.GetUserByID(&currentUser, userID); err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}

	if currentUser.HouseholdID == 0 {
		return nil, fmt.Errorf("current user does not belong to a household")
	}

	// Get all users from that household
	var householdUsers []model.User
	if err := uu.ur.GetUsersByHouseholdID(&householdUsers, currentUser.HouseholdID); err != nil {
		return nil, fmt.Errorf("failed to get household users: %w", err)
	}

	// Format the response
	var resUsers []api.UserResponse
	for _, user := range householdUsers {
		id := int(user.ID)
		emailStr := user.Email
		email := openapi_types.Email(emailStr)
		name := user.Name
		image := user.Image
		admin := user.Admin
		createdAt := user.CreatedAt

		resUsers = append(resUsers, api.UserResponse{
			Id:        id,
			Email:     email,
			Name:      name,
			Image:     &image,
			Admin:     admin,
			CreatedAt: createdAt,
		})
	}

	return resUsers, nil
}

func (uu *userUsecase) JoinHousehold(userID uint, inviteCode string) error {
	// Find the household by invite code
	var household model.Household
	if err := uu.hr.GetHouseholdByInviteCode(&household, inviteCode); err != nil {
		return fmt.Errorf("invalid invite code: %w", err)
	}

	// Get the current user
	var user model.User
	if err := uu.ur.GetUserByID(&user, userID); err != nil {
		return fmt.Errorf("could not find user: %w", err)
	}

	// Update user's household
	user.HouseholdID = household.ID
	if err := uu.ur.UpdateUser(&user, userID); err != nil {
		return fmt.Errorf("failed to update user's household: %w", err)
	}

	return nil
}

// generateCSRFToken は新しいCSRFトークンを生成し、保存します（内部メソッド）
func (uu *userUsecase) generateCSRFToken(sessionID string) (string, error) {
	token := utils.GenerateRandomString(32)
	uu.tokenStore.SaveToken(sessionID, model.CSRFToken{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24時間有効
	})
	return token, nil
}

// GetOrGenerateCSRFToken は既存のトークンを返すか、新しいトークンを生成します
func (uu *userUsecase) GetOrGenerateCSRFToken(sessionID string) (string, error) {
	// 既存のトークンを確認
	if token, exists := uu.tokenStore.GetToken(sessionID); exists {
		return token, nil
	}

	// 新しいトークンを生成
	return uu.generateCSRFToken(sessionID)
}

// ValidateCSRFToken はCSRFトークンを検証します
func (uu *userUsecase) ValidateCSRFToken(sessionID, token string) bool {
	return uu.tokenStore.ValidateToken(sessionID, token)
}
