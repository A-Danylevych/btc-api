package repository

import btcapi "github.com/A-Danylevych/btc-api"

type Authorization interface {
	CreateUser(user btcapi.User) (int, error)
	GetUserId(user btcapi.User) (int, error)
}

type Repository struct {
	Authorization
}

func NewRepository(filePath string) *Repository {
	return &Repository{
		Authorization: NewAuthJson(filePath),
	}
}
