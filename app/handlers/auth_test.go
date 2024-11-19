package handlers

import (
	"encoding/json"
	"knb/app/handlers/requests"
	"knb/app/handlers/responses"
	"knb/tests/fixtures"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	authRegistrationUrl = "/auth/registration"
	authLoginUrl        = "/auth/login"
)

type registrationTestCase struct {
	requestBody *requests.AuthRegistrationRequest
	*expectedError
	name string
}

type loginTestCase struct {
	requestBody *requests.AuthLoginRequest
	*expectedError
	name string
}

func TestAuthRegistration(t *testing.T) {
	layers := preparationForTest(t)

	if err := fixtures.NewFixtures(layers.db, layers.service).LoadPlayersFixture(); err != nil {
		t.Errorf("Failed to load fixtures, %s", err)
	}

	registrationFailedTestCases := []registrationTestCase{
		{
			expectedError: &expectedError{
				code:    http.StatusBadRequest,
				message: "Request is empty.",
			},
			name: "empty request body",
		},
		{
			requestBody: &requests.AuthRegistrationRequest{
				Login: "test@test.com",
			},
			expectedError: &expectedError{
				code:    http.StatusBadRequest,
				message: "Field validation for 'Password' failed on the 'required' tag",
			},
			name: "request body without `password` parameter",
		},
		{
			requestBody: &requests.AuthRegistrationRequest{
				Password: "strong-pass",
			},
			expectedError: &expectedError{
				code:    http.StatusBadRequest,
				message: "Field validation for 'Login' failed on the 'required' tag",
			},
			name: "request body without `login` parameter",
		},
		{
			requestBody: &requests.AuthRegistrationRequest{
				Login:    fixtures.AlreadyExistPlayer1Email,
				Password: "strong-pass",
			},
			expectedError: &expectedError{
				code:    http.StatusConflict,
				message: "This login already exists",
			},
			name: "`login` already exists",
		},
	}

	for _, tCase := range registrationFailedTestCases {
		t.Run(tCase.name, func(tt *testing.T) {
			var body []byte
			if tCase.requestBody != nil {
				body, _ = json.Marshal(tCase.requestBody)
			}

			resBody, resCode := sendRequestAndGetResponse(requestData{
				router:      layers.router,
				requestBody: body,
				method:      http.MethodPost,
				url:         authRegistrationUrl,
			})

			if tCase.expectedError != nil {
				if isNotError := assert.NotNil(tt, resBody); !isNotError {
					return
				}

				var resErr responseError
				err := json.Unmarshal(resBody, &resErr)
				if isNotError := assert.NoError(tt, err); !isNotError {
					return
				}
				if isNotError := assert.Equal(tt, tCase.expectedError.code, resCode); !isNotError {
					return
				}
				assert.Equal(tt, tCase.expectedError.message, resErr.Message)
			}
		})
	}

	registrationSuccessTestCases := []registrationTestCase{
		{
			requestBody: &requests.AuthRegistrationRequest{
				Login:    "new-user@test.com",
				Password: fixtures.AlreadyExistPlayer1Password,
			},
			name: "user created",
		},
	}

	for _, tCase := range registrationSuccessTestCases {
		t.Run(tCase.name, func(tt *testing.T) {
			body, _ := json.Marshal(tCase.requestBody)
			resBody, resCode := sendRequestAndGetResponse(requestData{
				router:      layers.router,
				requestBody: body,
				method:      http.MethodPost,
				url:         authRegistrationUrl,
			})

			var response responses.AuthRegistrationResponse
			err := json.Unmarshal(resBody, &response)
			if !assert.NoError(tt, err) {
				t.Errorf("Failed to create player, %s", err)
			}
			assert.Equal(tt, http.StatusCreated, resCode)

			result, err := layers.repository.Player.FindById(response.ID)
			if err != nil {
				t.Errorf("Failed to get created player, %s", err)
			}
			assert.Equal(tt, result.Email, tCase.requestBody.Login)
		})
	}

	if err := layers.bootstrap.TeardownTestDB(); err != nil {
		t.Errorf("Failed to teardown test DB, %s", err)
	}
}

