package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"knb/app/handlers/responses"
	"knb/tests/fixtures"
	"net/http"
	"testing"
)

const (
	gameNewGameUrl   = "/game/new"
	gameJoinGameUrl  = "/game/join/"
	gameStartGameUrl = "/game/start/"

	nonExistingGameId = "2485e769-aee9-486a-bc66-4ca964d7e617"
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

	existingAuthToken, err := layers.service.Security.GenerateAuthToken(uuid.MustParse(fixtures.Player1Uuid))
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

			game, err := layers.service.Game.FindGame(response.ID)
			if !assert.NoError(tt, err) {
				t.Errorf("Failed to get created game, %s", err)
			}
			assert.Equal(tt, response.ID, game.ID)
			assert.Equal(tt, 1, len(game.Players))
			assert.Equal(tt, fixtures.Player1Uuid, game.Players[0].ID.String())
		})
	}

	if err := layers.bootstrap.TeardownTestDB(); err != nil {
		t.Errorf("Failed to teardown test DB, %s", err)
	}
}

type gameJoinGameTestCase struct {
	*expectedError
	headers []*testRequestHeader
	gameId  string
	name    string
}

func TestGameJoinGame(t *testing.T) {
	layers := preparationForTest(t)

	fixture := fixtures.NewFixtures(layers.db, layers.service)
	if err := fixture.LoadPlayersFixture(); err != nil {
		t.Errorf("Failed to load player fixtures, %s", err)
	}
	if err := fixture.LoadGamesFixture(); err != nil {
		t.Errorf("Failed to load game fixtures, %s", err)
	}

	playerUuid, err := uuid.NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	nonExistingAuthToken, err := layers.service.Security.GenerateAuthToken(playerUuid)
	if err != nil {
		t.Fatal(err)
	}
	existingAuthToken, err := layers.service.Security.GenerateAuthToken(uuid.MustParse(fixtures.Player1Uuid))
	if err != nil {
		t.Fatal(err)
	}

	gameJoinGameFailedTestCases := []gameJoinGameTestCase{
		{
			expectedError: &expectedError{
				code:    http.StatusBadRequest,
				message: "game id is invalid",
			},
			headers: []*testRequestHeader{
				{
					key:   authorizationToken,
					value: nonExistingAuthToken,
				},
			},
			gameId: "wrong-id",
			name:   "id param is invalid",
		},
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
			gameId: nonExistingGameId,
			name:   "non-existing player",
		},
		{
			expectedError: &expectedError{
				code:    http.StatusNotFound,
				message: fmt.Sprintf("game with id %s not found", nonExistingGameId),
			},
			headers: []*testRequestHeader{
				{
					key:   authorizationToken,
					value: existingAuthToken,
				},
			},
			gameId: nonExistingGameId,
			name:   "game not found",
		},
		{
			expectedError: &expectedError{
				code:    http.StatusBadRequest,
				message: "you already joined to this game",
			},
			headers: []*testRequestHeader{
				{
					key:   authorizationToken,
					value: existingAuthToken,
				},
			},
			gameId: fixtures.GamePlannedUuid,
			name:   "existing player",
		},
	}

	for _, tCase := range gameJoinGameFailedTestCases {
		t.Run(tCase.name, func(tt *testing.T) {
			resBody, resCode := sendRequestAndGetResponse(requestData{
				router:  layers.router,
				headers: tCase.headers,
				method:  http.MethodPost,
				url:     fmt.Sprintf("%s%s", gameJoinGameUrl, tCase.gameId),
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

	playerAuthToken, err := layers.service.Security.GenerateAuthToken(uuid.MustParse(fixtures.Player2Uuid))
	if err != nil {
		t.Fatal(err)
	}

	gameJoinGameSuccessTestCases := []gameJoinGameTestCase{
		{
			expectedError: nil,
			headers: []*testRequestHeader{
				{
					key:   authorizationToken,
					value: playerAuthToken,
				},
			},
			gameId: fixtures.GamePlannedUuid,
			name:   "success join to game",
		},
	}

	for _, tCase := range gameJoinGameSuccessTestCases {
		t.Run(tCase.name, func(tt *testing.T) {
			resBody, resCode := sendRequestAndGetResponse(requestData{
				router:  layers.router,
				headers: tCase.headers,
				method:  http.MethodPost,
				url:     fmt.Sprintf("%s%s", gameJoinGameUrl, tCase.gameId),
			})

			var response responses.GameJoinGameResponse
			err := json.Unmarshal(resBody, &response)
			if !assert.NoError(tt, err) {
				t.Errorf("Failed to create game, %s", err)
			}
			assert.Equal(tt, http.StatusOK, resCode)

			game, err := layers.service.Game.FindGame(response.ID)
			if !assert.NoError(tt, err) {
				t.Errorf("Failed to get game, %s", err)
			}
			assert.Equal(tt, response.ID, game.ID)
			assert.Equal(tt, 2, len(game.Players))
			assert.Equal(tt, 2, len(response.Players))
			assert.Equal(tt, fixtures.Player1Uuid, game.Players[0].ID.String())
			assert.Equal(tt, fixtures.Player1Uuid, response.Players[0].ID.String())
			assert.Equal(tt, fixtures.Player1DisplayName, response.Players[0].Name)
			assert.Equal(tt, fixtures.Player2Uuid, game.Players[1].ID.String())
			assert.Equal(tt, fixtures.Player2Uuid, response.Players[1].ID.String())
			assert.Equal(tt, fixtures.Player2DisplayName, response.Players[1].Name)
		})
	}

	if err := layers.bootstrap.TeardownTestDB(); err != nil {
		t.Errorf("Failed to teardown test DB, %s", err)
	}
}

type gameStartGameTestCaseFailed struct {
	*expectedError
	headers []*testRequestHeader
	gameId  string
	name    string
}

type gameStartGameTestCaseSuccess struct {
	checkGameStatus bool
	headers         []*testRequestHeader
	gameId          string
	name            string
}

func TestGameStartGame(t *testing.T) {
	layers := preparationForTest(t)

	fixture := fixtures.NewFixtures(layers.db, layers.service)
	if err := fixture.LoadPlayersFixture(); err != nil {
		t.Errorf("Failed to load player fixtures, %s", err)
	}
	if err := fixture.LoadGamesFixture(); err != nil {
		t.Errorf("Failed to load game fixtures, %s", err)
	}

	playerUuid, err := uuid.NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	nonExistingAuthToken, err := layers.service.Security.GenerateAuthToken(playerUuid)
	if err != nil {
		t.Fatal(err)
	}
	playerOneAuthToken, err := layers.service.Security.GenerateAuthToken(uuid.MustParse(fixtures.Player1Uuid))
	if err != nil {
		t.Fatal(err)
	}
	playerTwoAuthToken, err := layers.service.Security.GenerateAuthToken(uuid.MustParse(fixtures.Player2Uuid))
	if err != nil {
		t.Fatal(err)
	}
	playerThreeAuthToken, err := layers.service.Security.GenerateAuthToken(uuid.MustParse(fixtures.Player3Uuid))
	if err != nil {
		t.Fatal(err)
	}

	startGameFiledTestCases := []gameStartGameTestCaseFailed{
		{
			expectedError: &expectedError{
				code:    http.StatusBadRequest,
				message: "game id is invalid",
			},
			headers: []*testRequestHeader{
				{
					key:   authorizationToken,
					value: nonExistingAuthToken,
				},
			},
			gameId: "wrong-id",
			name:   "id param is invalid",
		},
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
			gameId: nonExistingGameId,
			name:   "non-existing player",
		},
		{
			expectedError: &expectedError{
				code:    http.StatusNotFound,
				message: fmt.Sprintf("game with id %s not found", nonExistingGameId),
			},
			headers: []*testRequestHeader{
				{
					key:   authorizationToken,
					value: playerOneAuthToken,
				},
			},
			gameId: nonExistingGameId,
			name:   "game not found",
		},
		{
			expectedError: &expectedError{
				code:    http.StatusBadRequest,
				message: "the game can't start yet",
			},
			headers: []*testRequestHeader{
				{
					key:   authorizationToken,
					value: playerOneAuthToken,
				},
			},
			gameId: fixtures.GamePlannedUuid,
			name:   "game has scheduled",
		},
		{
			expectedError: &expectedError{
				code:    http.StatusBadRequest,
				message: "the game has already started",
			},
			headers: []*testRequestHeader{
				{
					key:   authorizationToken,
					value: playerTwoAuthToken,
				},
			},
			gameId: fixtures.GameStartedUuid,
			name:   "game already started",
		},
		{
			expectedError: &expectedError{
				code:    http.StatusBadRequest,
				message: "the game already over",
			},
			headers: []*testRequestHeader{
				{
					key:   authorizationToken,
					value: playerTwoAuthToken,
				},
			},
			gameId: fixtures.GameFinishedUuid,
			name:   "game already finished",
		},
		{
			expectedError: &expectedError{
				code:    http.StatusBadRequest,
				message: "not enough players",
			},
			headers: []*testRequestHeader{
				{
					key:   authorizationToken,
					value: playerOneAuthToken,
				},
			},
			gameId: fixtures.GameWaitingUuid,
			name:   "not enough players",
		},
		{
			expectedError: &expectedError{
				code:    http.StatusForbidden,
				message: "you can't participate in this game",
			},
			headers: []*testRequestHeader{
				{
					key:   authorizationToken,
					value: playerThreeAuthToken,
				},
			},
			gameId: fixtures.GameWaitingUuid,
			name:   "start play by non-participant",
		},
		// дальнейшие тест-кейсы необходимо прописать при добавлении условий игры (а именно на max/min игроков)
	}

	for _, tCase := range startGameFiledTestCases {
		t.Run(tCase.name, func(tt *testing.T) {
			resBody, resCode := sendRequestAndGetResponse(requestData{
				router:  layers.router,
				headers: tCase.headers,
				method:  http.MethodPost,
				url:     fmt.Sprintf("%s%s", gameStartGameUrl, tCase.gameId),
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

	startGameSuccessTestCases := []gameStartGameTestCaseSuccess{
		{
			headers: []*testRequestHeader{
				{
					key:   authorizationToken,
					value: playerOneAuthToken,
				},
			},
			gameId: fixtures.GameWaitingUuid,
			name:   "successfully start the game",
		},
		{
			checkGameStatus: true,
			headers: []*testRequestHeader{
				{
					key:   authorizationToken,
					value: playerTwoAuthToken,
				},
			},
			gameId: fixtures.GameWaitingUuid,
			name:   "successfully start the game",
		},
	}

	for _, tCase := range startGameSuccessTestCases {
		t.Run(tCase.name, func(tt *testing.T) {
			_, resCode := sendRequestAndGetResponse(requestData{
				router:  layers.router,
				headers: tCase.headers,
				method:  http.MethodPost,
				url:     fmt.Sprintf("%s%s", gameStartGameUrl, tCase.gameId),
			})

			assert.Equal(tt, http.StatusNoContent, resCode)

			/*
				+ 1. проверить на statusOk
				- 2. проверить статус готовности в game_players
				- 3. при наличии достаточного количества игроков (checkGameStatus=true) проверить смену статуса у самой игры
			*/
		})
	}
}
