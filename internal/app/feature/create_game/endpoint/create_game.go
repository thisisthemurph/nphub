package endpoint

import (
	"errors"
	"fmt"
	"nphud/internal/app/feature/create_game/command"
	"nphud/internal/app/feature/create_game/view"
	"nphud/internal/app/shared/ui"
	"nphud/pkg/np"
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
	// TODO: If the game exists but for a game with the current number and player uid, update the key
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		req := &CreateGameRequest{}
		if err := c.Bind(req); err != nil {
			return err
		}

		form := view.NewCreateGameFormProps(
			strings.TrimSpace(req.GameNumber),
			strings.TrimSpace(req.APIKey),
		)

		if !form.Validate() {
			return ui.Render(c, view.CreateGameForm(form))
		}

		game := np.New(form.Number, form.Key)
		if !game.Validate() {
			// The given number and api key combination is not correct.
			form.Errors.Set("key", "The given game number and API key does not match any current games.")
			return ui.Render(c, view.CreateGameForm(form))
		}

		cmd := command.NewCreateGameCommand(game.Number, game.APIKey)
		_, err := mediatr.Send[*command.CreateGameCommand, command.CreateGameResult](ctx, cmd)
		if err != nil {
			return err
		}

		redirectTarget := fmt.Sprintf("game/%s/%s", game.Number, game.APIKey)
		c.Response().Header().Set("HX-Redirect", redirectTarget)
		return nil
	}
}