func TestAuthLogin(t *testing.T) {
	layers := preparationForTest(t)

	if err := fixtures.NewFixtures(layers.db, layers.service).LoadPlayersFixture(); err != nil {
		t.Errorf("Failed to load fixtures, %s", err)
	}

	loginFailedTestCases := []loginTestCase{
		{
			expectedError: &expectedError{
				code:    http.StatusBadRequest,
				message: "Request is empty.",
			},
			name: "empty request body",
		},
		{
			requestBody: &requests.AuthLoginRequest{
				Login: "test@test.com",
			},
			expectedError: &expectedError{
				code:    http.StatusBadRequest,
				message: "Field validation for 'Password' failed on the 'required' tag",
			},
			name: "request body without `password` parameter",
		},
		{
			requestBody: &requests.AuthLoginRequest{
				Password: "strong-pass",
			},
			expectedError: &expectedError{
				code:    http.StatusBadRequest,
				message: "Field validation for 'Login' failed on the 'required' tag",
			},
			name: "request body without `login` parameter",
		},
		{
			requestBody: &requests.AuthLoginRequest{
				Login:    "wrong@wrong.com",
				Password: "strong-pass",
			},
			expectedError: &expectedError{
				code:    http.StatusUnauthorized,
				message: "Login or Password are incorrect",
			},
			name: "login is failed",
		},
		{
			requestBody: &requests.AuthLoginRequest{
				Login:    fixtures.AlreadyExistPlayer1Email,
				Password: "wrong-pass",
			},
			expectedError: &expectedError{
				code:    http.StatusUnauthorized,
				message: "Login or Password are incorrect",
			},
			name: "password is failed",
		},
	}

	for _, tCase := range loginFailedTestCases {
		t.Run(tCase.name, func(tt *testing.T) {
			var body []byte
			if tCase.requestBody != nil {
				body, _ = json.Marshal(tCase.requestBody)
			}

			resBody, resCode := sendRequestAndGetResponse(requestData{
				router:      layers.router,
				requestBody: body,
				method:      http.MethodPost,
				url:         authLoginUrl,
			})

			if tCase.expectedError != nil {
				if isNotError := assert.NotNil(tt, resBody); !isNotError {
					return
				}

				var resErr responseError
				err := json.Unmarshal(resBody, &resErr)
				if isNotError := assert.NoError(tt, err); !isNotError {
					return
				}
				if isNotError := assert.Equal(tt, tCase.expectedError.code, resCode); !isNotError {
					return
				}
				assert.Equal(tt, tCase.expectedError.message, resErr.Message)
			}
		})
	}

	loginSuccessTestCases := []loginTestCase{
		{
			requestBody: &requests.AuthLoginRequest{
				Login:    fixtures.AlreadyExistPlayer1Email,
				Password: fixtures.AlreadyExistPlayer1Password,
			},
			name: "login success",
		},
	}

	for _, tCase := range loginSuccessTestCases {
		t.Run(tCase.name, func(tt *testing.T) {
			body, _ := json.Marshal(tCase.requestBody)
			resBody, resCode := sendRequestAndGetResponse(requestData{
				router:      layers.router,
				requestBody: body,
				method:      http.MethodPost,
				url:         authLoginUrl,
			})

			var response responses.AuthLoginResponse
			err := json.Unmarshal(resBody, &response)
			if !assert.NoError(tt, err) {
				t.Errorf("Failed to login player, %s", err)
			}
			assert.Equal(tt, http.StatusOK, resCode)

			playerId, err := layers.service.Security.ParseAuthToken(response.Token)
			if err != nil {
				t.Errorf("Failed to get created player, %s", err)
			}
			assert.Equal(tt, playerId, fixtures.AlreadyExistPlayer1Uuid)
		})
	}

	if err := layers.bootstrap.TeardownTestDB(); err != nil {
		t.Errorf("Failed to teardown test DB, %s", err)
	}
}
