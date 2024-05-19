package endpoint

import (
	"errors"
	"fmt"
	"nphud/internal/app/feature/create_game/command"
	"nphud/internal/app/feature/create_game/view"
	"nphud/internal/app/shared/ui"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type CreateGameRequest struct {
	GameNumber string `json:"game_number" form:"game_number"`
	APIKey     string `json:"api_key" form:"api_key"`
}

var (
	ErrGameNumberRequired = errors.New("the game number is required")
	ErrGameAPIKeyRequired = errors.New("the API key is required")
)

func (ep *createGameEndpoint) createGameHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		request := &CreateGameRequest{}
		if err := c.Bind(request); err != nil {
			return err
		}

		form := view.NewCreateGameFormProps(
			strings.TrimSpace(request.GameNumber),
			strings.TrimSpace(request.APIKey),
		)

		if !form.Validate() {
			return ui.Render(c, view.CreateGameForm(form))
		}

		cmd := command.NewCreateGameCommand(request.GameNumber, request.APIKey)
		result, err := mediatr.Send[*command.CreateGameCommand, command.CreateGameResult](ctx, cmd)
		if err != nil {
			return err
		}

		redirectTarget := fmt.Sprintf("game/%d", result.GameID)
		c.Response().Header().Set("HX-Redirect", redirectTarget)
		return nil
	}
}
