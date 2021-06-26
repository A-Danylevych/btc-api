package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"net/mail"
	"time"

	btcapi "github.com/A-Danylevych/btc-api"
	"github.com/A-Danylevych/btc-api/pkg/repository"
	"github.com/dgrijalva/jwt-go"
)

const (
	salt       = "qwertysjfwqheqssdkp8763"
	signingKey = "sedeskaj212432lsdaealee"
	tokenTTL   = time.Hour
)

type AuthService struct {
	repo repository.Authorization
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

//Verifies that the email format is correct and passes user information to the record, after hashing the password.
//Returns the created user ID or an error
func (a *AuthService) CreateUser(user btcapi.User) (int, error) {
	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return user.Id, err
	}

	user.Password = generatePasswordHash(user.Password)
	return a.repo.CreateUser(user)
}

//Checks if there is such a user. Returns authorization token or error
func (a *AuthService) GenerateToken(user btcapi.User) (string, error) {
	user.Password = generatePasswordHash(user.Password)
	id, err := a.repo.GetUserId(user)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		id,
	})

	return token.SignedString([]byte(signingKey))
}

//The token is parsed in the ID of the user to whom it belongs.
func (a *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

//Hashes the password
func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
