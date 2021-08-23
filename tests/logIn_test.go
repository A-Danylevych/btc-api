package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func (s *APITestSuite) TestLogIn() {
	router := s.handlers.InitRouters()
	r := s.Require()
	email, password := "already@registered.mail", "12345"

	logInData := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)

	req := httptest.NewRequest("POST", "/user/login", bytes.NewBufferString(logInData))

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		s.FailNow(err.Error())
	}
	logInResponse := struct {
		Token string `json:"token"`
	}{}
	json.Unmarshal(body, &logInResponse)
	if err != nil {
		s.FailNow(err.Error())
	}
	r.Equal(http.StatusOK, resp.Code)

}

func (s *APITestSuite) TestLogInNotRegistered() {
	router := s.handlers.InitRouters()
	r := s.Require()
	email, password := "not@registered.mail", "12345"

	logInData := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)

	req := httptest.NewRequest("POST", "/user/login", bytes.NewBufferString(logInData))

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
	expectedErrorMessage := "no such user"
	r.Equal(expectedErrorMessage, logInResponse.Message)

}
