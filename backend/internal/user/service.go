package user

import (
	"context"
	"fmt"
	"os"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo MongoUserRepository
}

type UserFilter struct {
	Id       string
	Username string
}

func NewService(repo MongoUserRepository) *Service {
	return &Service{repo}
}

func checkPasswordSafety(password string) bool {
	if len(password) < 16 {
		return false
	}

	var (
		hasUppercase = false
		HasLowercase = false
		hasSymbol    = false
		hasDigit     = false
	)

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUppercase = true
		} else if unicode.IsLower(char) {
			HasLowercase = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		} else {
			hasSymbol = true
		}
	}

	return hasUppercase && HasLowercase && hasSymbol && hasDigit
}

func (s *Service) CreateNewUser(ctx context.Context, user *User) error {
	if user.Username == "" {
		return ErrUsernameNotValid
	}

	if !checkPasswordSafety(user.Password) {
		return ErrPasswordNotValid
	}

	_, err := s.repo.GetByUsername(ctx, user.Username)
	if err == nil {
		return ErrUsernameTaken
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return ErrCouldNotHashPwd
	}

	user.Password = string(passwordHash)

	id, err := s.repo.Create(ctx, user)
	if err != nil {
		return ErrCreatingUser
	}

	fmt.Println("created user with id: ", id)
	return nil
}

func (s *Service) LoginUser(ctx context.Context, username string, pwd string) (*string, error) {
	foundUser, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if !foundUser.Enabled {
		return nil, ErrUserNotAuthorized
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(pwd))
	if err != nil {
		return nil, ErrPasswordNotValid
	}

	tokenGenerator := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  foundUser.ID.Hex(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := tokenGenerator.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return nil, ErrTokenGeneration
	}

	return &token, nil
}

func (s *Service) GetUserById(ctx context.Context, id string) (*User, error) {
	foundUser, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return foundUser, nil
}
