package handlers

import (
	"context"
	"github.com/eskpil/aarhus/internal/core"

	ticketService "github.com/eskpil/aarhus/internal/core/services/ticket"

	"github.com/eskpil/aarhus/internal/core/state"
	"github.com/eskpil/aarhus/internal/slog"
	"github.com/eskpil/aarhus/pkg/contracts"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func HandleGetTicket(s *state.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(core.UpgradeContext(c), 10*time.Second)
		defer cancel()

		ticketId := c.Param("ticketId")
		ticket, err := ticketService.FindTicket(ctx, s, ticketId)
		if err != nil {
			slog.Errorc(ctx, "could not find ticket", err)
			return c.JSON(http.StatusNotFound, contracts.Error{
				Status:  http.StatusNotFound,
				Message: "could not find ticket",
			})
		}

		return c.JSON(http.StatusOK, ticket)
	}
}

func HandleCreateTicket(s *state.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(core.UpgradeContext(c), 10*time.Second)
		defer cancel()

		input := new(contracts.TicketInput)
		if err := c.Bind(input); err != nil {
			slog.Errorc(ctx, "could not bind input", err)
			return c.JSON(http.StatusBadRequest, contracts.Error{
				Status:  http.StatusBadRequest,
				Message: "invalid json",
			})
		}

		ticket, err := ticketService.Create(ctx, s, input)
		if err != nil {
			slog.Errorc(ctx, "could not create ticket", err)
			return c.JSON(http.StatusInternalServerError, contracts.Error{
				Status:  http.StatusInternalServerError,
				Message: "could not create ticket",
			})
		}

		return c.JSON(http.StatusOK, ticket)
	}
}
