package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"knb/app/handlers/responses"
	"knb/tests/fixtures"
	"net/http"
	"testing"
)

const (
	gameNewGameUrl = "/game/new"
)

type gameNewGameTestCase struct {
	*expectedError
	headers []*testRequestHeader
	name    string
}

func TestGameNewGame(t *testing.T) {
	layers := preparationForTest(t)

	if err := fixtures.NewFixtures(layers.db, layers.service).LoadPlayersFixture(); err != nil {
		t.Errorf("Failed to load fixtures, %s", err)
	}

	playerUuid, err := uuid.NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	nonExistingAuthToken, err := layers.service.Security.GenerateAuthToken(playerUuid)
	if err != nil {
		t.Fatal(err)
	}

	gameNewGameFailedTestCases := []gameNewGameTestCase{
		{
			expectedError: &expectedError{
				code:    http.StatusUnauthorized,
				message: "Unauthorized",
			},
			headers: []*testRequestHeader{
				{
					key:   authorizationToken,
					value: nonExistingAuthToken,
				},
			},
			name: "non-existing player",
		},
	}

	for _, tCase := range gameNewGameFailedTestCases {
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

	existingAuthToken, err := layers.service.Security.GenerateAuthToken(uuid.MustParse(fixtures.AlreadyExistUuid))
	if err != nil {
		t.Fatal(err)
	}

	gameNewGameSuccessTestCases := []gameNewGameTestCase{
		{
			expectedError: nil,
			headers: []*testRequestHeader{
				{
					key:   authorizationToken,
					value: existingAuthToken,
				},
			},
			name: "success new game",
		},
	}

	for _, tCase := range gameNewGameSuccessTestCases {
		t.Run(tCase.name, func(tt *testing.T) {
			resBody, resCode := sendRequestAndGetResponse(requestData{
				router:  layers.router,
				headers: tCase.headers,
				method:  http.MethodPost,
				url:     gameNewGameUrl,
			})

			var response responses.GameNewGameResponse
			err := json.Unmarshal(resBody, &response)
			if !assert.NoError(tt, err) {
				t.Errorf("Failed to create game, %s", err)
			}
			assert.Equal(tt, http.StatusCreated, resCode)

			game, err := layers.service.Game.FindGame(response.Id)
			if !assert.NoError(tt, err) {
				t.Errorf("Failed to get created game, %s", err)
			}
			assert.Equal(tt, response.Id, game.ID)
			assert.Equal(tt, 1, len(game.Players))
			assert.Equal(tt, fixtures.AlreadyExistUuid, game.Players[0].ID.String())
		})
	}

	if err := layers.bootstrap.TeardownTestDB(); err != nil {
		t.Errorf("Failed to teardown test DB, %s", err)
	}
}
