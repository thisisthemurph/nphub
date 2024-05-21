package endpoint

import (
	"errors"
	"nphud/cmd/app/setup/contract"
	"nphud/cmd/app/setup/contract/params"
	"nphud/internal/app/feature/manage_game/command"
	"nphud/internal/app/feature/manage_game/model"
	"nphud/internal/app/feature/manage_game/view"
	"nphud/internal/app/shared/ui"
	npmodel "nphud/pkg/np/model"
	"strconv"

	"github.com/google/uuid"
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
	ep.GameGroup.GET("/:gameNumber/:gameApiKey", ep.handler())
	ep.GameGroup.GET("/:gameExternalId/player/:playerUid", ep.playerDataHandler())
}

func (ep *manageGameEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		gameNumber := c.Param("gameNumber")
		apiKey := c.Param("gameApiKey")

		grCmd := command.NewGetGameByNumberAndAPIKeyQuery(gameNumber, apiKey)
		gameResult, err := mediatr.Send[*command.GetGameByNumberAndAPIKeyQuery, model.Game](c.Request().Context(), grCmd)
		if err != nil {
			if errors.Is(err, command.ErrGameNotFound) {
				return ui.Render(c, view.GameNotFoundPage())
			}
			return err
		}

		ssCmd := command.NewGetLatestScanningDataQuery(gameResult.Number, gameResult.PlayerUID)
		snapshotResult, err := mediatr.Send[*command.GetLatestScanningDataQuery, npmodel.ScanningData](c.Request().Context(), ssCmd)
		if err != nil {
			return err
		}

		viewBag := view.ManageGameViewBag{
			Game:               gameResult,
			LatestScanningData: snapshotResult,
		}

		return ui.Render(c, view.ManageGamePage(viewBag))
	}
}

func (ep *manageGameEndpoint) playerDataHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		gameExternalId, err := uuid.Parse(c.Param("gameExternalId"))
		if err != nil {
			return err
		}

		playerUID, err := strconv.Atoi(c.Param("playerUid"))
		if err != nil {
			return err
		}

		grCmd := command.NewGetGameByExternalIDQuery(gameExternalId)
		gameResult, err := mediatr.Send[*command.GetGameByExternalIdQuery, model.Game](c.Request().Context(), grCmd)
		if err != nil {
			return err
		}

		ssCmd := command.NewGetLatestScanningDataQuery(gameResult.Number, gameResult.PlayerUID)
		snapshotResult, err := mediatr.Send[*command.GetLatestScanningDataQuery, npmodel.ScanningData](c.Request().Context(), ssCmd)
		if err != nil {
			return err
		}

		player, _ := snapshotResult.Players.Get(playerUID)
		homeStar, _ := snapshotResult.Stars.Get(player.HomeStarUID)

		bag := view.PlayerDataViewBag{
			HomeStar: homeStar,
			Player:   player,
		}

		return ui.Render(c, view.PlayerData(bag))
	}
}
