package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	btcapi "github.com/A-Danylevych/btc-api"
	"github.com/A-Danylevych/btc-api/pkg/service"
	mock_service "github.com/A-Danylevych/btc-api/pkg/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestHandlerCreate(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user btcapi.User)
	log := logrus.New()

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            btcapi.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email":"Test","password":"test"}`,
			inputUser: btcapi.User{
				Email:    "Test",
				Password: "test",
			},
			mockBehavior: func(s *mock_service.MockAuthorization,
				user btcapi.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:      "Empty Fields",
			inputBody: `{"email":"Test"}`,
			mockBehavior: func(s *mock_service.MockAuthorization,
				user btcapi.User) {
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"email":"Test","password":"test"}`,
			inputUser: btcapi.User{
				Email:    "Test",
				Password: "test",
			},
			mockBehavior: func(s *mock_service.MockAuthorization,
				user btcapi.User) {
				s.EXPECT().CreateUser(user).Return(1, errors.New("service failure"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			service := &service.Service{Authorization: auth}
			handler := NewHandler(service, log)

			r := gin.New()
			r.POST("/create", handler.create)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/create",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandlerLogIn(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user btcapi.User)
	log := logrus.New()

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            btcapi.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email":"Test","password":"test"}`,
			inputUser: btcapi.User{
				Email:    "Test",
				Password: "test",
			},
			mockBehavior: func(s *mock_service.MockAuthorization,
				user btcapi.User) {
				s.EXPECT().GenerateToken(user).Return("1", nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"token":"1"}`,
		},
		{
			name:      "Empty Fields",
			inputBody: `{"email":"Test"}`,
			mockBehavior: func(s *mock_service.MockAuthorization,
				user btcapi.User) {
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"email":"Test","password":"test"}`,
			inputUser: btcapi.User{
				Email:    "Test",
				Password: "test",
			},
			mockBehavior: func(s *mock_service.MockAuthorization,
				user btcapi.User) {
				s.EXPECT().GenerateToken(user).Return("1", errors.New("service failure"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			service := &service.Service{Authorization: auth}
			handler := NewHandler(service, log)

			r := gin.New()
			r.POST("/login", handler.logIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/login",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
