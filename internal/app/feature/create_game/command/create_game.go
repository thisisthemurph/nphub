package command

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"nphud/internal/shared/service"
	"nphud/pkg/np"
	"nphud/pkg/np/model"
	"nphud/pkg/util"
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

	gameRowId, err := c.upsertGameInDatabase(ctx, cmd.Number, cmd.APIKey, snapshotFileName, snapshot.ScanningData)
	if err != nil {
		return CreateGameResult{}, err
	}

	return CreateGameResult{GameID: gameRowId, SnapshotFileName: snapshotFileName}, nil
}

func (c *CreateGameCommandHandler) upsertGameInDatabase(ctx context.Context, gameNumber, apiKey, snapshotFileName string, scanning model.ScanningData) (int64, error) {
	var stmt string
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("could not begin transaction", "err", err)
		return 0, err
	}
	defer tx.Rollback()

	// Determine if the game already exists for the player
	var existingGameID int64
	var existingGameApiKey string
	stmt = `select id, api_key from games where number = ? and player_uid = ?;`
	err = c.db.QueryRow(stmt, gameNumber, scanning.PlayerUID).Scan(&existingGameID, &existingGameApiKey)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, err
		}
	}

	// Upsert the new games row
	var gameRowId int64
	gameExists := existingGameApiKey != ""
	if gameExists {
		err = updateExistingGameRow(tx, gameNumber, apiKey, existingGameApiKey)
		gameRowId = existingGameID
	} else {
		gameRowId, err = insertNewGameRow(tx, gameNumber, apiKey, scanning)
	}
	if err != nil {
		return 0, err
	}

	// Insert the new snapshots row
	stmt = `insert into snapshots (game_id, path, created_at) values (?, ?, ?);`
	if _, err = tx.Exec(stmt, gameRowId, snapshotFileName, time.Now().UnixMilli()); err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		slog.Error("could not commit transaction", "err", err)
		return 0, err
	}

	return gameRowId, nil
}

func insertNewGameRow(
	tx *sql.Tx,
	gameNumber string,
	apiKey string,
	scanning model.ScanningData,
) (int64, error) {
	nextTickTime, err := np.CalculateNextTickTime(scanning.StartTime, scanning.TickRate)
	if err != nil {
		return 0, err
	}

	stmt := `
	insert into games (
		name,
		number,
		api_key,
		player_uid,
		start_time,
		tick_rate,
		production_rate,
		started,
		paused,
		game_over,
		next_snapshot_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	res, err := tx.Exec(
		stmt,
		scanning.Name,
		gameNumber,
		apiKey,
		scanning.PlayerUID,
		scanning.StartTimeRaw,
		scanning.TickRate,
		scanning.ProductionRate,
		util.BoolToInt(scanning.Started),
		util.BoolToInt(scanning.Paused),
		util.BoolToInt(scanning.GameOver),
		nextTickTime.UnixMilli(),
	)

	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func updateExistingGameRow(tx *sql.Tx, gameNumber, apiKey, existingGameApiKey string) error {
	stmt := `update games set api_key = ? where number = ? and api_key = ?;`
	_, err := tx.Exec(stmt, apiKey, gameNumber, existingGameApiKey)
	return err
}
