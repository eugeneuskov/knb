package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"knb/app/handlers/requests"
	"knb/app/handlers/responses"
	"knb/app/repositories"
	"knb/app/services"
	"knb/tests"
	"knb/tests/fixtures"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	authTestEnvFilePath = "../../.env.test"

	routesMode      = "debug"
	registrationUrl = "/auth/registration"
	loginUrl        = "/auth/login"
)

type responseError struct {
	Message string `json:"message"`
}

type expectedError struct {
	code    int
	message string
}

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
	bootstrapTest := tests.NewBootstrapTest(authTestEnvFilePath)
	if err := bootstrapTest.SetupTestDB(); err != nil {
		t.Errorf("Failed to setup test DB, %s", err)
	}

	testDb := bootstrapTest.DB()
	repositoryMap := repositories.NewRepository(testDb)
	serviceMap := services.NewService(repositoryMap, bootstrapTest.Config())

	if err := fixtures.NewFixtures(testDb, serviceMap).LoadPlayersFixture(); err != nil {
		t.Errorf("Failed to load fixtures, %s", err)
	}

	router := NewHandler(serviceMap).InitRoutes(routesMode)

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
				Login:    fixtures.AlreadyExistEmail,
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
			resBody, resCode := sendRequestAndGetResponse(
				router,
				tCase.requestBody,
				registrationUrl,
			)

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
				Password: fixtures.AlreadyExistPassword,
			},
			name: "user created",
		},
	}

	for _, tCase := range registrationSuccessTestCases {
		t.Run(tCase.name, func(tt *testing.T) {
			resBody, resCode := sendRequestAndGetResponse(
				router,
				tCase.requestBody,
				registrationUrl,
			)

			var response responses.AuthRegistrationResponse
			err := json.Unmarshal(resBody, &response)
			if !assert.NoError(tt, err) {
				t.Errorf("Failed to create player, %s", err)
			}
			assert.Equal(tt, http.StatusCreated, resCode)

			result, err := repositoryMap.Player.FindById(response.ID)
			if err != nil {
				t.Errorf("Failed to get created player, %s", err)
			}
			assert.Equal(tt, result.Email, tCase.requestBody.Login)
		})
	}

	if err := bootstrapTest.TeardownTestDB(); err != nil {
		t.Errorf("Failed to teardown test DB, %s", err)
	}
}

func TestAuthLogin(t *testing.T) {
	bootstrapTest := tests.NewBootstrapTest(authTestEnvFilePath)
	if err := bootstrapTest.SetupTestDB(); err != nil {
		t.Errorf("Failed to setup test DB, %s", err)
	}

	testDB := bootstrapTest.DB()
	repositoryMap := repositories.NewRepository(testDB)
	serviceMap := services.NewService(repositoryMap, bootstrapTest.Config())

	if err := fixtures.NewFixtures(testDB, serviceMap).LoadPlayersFixture(); err != nil {
		t.Errorf("Failed to load fixtures, %s", err)
	}

	router := NewHandler(serviceMap).InitRoutes(routesMode)

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
				Login:    fixtures.AlreadyExistEmail,
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
			resBody, resCode := sendRequestAndGetResponse(
				router,
				tCase.requestBody,
				loginUrl,
			)

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
				Login:    fixtures.AlreadyExistEmail,
				Password: fixtures.AlreadyExistPassword,
			},
			name: "login success",
		},
	}

	for _, tCase := range loginSuccessTestCases {
		t.Run(tCase.name, func(tt *testing.T) {
			resBody, resCode := sendRequestAndGetResponse(
				router,
				tCase.requestBody,
				loginUrl,
			)

			var response responses.AuthLoginResponse
			err := json.Unmarshal(resBody, &response)
			if !assert.NoError(tt, err) {
				t.Errorf("Failed to login player, %s", err)
			}
			assert.Equal(tt, http.StatusOK, resCode)

			playerId, err := serviceMap.Auth.ParseToken(response.Token)
			if err != nil {
				t.Errorf("Failed to get created player, %s", err)
			}
			assert.Equal(tt, playerId, fixtures.AlreadyExistUuid)
		})
	}

	if err := bootstrapTest.TeardownTestDB(); err != nil {
		t.Errorf("Failed to teardown test DB, %s", err)
	}
}

type authRequest interface {
	*requests.AuthRegistrationRequest |
		*requests.AuthLoginRequest
}

func sendRequestAndGetResponse[R authRequest](
	router *gin.Engine,
	requestBody R,
	url string,
) ([]byte, int) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, url, getRequestBody(requestBody))
	router.ServeHTTP(w, req)

	return w.Body.Bytes(), w.Code
}

func getRequestBody[R authRequest](requestBody R) *bytes.Buffer {
	body, _ := json.Marshal(requestBody)

	switch requestBody {
	case nil:
		return bytes.NewBuffer(nil)
	default:
		return bytes.NewBuffer(body)
	}
}
