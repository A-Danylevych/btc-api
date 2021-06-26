package repository

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	btcapi "github.com/A-Danylevych/btc-api"
)

type AuthJson struct {
	filePath string
}

func NewAuthJson(filePath string) *AuthJson {
	return &AuthJson{filePath: filePath}
}

//Adding a user to the end of the json file, assigning an id and checking if the user is already registered
func (a *AuthJson) CreateUser(user btcapi.User) (int, error) {

	file, err := ioutil.ReadFile(a.filePath)
	if err != nil {
		return user.Id, err
	}

	users := []btcapi.User{}
	json.Unmarshal(file, &users)

	for _, u := range users {
		if user.Email == u.Email {
			return user.Id, errors.New("email is already registered")
		}
	}
	user.Id = len(users) + 1
	users = append(users, user)
	dataBytes, err := json.Marshal(users)

	if err != nil {
		return user.Id, err
	}

	err = ioutil.WriteFile(a.filePath, dataBytes, 0644)
	if err != nil {
		return user.Id, err
	}
	return user.Id, nil
}

//returns user id by email and password
func (a *AuthJson) GetUserId(user btcapi.User) (int, error) {
	file, err := ioutil.ReadFile(a.filePath)
	if err != nil {
		return user.Id, err
	}
	users := []btcapi.User{}

	json.Unmarshal(file, &users)

	for _, u := range users {
		if user.Email == u.Email && user.Password == u.Password {
			return u.Id, nil
		}
	}

	return user.Id, errors.New("no such user")
}
