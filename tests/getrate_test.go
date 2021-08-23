package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *APITestSuite) TestGetRateOK() {

	router := s.handlers.InitRouters()

	r := s.Require()

	token, err := GetToken(router)
	if err != nil {
		s.FailNow(err.Error())
	}
	req := httptest.NewRequest("GET", "/btcRate", nil)

	headerName := "Authorization"
	headerValue := "Bearer " + token
	req.Header.Add(headerName, headerValue)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	r.Equal(http.StatusOK, resp.Code)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.FailNow(err.Error())
	}
	GetRateResponse := struct {
		Id   int     `json:"user_id"`
		Rate float64 `json:"rate"`
	}{}
	json.Unmarshal(body, &GetRateResponse)

	rate, _ := GetRate()

	expectedUserId := 1
	r.Equal(expectedUserId, GetRateResponse.Id)
	r.Equal(rate, GetRateResponse.Rate)

}

func (s *APITestSuite) TestGetRateWhithoutToken() {

	router := s.handlers.InitRouters()

	r := s.Require()

	req := httptest.NewRequest("GET", "/btcRate", nil)

	headerName := "Authorization"
	headerValue := "Bearer "
	req.Header.Add(headerName, headerValue)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	r.Equal(http.StatusUnauthorized, resp.Code)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.FailNow(err.Error())
	}
	GetRateResponse := struct {
		Message string `json:"message"`
	}{}
	json.Unmarshal(body, &GetRateResponse)

	expectedErrorMessage := "token contains an invalid number of segments"

	r.Equal(expectedErrorMessage, GetRateResponse.Message)

}

func GetToken(s *gin.Engine) (string, error) {
	router := s
	email, password := "already@registered.mail", "12345"

	logInData := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)

	req := httptest.NewRequest("POST", "/user/login", bytes.NewBufferString(logInData))

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}
	logInResponse := struct {
		Token string `json:"token"`
	}{}
	json.Unmarshal(body, &logInResponse)

	return logInResponse.Token, err

}

func GetRate() (float64, error) {

	req, err := http.NewRequest("GET", apiUrl, nil)

	if err != nil {
		return 0, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	type Data struct {
		Base     string `json:"base"`
		Currency string `json:"currency"`
		Amount   string `json:"amount"`
	}
	response := struct {
		Data Data `json:"data"`
	}{}

	json.Unmarshal(body, &response)

	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(response.Data.Amount, 64)
}
