package service

import (
	"errors"
	"testing"

	btcapi "github.com/A-Danylevych/btc-api"
	mock_repository "github.com/A-Danylevych/btc-api/pkg/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {

	type mockBehavior func(s *mock_repository.MockAuthorization, user btcapi.User)

	testTable := []struct {
		name           string
		inputUser      btcapi.User
		mockBehavior   mockBehavior
		expectedUserId int
		expectedError  error
	}{
		{
			name: "OK",
			inputUser: btcapi.User{
				Email:    "correct@mail.com",
				Password: "test",
			},
			mockBehavior: func(s *mock_repository.MockAuthorization,
				user btcapi.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedUserId: 1,
			expectedError:  nil,
		},
		{
			name: "Invalid mail",
			inputUser: btcapi.User{
				Email:    "notcorrectmail.com",
				Password: "something",
			},
			mockBehavior: func(s *mock_repository.MockAuthorization,
				user btcapi.User) {
			},
			expectedUserId: 0,
			expectedError:  errors.New("mail: missing '@' or angle-addr"),
		},
		{
			name: "Service Failure",
			inputUser: btcapi.User{
				Email:    "correct@mail.com",
				Password: "test",
			},
			mockBehavior: func(s *mock_repository.MockAuthorization,
				user btcapi.User) {
				s.EXPECT().CreateUser(user).Return(0, errors.New("repository error"))
			},
			expectedUserId: 0,
			expectedError:  errors.New("repository error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockAuthorization(c)
			service := NewAuthService(repo)
			inputPassword := testCase.inputUser.Password
			testCase.inputUser.Password = generatePasswordHash(inputPassword)
			testCase.mockBehavior(repo, testCase.inputUser)
			testCase.inputUser.Password = inputPassword

			userId, err := service.CreateUser(testCase.inputUser)

			assert.Equal(t, testCase.expectedUserId, userId)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}
