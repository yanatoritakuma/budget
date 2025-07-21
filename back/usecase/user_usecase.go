package usecase

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yanatoritakuma/budget/back/model"
	"github.com/yanatoritakuma/budget/back/repository"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	SignUp(user model.User) (model.UserResponse, error)
	Login(user model.User) (string, error)
	GetLoggedInUser(tokenString string) (*model.UserResponse, error)
	UpdateUser(user model.User, id uint) (model.UserResponse, error)
	DeleteUser(id uint) error
}

type userUsecase struct {
	ur repository.IUserRepository
}

func NewUserUsecase(ur repository.IUserRepository) IUserUsecase {
	return &userUsecase{ur}
}

func (uu *userUsecase) SignUp(user model.User) (model.UserResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return model.UserResponse{}, err
	}
	newUser := model.User{Email: user.Email, Password: string(hash), Name: user.Name, Image: user.Image}
	if err := uu.ur.CreateUser(&newUser); err != nil {
		return model.UserResponse{}, err
	}
	resUser := model.UserResponse{
		ID:        newUser.ID,
		Email:     newUser.Email,
		Name:      newUser.Name,
		Image:     newUser.Image,
		Admin:     newUser.Admin,
		CreatedAt: newUser.CreatedAt,
	}
	return resUser, nil
}

func (uu *userUsecase) Login(user model.User) (string, error) {
	storedUser := model.User{}
	if err := uu.ur.GetUserByEmail(&storedUser, user.Email); err != nil {
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

func (uu *userUsecase) GetLoggedInUser(tokenString string) (*model.UserResponse, error) {
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
		return &model.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Image:     user.Image,
			Admin:     user.Admin,
			CreatedAt: user.CreatedAt,
		}, nil
	} else {
		return nil, fmt.Errorf("invalid JWT token")
	}
}

func (uu *userUsecase) UpdateUser(user model.User, id uint) (model.UserResponse, error) {
	if err := uu.ur.UpdateUser(&user, id); err != nil {
		return model.UserResponse{}, err
	}
	resUser := model.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Image:     user.Image,
		Admin:     user.Admin,
		CreatedAt: user.CreatedAt,
	}
	return resUser, nil
}

func (uu *userUsecase) DeleteUser(id uint) error {
	if err := uu.ur.DeleteUser(id); err != nil {
		return err
	}
	return nil
}
