package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func (s *APITestSuite) TestCreateNewUser() {
	router := s.handlers.InitRouters()
	r := s.Require()
	email, password := "to@registe.mail", "12345"

	logInData := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)

	req := httptest.NewRequest("POST", "/user/create", bytes.NewBufferString(logInData))

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		s.FailNow(err.Error())
	}
	logInResponse := struct {
		UserId int `json:"id"`
	}{}
	json.Unmarshal(body, &logInResponse)
	if err != nil {
		s.FailNow(err.Error())
	}

	r.Equal(http.StatusOK, resp.Code)
	expectedUserId := 4
	r.Equal(expectedUserId, logInResponse.UserId)
}

func (s *APITestSuite) TestCreateNewUserWrongMail() {
	router := s.handlers.InitRouters()
	r := s.Require()
	email, password := "toregiste.gr", "12345"

	logInData := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)

	req := httptest.NewRequest("POST", "/user/create", bytes.NewBufferString(logInData))

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		s.FailNow(err.Error())
	}
	logInResponse := struct {
		Message string `json:"message"`
	}{}
	json.Unmarshal(body, &logInResponse)
	if err != nil {
		s.FailNow(err.Error())
	}
	r.Equal(http.StatusInternalServerError, resp.Code)
	expectedMailError := "mail: missing '@' or angle-addr"
	r.Equal(expectedMailError, logInResponse.Message)
}
