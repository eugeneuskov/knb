package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type middlewareTestCase struct {
	*expectedError
	headers []*testRequestHeader
	name    string
}

func TestMiddleware(t *testing.T) {
	layers := preparationForTest(t)

	middlewareFailedTestCases := []middlewareTestCase{
		{
			expectedError: &expectedError{
				code:    http.StatusUnauthorized,
				message: "empty 'Access-Token' header",
			},
			headers: nil,
			name:    "empty headers",
		},
		{
			expectedError: &expectedError{
				code:    http.StatusUnauthorized,
				message: "empty 'Access-Token' header",
			},
			headers: []*testRequestHeader{
				{
					key:   authorizationToken,
					value: "",
				},
			},
			name: "empty Access-Token header",
		},
		{
			expectedError: &expectedError{
				code:    http.StatusUnauthorized,
				message: "invalid 'Access-Token' header",
			},
			headers: []*testRequestHeader{
				{
					key:   authorizationToken,
					value: " ",
				},
			},
			name: "empty string Access-Token header",
		},
		{
			expectedError: &expectedError{
				code:    http.StatusUnauthorized,
				message: "Unauthorized",
			},
			headers: []*testRequestHeader{
				{
					key:   authorizationToken,
					value: "wrong-access-token",
				},
			},
			name: "wrong Access-Token header",
		},
	}

	for _, tCase := range middlewareFailedTestCases {
		t.Run(tCase.name, func(tt *testing.T) {
			resBody, resCode := sendRequestAndGetResponse(requestData{
				router:  layers.router,
				headers: tCase.headers,
				method:  http.MethodPost,
				url:     gameNewGameUrl,
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

	playerUuid, err := uuid.NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	successAuthToken, err := layers.service.Security.GenerateAuthToken(playerUuid)
	if err != nil {
		t.Fatal(err)
	}

	middlewareSuccessTestCases := []middlewareTestCase{
		{
			headers: []*testRequestHeader{
				{
					key:   authorizationToken,
					value: successAuthToken,
				},
			},
			name: "success",
		},
	}

	for _, tCase := range middlewareSuccessTestCases {
		t.Run(tCase.name, func(tt *testing.T) {
			_, _ = sendRequestAndGetResponse(requestData{
				router:  layers.router,
				headers: tCase.headers,
				method:  http.MethodPost,
				url:     authLoginUrl,
			})

			playerId, err := layers.service.Security.ParseAuthToken(successAuthToken)
			if err != nil {
				t.Errorf("Failed to parse player token, %s", err)
			}
			assert.Equal(tt, playerUuid.String(), playerId)
		})
	}
}
