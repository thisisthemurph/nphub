package endpoint

import (
	"nphud/cmd/app/setup/contract"
	"nphud/cmd/app/setup/contract/params"
	"nphud/internal/app/feature/manage_game/command"
	"nphud/internal/app/feature/manage_game/model"
	"nphud/internal/app/feature/manage_game/view"
	"nphud/internal/app/shared/ui"
	npmodel "nphud/pkg/np/model"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
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

		grCmd := command.NewGetGameByRowIDQuery(gameId)
		gameResult, err := mediatr.Send[*command.GetGameByRowIDQuery, model.Game](c.Request().Context(), grCmd)
		if err != nil {
			return err
		}

		ssCmd := command.NewGetLatestScanningDataQuery(gameResult.Number, gameResult.PlayerUID)
		snapshotResult, err := mediatr.Send[*command.GetLatestScanningDataQuery, npmodel.ScanningData](c.Request().Context(), ssCmd)
		if err != nil {
			return err
		}

		viewBag := view.ManageGameViewBag{
			GameName:           gameResult.Name,
			GameNumber:         gameResult.Number,
			LatestScanningData: snapshotResult,
		}

		return ui.Render(c, view.ManageGamePage(viewBag))
	}
}
