package application_builder

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di"
	"log"
	"log/slog"
	"nphud/cmd/app/setup/contract"
	"nphud/cmd/app/setup/contract/params"
	endpoint "nphud/internal/app/feature/create_game/endpoint"
	endpoint2 "nphud/internal/app/feature/manage_game/endpoint"
	"nphud/internal/shared/service"
	"nphud/pkg/config"
	"nphud/pkg/store"
	"os"
)

const (
	ConfigDependencyKey string = "config"
	EchoDependencyKey   string = "echo"
	LoggerDependencyKey string = "log"

	RoutesDependencyKey          string = "routes"
	GameRoutesGroupDependencyKey string = "game_route_group"
)

// AddCore adds the basic dependencies to the application.
func (b *ApplicationBuilder) AddCore() {
	logDep := di.Def{
		Name:  LoggerDependencyKey,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return b.Logger, nil
		},
	}

	configDep := di.Def{
		Name:  ConfigDependencyKey,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return config.NewAppConfig(os.Getenv), nil
		},
	}

	if err := b.Services.Add(logDep); err != nil {
		log.Fatal(err)
	}

	if err := b.Services.Add(configDep); err != nil {
		log.Fatal(err)
	}
}

// AddInfrastructure adds additional libraries and infrastructure to the application.
func (b *ApplicationBuilder) AddInfrastructure() {
	if err := addEcho(b.Services); err != nil {
		log.Fatal(err)
	}

	configDatabase := di.Def{
		Name:  "db",
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get(ConfigDependencyKey).(*config.AppConfig)
			return store.GetOrCreate(cfg.Database.FullPath)
		},
	}

	if err := b.Services.Add(configDatabase); err != nil {
		log.Fatal(err)
	}
}

func (b *ApplicationBuilder) AddServices() {
	snapshotFileServiceDep := di.Def{
		Name:  "snapshot-file-service",
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewSnapshotFileService("snapshots"), nil
		},
	}

	if err := b.Services.Add(snapshotFileServiceDep); err != nil {
		log.Fatal(err)
	}
}

// AddRoutes adds the different route endpoints to the application.
func (b *ApplicationBuilder) AddRoutes() {
	routesDep := di.Def{
		Name:  RoutesDependencyKey,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			gameRouteParams := ctn.Get(GameRoutesGroupDependencyKey).(*params.GameRouteParams)

			manageGameEndpoint := endpoint2.NewManageGameEndpoint(gameRouteParams)
			createGameEndpoint := endpoint.NewCreateGameEndpoint(gameRouteParams)
			endpoints := []contract.Endpoint{
				manageGameEndpoint,
				createGameEndpoint,
			}

			return endpoints, nil
		},
	}

	if err := b.Services.Add(routesDep); err != nil {
		log.Fatal(err)
	}
}

func addEcho(container *di.Builder) error {
	echoDep := di.Def{
		Name:  EchoDependencyKey,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return echo.New(), nil
		},
	}

	gameGroupDep := di.Def{
		Name:  GameRoutesGroupDependencyKey,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			e := ctn.Get(EchoDependencyKey).(*echo.Echo)
			logger := ctn.Get(LoggerDependencyKey).(*slog.Logger)

			games := e.Group("/game")
			gameRouteParams := &params.GameRouteParams{
				Logger:    logger,
				Validator: validator.New(),
				GameGroup: games,
			}

			return gameRouteParams, nil
		},
	}

	if err := container.Add(echoDep); err != nil {
		return err
	}

	if err := container.Add(gameGroupDep); err != nil {
		return err
	}

	return nil
}
