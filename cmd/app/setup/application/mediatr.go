package application

import (
	"database/sql"
	"nphud/cmd/app/setup/bahaviours"
	"nphud/internal/app/feature/create_game/command"
	command2 "nphud/internal/app/feature/manage_game/command"
	"nphud/internal/app/feature/manage_game/model"
	"nphud/internal/shared/service"
	npmodel "nphud/pkg/np/model"

	"github.com/mehdihadeli/go-mediatr"
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

	getGameByNumberAndApiKeyHandler := command2.NewGetGameByNumberAndAPIKeyQueryHandler(db)
	err = mediatr.RegisterRequestHandler[*command2.GetGameByNumberAndAPIKeyQuery, model.Game](getGameByNumberAndApiKeyHandler)
	if err != nil {
		return err
	}

	getGameByRowIdQueryHandler := command2.NewGetGameByRowIDQueryHandler(db)
	err = mediatr.RegisterRequestHandler[*command2.GetGameByExternalIdQuery, model.Game](getGameByRowIdQueryHandler)
	if err != nil {
		return err
	}

	getLatestScanningDataQueryHandler := command2.NewGetLatestScanningDataQueryHandler(db, snapshotFileService)
	err = mediatr.RegisterRequestHandler[*command2.GetLatestScanningDataQuery, npmodel.ScanningData](getLatestScanningDataQueryHandler)
	if err != nil {
		return err
	}

	return nil
}
