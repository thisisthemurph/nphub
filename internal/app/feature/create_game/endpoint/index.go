package endpoint

import (
	"github.com/labstack/echo/v4"
	"nphud/cmd/app/setup/contract"
	"nphud/cmd/app/setup/contract/params"
	"nphud/internal/app/feature/create_game/view"
	"nphud/internal/app/shared/ui"
)

type createGameEndpoint struct {
	params.GameRouteParams
}

func NewCreateGameEndpoint(params *params.GameRouteParams) contract.Endpoint {
	return &createGameEndpoint{
		GameRouteParams: *params,
	}
}

func (ep *createGameEndpoint) MapEndpoint() {
	ep.GameGroup.GET("", ep.indexHandler())
	ep.GameGroup.POST("", ep.createGameHandler())
}

func (ep *createGameEndpoint) indexHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return ui.Render(c, view.CreateGamePage())
	}
}
