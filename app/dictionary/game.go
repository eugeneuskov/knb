package dictionary

type GameStatus string

const (
	GameStatusPlanned  GameStatus = "planned"
	GameStatusWaiting  GameStatus = "waiting"
	GameStatusStarted  GameStatus = "started"
	GameStatusFinished GameStatus = "finished"
)
