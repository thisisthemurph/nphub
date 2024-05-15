package application

import (
	"github.com/labstack/echo/v4"
	"nphud/cmd/app/setup/contract"
)

func (app *Application) MapEndpoints() {
	endpoints := app.Container.Get("routes").([]contract.Endpoint)
	for _, endpoint := range endpoints {
		endpoint.MapEndpoint()
	}

	e := app.Container.Get("echo").(*echo.Echo)
	e.Static("/public", "internal/app/shared/static")
}
