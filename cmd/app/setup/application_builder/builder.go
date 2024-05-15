package application_builder

import (
	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di"
	"log"
	"log/slog"
	"nphud/cmd/app/setup/application"
	"nphud/pkg/config"
	"os"
)

type ApplicationBuilder struct {
	Services *di.Builder
	Logger   *slog.Logger
}

func NewApplicationBuilder() *ApplicationBuilder {
	logger := createLogger()

	builder, err := di.NewBuilder()
	if err != nil {
		log.Fatal(err)
	}

	return &ApplicationBuilder{
		Services: builder,
		Logger:   logger,
	}
}

// Build creates and returns the new application.Application object
func (b *ApplicationBuilder) Build() *application.Application {
	container := b.Services.Build()

	e := container.Get("echo").(*echo.Echo)
	cfg := container.Get("config").(*config.AppConfig)
	return application.New(
		container,
		e,
		b.Logger,
		cfg,
	)
}

func createLogger() *slog.Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return logger
}
