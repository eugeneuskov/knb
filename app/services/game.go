package services

import (
	"errors"
	"github.com/google/uuid"
	"knb/app/entities"
	customErrors "knb/app/errors"
	"knb/app/interfaces"
)

type gameService struct {
	gameRepository   interfaces.GameRepository
	playerRepository interfaces.RepositoryPlayer
}

func newGameService(
	gameRepository interfaces.GameRepository,
	playerRepository interfaces.RepositoryPlayer,
) *gameService {
	return &gameService{
		gameRepository,
		playerRepository,
	}
}

func (g *gameService) NewGameRequest(playerOwnerId uuid.UUID) (*entities.Game, error) {
	_, err := g.playerRepository.FindById(playerOwnerId)
	if err != nil {
		var notFoundErr *customErrors.NotFoundError
		if errors.As(err, &notFoundErr) {
			return nil, customErrors.NewWrongLoginError("Unauthorized")
		}

		return nil, err
	}

	return g.gameRepository.CreateGame(playerOwnerId)
}

func (g *gameService) FindGame(gameId uuid.UUID) (*entities.Game, error) {
	return g.gameRepository.FindById(gameId)
}

func (g *gameService) JoinGame(playerId uuid.UUID, gameId uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (g *gameService) StartGame(gameId uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
