package tests

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	btcapi "github.com/A-Danylevych/btc-api"
	"github.com/A-Danylevych/btc-api/pkg/handler"
	"github.com/A-Danylevych/btc-api/pkg/repository"
	"github.com/A-Danylevych/btc-api/pkg/service"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

const (
	fileName = "test.json"
	apiUrl   = "https://api.coinbase.com/v2/prices/spot?currency=UAH"
	port     = "8000"
)

type APITestSuite struct {
	suite.Suite

	handlers *handler.Handler
	services *service.Service
	repos    *repository.Repository
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func (s *APITestSuite) SetupSuite() {
	_, err := os.Create(fileName)
	if err != nil {
		s.FailNow("File system error ")
	}

	if err := s.populateDB(); err != nil {
		s.FailNow("Failed to populate DB", err)
	}

	s.initDeps()
}

func (s *APITestSuite) TearDownSuite() {
	os.Remove(fileName)
}

func (s *APITestSuite) initDeps() {
	log := logrus.New()
	s.repos = repository.NewRepository(fileName)
	s.services = service.NewService(s.repos, apiUrl)
	s.handlers = handler.NewHandler(s.services, log)
}

func (s *APITestSuite) populateDB() error {

	users := []btcapi.User{
		{
			Id:       1,
			Email:    "already@registered.mail",
			Password: "12345",
		},
		{
			Id:       2,
			Email:    "second@test.mail",
			Password: "12345",
		},
		{
			Id:       3,
			Email:    "third@test.mail",
			Password: "12345",
		},
	}

	for index, user := range users {
		users[index].Password = generatePasswordHash(user.Password)
	}
	dataBytes, err := json.Marshal(users)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fileName, dataBytes, 0644)

	return err
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	const salt = "qwertysjfwqheqssdkp8763"
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
