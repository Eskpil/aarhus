package handlers

import (
	"context"

	serverService "github.com/eskpil/aarhus/internal/core/services/server"

	"github.com/eskpil/aarhus/internal/core"
	"github.com/eskpil/aarhus/internal/core/state"
	"github.com/eskpil/aarhus/internal/slog"
	"github.com/eskpil/aarhus/pkg/contracts"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func HandleServerCreate(s *state.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(core.UpgradeContext(c), 10*time.Second)
		defer cancel()

		input := new(contracts.ServerCreateInput)
		if err := c.Bind(input); err != nil {
			slog.Errorc(ctx, "could not bind input", err)
			return c.JSON(http.StatusBadRequest, contracts.Error{
				Status:  http.StatusBadRequest,
				Message: "invalid json",
			})
		}

		server, err := serverService.Create(ctx, s, input)
		if err != nil {
			slog.Errorc(ctx, "could not create server", err)
			return c.JSON(http.StatusInternalServerError, contracts.Error{
				Status:  http.StatusInternalServerError,
				Message: "could not create server",
			})
		}

		return c.JSON(http.StatusOK, server)
	}
}

func HandleGetAllServers(s *state.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(core.UpgradeContext(c), 10*time.Second)
		defer cancel()

		servers, err := serverService.GetAll(ctx, s)
		if err != nil {
			slog.Errorc(ctx, "could not get all servers", err)
			return c.JSON(http.StatusInternalServerError, contracts.Error{
				Status:  http.StatusInternalServerError,
				Message: "could not get all servers",
			})
		}

		return c.JSON(http.StatusOK, servers)
	}
}

func HandleGetServer(s *state.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(core.UpgradeContext(c), 10*time.Second)
		defer cancel()

		id := c.Param("serverId")

		server, err := serverService.GetById(ctx, s, id)
		if err != nil {
			slog.Errorc(ctx, "could not get server", err)
			return c.JSON(http.StatusInternalServerError, contracts.Error{
				Status:  http.StatusInternalServerError,
				Message: "could not get server",
			})
		}

		return c.JSON(http.StatusOK, server)
	}
}
