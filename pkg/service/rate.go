package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/A-Danylevych/btc-api/pkg/ratemicroservice/endpoints"
)

type RateService struct {
	url string
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
	response := endpoints.GetResponse{}
	json.Unmarshal(body, &response)

	if err != nil {
		return 0, err
	}

	if response.Err != "" {
		return response.Rate, errors.New(response.Err)
	}
	return response.Rate, nil
}

func NewRateService(url string) *RateService {
	return &RateService{url: url}
}
