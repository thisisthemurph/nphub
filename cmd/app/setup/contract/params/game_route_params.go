package params

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"log/slog"
)

type GameRouteParams struct {
	Logger    *slog.Logger
	GameGroup *echo.Group
	Validator *validator.Validate
}
