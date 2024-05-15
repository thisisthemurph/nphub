package command

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"nphud/internal/shared/service"
	"nphud/pkg/np"
	"nphud/pkg/np/model"
	"time"
)

type CreateGameCommand struct {
	Number string `validate:"required"`
	APIKey string `validate:"required"`
}

func NewCreateGameCommand(number string, apiKey string) *CreateGameCommand {
	return &CreateGameCommand{
		Number: number,
		APIKey: apiKey,
	}
}

type CreateGameResult struct {
	GameID           int64
	SnapshotFileName string
}

type CreateGameCommandHandler struct {
	db                  *sql.DB
	snapshotFileService service.SnapshotFileService
}

func NewCreateGameCommandHandler(db *sql.DB, snapshotFileService service.SnapshotFileService) *CreateGameCommandHandler {
	return &CreateGameCommandHandler{
		db:                  db,
		snapshotFileService: snapshotFileService,
	}
}

func (c *CreateGameCommandHandler) Handle(ctx context.Context, cmd *CreateGameCommand) (CreateGameResult, error) {
	// TODO: Validate the command using validator

	game := np.New(cmd.Number, cmd.APIKey)
	snapshotBytes, err := game.TakeSnapshot()
	if err != nil {
		return CreateGameResult{}, err
	}

	snapshotFileName, err := c.snapshotFileService.Save(cmd.Number, snapshotBytes)
	if err != nil {
		return CreateGameResult{}, err
	}

	var snapshot np.APIResponse
	if err = json.Unmarshal(snapshotBytes, &snapshot); err != nil {
		return CreateGameResult{}, err
	}

	gameRowId, err := c.insertNewGameInDatabase(ctx, cmd.Number, cmd.APIKey, snapshotFileName, snapshot.ScanningData)
	if err != nil {
		return CreateGameResult{}, err
	}

	return CreateGameResult{GameID: gameRowId, SnapshotFileName: snapshotFileName}, nil
}

func (c *CreateGameCommandHandler) insertNewGameInDatabase(ctx context.Context, gameNumber, apiKey, snapshotFileName string, scanning model.ScanningData) (int64, error) {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("could not begin transaction", "err", err)
		return 0, err
	}
	defer tx.Rollback()

	stmt := `
		insert into games (
			number,
			api_key,
			player_uid,
			start_time,
			tick_rate,
			production_rate
		) VALUES (?, ?, ?, ?, ?, ?);`

	res, err := tx.Exec(
		stmt,
		gameNumber,
		apiKey,
		scanning.PlayerUID,
		scanning.StartTimeRaw,
		scanning.TickRate,
		scanning.ProductionRate,
	)

	if err != nil {
		return 0, err
	}

	gameRowId, err := res.LastInsertId()

	stmt = "insert into snapshots (game_id, path, created_at) values (?, ?, ?);"
	if _, err = tx.Exec(stmt, gameRowId, snapshotFileName, time.Now().UnixMilli()); err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		slog.Error("could not commit transaction", "err", err)
		return 0, err
	}

	return gameRowId, nil
}
