package command

import (
	"context"
	"database/sql"
	"nphud/internal/app/feature/manage_game/model"
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
	stmt := "select number, player_uid, api_key, start_time, tick_rate, production_rate from games where id = ?;"

	var startTimeMillis int64
	err := h.db.QueryRowContext(ctx, stmt, cmd.RowID).Scan(
		&game.Number,
		&game.PlayerUID,
		&game.APIKey,
		&startTimeMillis,
		&game.TickRate,
		&game.ProductionRate,
	)

	game.StartTime = time.Unix(0, startTimeMillis*int64(time.Millisecond))

	if err != nil {
		return game, err
	}

	return game, nil
}
