package command

import (
	"context"
	"database/sql"
	"errors"
	"nphud/internal/app/feature/manage_game/model"
	"nphud/pkg/util"
	"time"
)

var ErrGameNotFound = errors.New("game not found")

type GetGameByNumberAndAPIKeyQuery struct {
	Number string
	APIKey string
}

func NewGetGameByNumberAndAPIKeyQuery(number, apiKey string) *GetGameByNumberAndAPIKeyQuery {
	return &GetGameByNumberAndAPIKeyQuery{
		Number: number,
		APIKey: apiKey,
	}
}

type GetGameByNumberAndAPIKeyQueryHandler struct {
	db *sql.DB
}

func NewGetGameByNumberAndAPIKeyQueryHandler(db *sql.DB) *GetGameByNumberAndAPIKeyQueryHandler {
	return &GetGameByNumberAndAPIKeyQueryHandler{
		db: db,
	}
}

func (h *GetGameByNumberAndAPIKeyQueryHandler) Handle(ctx context.Context, cmd *GetGameByNumberAndAPIKeyQuery) (model.Game, error) {
	var game model.Game
	stmt := `
	select
		external_id,
		name,
		number, 
		player_uid, 
		api_key, 
		start_time, 
		tick_rate,
		production_rate ,
		started,
		paused,
		game_over
	from games 
	where number = ? and api_key = ?;`

	var (
		started         int
		paused          int
		gameOver        int
		startTimeMillis int64
	)

	err := h.db.QueryRowContext(ctx, stmt, cmd.Number, cmd.APIKey).Scan(
		&game.ExternalId,
		&game.Name,
		&game.Number,
		&game.PlayerUID,
		&game.APIKey,
		&startTimeMillis,
		&game.TickRate,
		&game.ProductionRate,
		&started,
		&paused,
		&gameOver,
	)

	game.Started = util.IntToBool(started)
	game.Paused = util.IntToBool(paused)
	game.GameOver = util.IntToBool(gameOver)
	game.StartTime = time.Unix(0, startTimeMillis*int64(time.Millisecond))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return game, ErrGameNotFound
		}
		return game, err
	}

	return game, nil
}
