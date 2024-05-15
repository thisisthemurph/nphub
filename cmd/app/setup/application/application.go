package application

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di"
	"log"
	"log/slog"
	"net/http"
	"nphud/pkg/config"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Application struct {
	Container di.Container
	Echo      *echo.Echo
	Logger    *slog.Logger
	Cfg       *config.AppConfig
}

func New(
	container di.Container,
	e *echo.Echo,
	logger *slog.Logger,
	cfg *config.AppConfig,
) *Application {
	return &Application{
		Container: container,
		Echo:      e,
		Logger:    logger,
		Cfg:       cfg,
	}
}

func (app *Application) Run() {
	defaultDuration := time.Second * 20

	startCtx, cancel := context.WithTimeout(context.Background(), defaultDuration)
	defer cancel()
	app.Start(startCtx)
	<-app.Wait()

	stopCtx, cancel := context.WithTimeout(context.Background(), defaultDuration)
	defer cancel()
	app.Stop(stopCtx)
}

func (app *Application) Start(startCtx context.Context) {
	echoStartHook(startCtx, app)
}

func (app *Application) Stop(shutdownCtx context.Context) {
	echoStopHook(shutdownCtx, app)

	app.Container.Delete()
	log.Println("Graceful shutdown complete.")
}

func (app *Application) Wait() <-chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	return sigChan
}

func echoStartHook(startCtx context.Context, application *Application) {
	go func() {
		// When Shutdown is called, Serve, ListenAndServe, and ListenAndServeTLS immediately return ErrServerClosed. Make sure the program doesn't exit and waits instead for Shutdown to return.
		if err := application.Echo.Start(application.Cfg.ListenAddress); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()
}

func echoStopHook(stopCtx context.Context, application *Application) {
	if err := application.Echo.Shutdown(stopCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
}
