package endpoint

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"nphud/cmd/app/setup/contract"
	"nphud/cmd/app/setup/contract/params"
	"nphud/internal/app/feature/manage_game/command"
	"nphud/internal/app/feature/manage_game/model"
	"nphud/internal/app/feature/manage_game/view"
	"nphud/internal/app/shared/ui"
	"strconv"
)

type manageGameEndpoint struct {
	params.GameRouteParams
}

func NewManageGameEndpoint(gameRouteParams *params.GameRouteParams) contract.Endpoint {
	return &manageGameEndpoint{
		GameRouteParams: *gameRouteParams,
	}
}

func (ep *manageGameEndpoint) MapEndpoint() {
	ep.GameGroup.GET("/:gameId", ep.handler())
}

func (ep *manageGameEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		gameIdParam := c.Param("gameId")
		gameId, err := strconv.Atoi(gameIdParam)
		if err != nil {
			return err
		}

		cmd := command.NewGetGameByRowIDQuery(gameId)
		gameResult, err := mediatr.Send[*command.GetGameByRowIDQuery, model.Game](c.Request().Context(), cmd)
		if err != nil {
			return err
		}

		return ui.Render(c, view.ManageGamePage(gameResult))
	}
}
