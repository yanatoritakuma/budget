package usecase

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/yanatoritakuma/budget/back/domain/household"
	"github.com/yanatoritakuma/budget/back/domain/user"
	"github.com/yanatoritakuma/budget/back/internal/api"
	"github.com/yanatoritakuma/budget/back/model"
	"github.com/yanatoritakuma/budget/back/utils"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	SignUp(user api.SignUpRequest) (api.UserResponse, error)
	Login(user api.SignUpRequest) (string, error)
	GetLoggedInUser(tokenString string) (*api.UserResponse, error)
	UpdateUser(id uint, req api.UserUpdate) (api.UserResponse, error)
	DeleteUser(id uint) error
	GetHouseholdUsers(userID uint) ([]api.UserResponse, error)
	JoinHousehold(userID uint, inviteCode string) error
	GetOrGenerateCSRFToken(sessionID string) (string, error)
	ValidateCSRFToken(sessionID, token string) bool
}

type userUsecase struct {
	ur         user.IUserRepository
	hr         household.IHouseholdRepository
	tokenStore *model.TokenStore
}

func NewUserUsecase(ur user.IUserRepository, hr household.IHouseholdRepository) IUserUsecase {
	return &userUsecase{
		ur:         ur,
		hr:         hr,
		tokenStore: model.NewTokenStore(),
	}
}

func (uu *userUsecase) SignUp(req api.SignUpRequest) (api.UserResponse, error) {
	ctx := context.Background()

	// Create a new domain household
	domainHousehold, err := household.NewHousehold(
		fmt.Sprintf("%s's Household", req.Name),
		utils.GenerateRandomString(16),
	)
	if err != nil {
		return api.UserResponse{}, err
	}

	if err := uu.hr.Create(ctx, domainHousehold); err != nil {
		return api.UserResponse{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return api.UserResponse{}, err
	}

	domainUser, err := user.NewUser(
		string(req.Email),
		string(hash),
		req.Name,
		"",
		false,
		domainHousehold.ID,
	)
	if err != nil {
		return api.UserResponse{}, err
	}
	if req.Image != nil {
		domainUser.Image = *req.Image
	}

	if err := uu.ur.Create(ctx, domainUser); err != nil {
		return api.UserResponse{}, err
	}

	resUser := api.UserResponse{
		Id:        int(domainUser.ID),
		Email:     openapi_types.Email(domainUser.Email),
		Name:      domainUser.Name,
		Image:     &domainUser.Image,
		Admin:     domainUser.Admin,
		CreatedAt: domainUser.CreatedAt,
	}
	return resUser, nil
}

func (uu *userUsecase) Login(req api.SignUpRequest) (string, error) {
	ctx := context.Background()

	storedUser, err := uu.ur.FindByEmail(ctx, string(req.Email))
	if err != nil {
		return "", err
	}
	if storedUser == nil {
		return "", fmt.Errorf("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(req.Password))
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
	ctx := context.Background()

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

		domainUser, err := uu.ur.FindByID(ctx, uint(userID))
		if err != nil {
			return nil, err
		}
		if domainUser == nil {
			return nil, fmt.Errorf("user not found")
		}

		id := int(domainUser.ID)
		emailStr := domainUser.Email
		email := openapi_types.Email(emailStr)
		name := domainUser.Name
		image := domainUser.Image
		admin := domainUser.Admin
		createdAt := domainUser.CreatedAt

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

func (uu *userUsecase) UpdateUser(id uint, req api.UserUpdate) (api.UserResponse, error) {
	ctx := context.Background()

	existingUser, err := uu.ur.FindByID(ctx, id)
	if err != nil {
		return api.UserResponse{}, err
	}
	if existingUser == nil {
		return api.UserResponse{}, fmt.Errorf("user not found")
	}

	if req.Name != nil {
		existingUser.Name = *req.Name
	}
	if req.Image != nil {
		existingUser.Image = *req.Image
	}

	if err := uu.ur.Update(ctx, existingUser); err != nil {
		return api.UserResponse{}, err
	}

	resUser := api.UserResponse{
		Id:        int(existingUser.ID),
		Email:     openapi_types.Email(existingUser.Email),
		Name:      existingUser.Name,
		Image:     &existingUser.Image,
		Admin:     existingUser.Admin,
		CreatedAt: existingUser.CreatedAt,
	}
	return resUser, nil
}

func (uu *userUsecase) DeleteUser(id uint) error {
	ctx := context.Background()
	if err := uu.ur.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func (uu *userUsecase) GetHouseholdUsers(userID uint) ([]api.UserResponse, error) {
	ctx := context.Background()

	currentUser, err := uu.ur.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}
	if currentUser == nil {
		return nil, fmt.Errorf("current user not found")
	}

	if currentUser.HouseholdID == 0 {
		return nil, fmt.Errorf("current user does not belong to a household")
	}

	householdUsers, err := uu.ur.FindByHouseholdID(ctx, currentUser.HouseholdID)
	if err != nil {
		return nil, fmt.Errorf("failed to get household users: %w", err)
	}

	// Format the response
	var resUsers []api.UserResponse
	for _, domainUser := range householdUsers {
		id := int(domainUser.ID)
		emailStr := domainUser.Email
		email := openapi_types.Email(emailStr)
		name := domainUser.Name
		image := domainUser.Image
		admin := domainUser.Admin
		createdAt := domainUser.CreatedAt

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
	ctx := context.Background()

	domainHousehold, err := uu.hr.FindByInviteCode(ctx, inviteCode) // Use FindByInviteCode
	if err != nil {
		return fmt.Errorf("invalid invite code: %w", err)
	}
	if domainHousehold == nil {
		return fmt.Errorf("invalid invite code: household not found")
	}

	domainUser, err := uu.ur.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("could not find user: %w", err)
	}
	if domainUser == nil {
		return fmt.Errorf("user not found")
	}

	domainUser.HouseholdID = domainHousehold.ID
	if err := uu.ur.Update(ctx, domainUser); err != nil {
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
