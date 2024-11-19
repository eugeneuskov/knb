package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"knb/app/entities"
	customErrors "knb/app/errors"
	"knb/app/interfaces"
	"knb/app/repositories"
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
	if err := g.checkUser(playerOwnerId); err != nil {
		return nil, err
	}

	return g.gameRepository.CreateGame(playerOwnerId)
}

func (g *gameService) FindGame(gameId uuid.UUID) (*entities.Game, error) {
	return g.gameRepository.FindById(gameId)
}

func (g *gameService) JoinGame(playerId uuid.UUID, gameId uuid.UUID) (*entities.Game, error) {
	if err := g.checkUser(playerId); err != nil {
		return nil, err
	}
	game, err := g.getGame(gameId)
	if err != nil {
		return nil, err
	}
	for _, player := range game.Players {
		if player.ID == playerId {
			return nil, customErrors.NewBadRequestError("you already joined to this game")
		}
	}

	if err := g.gameRepository.AddPlayers(game, []uuid.UUID{playerId}); err != nil {
		return nil, err
	}

	return g.FindGame(gameId)
}

func (g *gameService) StartGame(gameId uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (g *gameService) checkUser(playerId uuid.UUID) error {
	_, err := g.playerRepository.FindById(playerId)
	if err != nil {
		var notFoundErr *customErrors.NotFoundError
		if errors.As(err, &notFoundErr) {
			return customErrors.NewWrongLoginError("Unauthorized")
		}

		return err
	}

	return nil
}

func (g *gameService) getGame(gameId uuid.UUID) (*entities.Game, error) {
	game, err := g.gameRepository.FindById(gameId)
	if err != nil {
		if err.Error() == repositories.RecordNotFoundError {
			return nil, customErrors.NewNotFoundError(fmt.Sprintf("game with id %s not found", gameId))
		}

		return nil, err
	}

	return game, nil
}
