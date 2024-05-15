package endpoint

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"nphud/internal/app/feature/create_game/command"
)

type CreateGameRequest struct {
	GameNumber string `json:"game_number" form:"game_number"`
	APIKey     string `json:"api_key" form:"api_key"`
}

func (ep *createGameEndpoint) createGameHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		request := &CreateGameRequest{}
		if err := c.Bind(request); err != nil {
			return err
		}

		// TODO: Validate the request

		cmd := command.NewCreateGameCommand(request.GameNumber, request.APIKey)
		result, err := mediatr.Send[*command.CreateGameCommand, command.CreateGameResult](c.Request().Context(), cmd)
		if err != nil {
			return err
		}

		redirectTarget := fmt.Sprintf("game/%d", result.GameID)
		c.Response().Header().Set("HX-Redirect", redirectTarget)
		return nil
	}
}
