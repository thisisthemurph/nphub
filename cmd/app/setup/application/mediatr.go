package application

import (
	"database/sql"
	"github.com/mehdihadeli/go-mediatr"
	"nphud/cmd/app/setup/bahaviours"
	"nphud/internal/app/feature/create_game/command"
	command2 "nphud/internal/app/feature/manage_game/command"
	"nphud/internal/app/feature/manage_game/model"
	"nphud/internal/shared/service"
)

// ConfigMediator sets up all mediator handlers with appropriate requests.
// All handlers must be manually added here.
func (app *Application) ConfigMediator() error {
	loggerPipeline := &bahaviours.RequestLoggerBehaviour{}
	if err := mediatr.RegisterRequestPipelineBehaviors(loggerPipeline); err != nil {
		return err
	}

	db := app.Container.Get("db").(*sql.DB)
	snapshotFileService := app.Container.Get("snapshot-file-service").(service.SnapshotFileService)

	createGameCommandHandler := command.NewCreateGameCommandHandler(db, snapshotFileService)
	err := mediatr.RegisterRequestHandler[*command.CreateGameCommand, command.CreateGameResult](createGameCommandHandler)
	if err != nil {
		return err
	}

	getGameByRowIdQueryHandler := command2.NewGetGameByRowIDQueryHandler(db)
	err = mediatr.RegisterRequestHandler[*command2.GetGameByRowIDQuery, model.Game](getGameByRowIdQueryHandler)

	return nil
}
