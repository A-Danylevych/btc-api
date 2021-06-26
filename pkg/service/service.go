package service

import (
	btcapi "github.com/A-Danylevych/btc-api"
	"github.com/A-Danylevych/btc-api/pkg/repository"
)

type Authorization interface {
	CreateUser(user btcapi.User) (int, error)
	GenerateToken(user btcapi.User) (string, error)
	ParseToken(token string) (int, error)
}

type BtcRate interface {
	GetRate() (float64, error)
}

type Service struct {
	Authorization
	BtcRate
}

func NewService(repos *repository.Repository, url string) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		BtcRate:       NewRateService(url),
	}
}
