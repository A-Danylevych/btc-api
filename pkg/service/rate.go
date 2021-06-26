package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type RateService struct {
	url string
}

type Response struct {
	Data Data `json:"data"`
}

type Data struct {
	Base     string `json:"base"`
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

//Makes a request for a third-party API. And returns bitcoin to hryvnia rate or error
func (s *RateService) GetRate() (float64, error) {

	req, err := http.NewRequest("GET", s.url, nil)

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
	response := Response{}
	json.Unmarshal(body, &response)

	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(response.Data.Amount, 64)
}

func NewRateService(url string) *RateService {
	return &RateService{url: url}
}
