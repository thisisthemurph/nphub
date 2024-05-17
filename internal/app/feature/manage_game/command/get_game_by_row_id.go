package command

import (
	"context"
	"database/sql"
	"nphud/internal/app/feature/manage_game/model"
	"nphud/pkg/util"
	"time"
)

type GetGameByRowIDQuery struct {
	RowID int `validate:"required,gte=1"`
}

func NewGetGameByRowIDQuery(rowID int) *GetGameByRowIDQuery {
	return &GetGameByRowIDQuery{
		RowID: rowID,
	}
}

type GetGameByRowIDQueryHandler struct {
	db *sql.DB
}

func NewGetGameByRowIDQueryHandler(db *sql.DB) *GetGameByRowIDQueryHandler {
	return &GetGameByRowIDQueryHandler{
		db: db,
	}
}

func (h *GetGameByRowIDQueryHandler) Handle(ctx context.Context, cmd *GetGameByRowIDQuery) (model.Game, error) {
	var game model.Game
	stmt := `
	select
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
	from games where id = ?;`

	var (
		started         int
		paused          int
		gameOver        int
		startTimeMillis int64
	)

	err := h.db.QueryRowContext(ctx, stmt, cmd.RowID).Scan(
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
		return game, err
	}

	return game, nil
}
