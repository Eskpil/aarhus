package auth

import (
	"context"
	"github.com/eskpil/aarhus/internal/node"
	"github.com/eskpil/aarhus/internal/slog"
	"github.com/eskpil/aarhus/pkg/contracts"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func Middleware(s *node.State) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
			defer cancel()

			ticketId := c.QueryParam("ticket")
			if ticketId != "" {
				return echo.NewHTTPError(http.StatusBadRequest, contracts.Error{
					Status:  http.StatusBadRequest,
					Message: "missing ticket",
				})
			}

			ticket, err := s.ValidateTicket(ctx, ticketId)
			if err != nil {
				slog.Errorc(ctx, "could not validate ticket", err)
				return echo.NewHTTPError(http.StatusInternalServerError, contracts.Error{
					Status:  http.StatusInternalServerError,
					Message: "could not validate ticket",
				})
			}

			c.Set("ticket", ticket)

			return nil
		}
	}
}
