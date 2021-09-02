package rate

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/go-kit/kit/log"
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

func NewRateService(url string) Service {
	return &RateService{url: url}
}

func (s *RateService) GetRate(_ context.Context) (float64, error) {

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

func (w *RateService) ServiceStatus(_ context.Context) (int, error) {
	logger.Log("Checking the Service health...")
	return http.StatusOK, nil
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}
