package handlers

import (
	"context"
	"github.com/eskpil/aarhus/internal/core"

	nodeService "github.com/eskpil/aarhus/internal/core/services/node"

	"github.com/eskpil/aarhus/internal/core/state"
	"github.com/eskpil/aarhus/internal/slog"
	"github.com/eskpil/aarhus/pkg/contracts"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func HandleHeartbeatPost(s *state.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(core.UpgradeContext(c), 10*time.Second)
		defer cancel()

		iden := c.Get("identity").(*core.Identity)
		if iden.Type != core.IdentityTypeNode {
			return c.JSON(http.StatusBadRequest, contracts.Error{
				Status:  http.StatusBadRequest,
				Message: "only intended for nodes",
			})
		}

		input := new(contracts.HeartbeatInput)
		if err := c.Bind(input); err != nil {
			slog.Errorc(ctx, "could not bind input", err)
			return c.JSON(http.StatusBadRequest, contracts.Error{
				Status:  http.StatusBadRequest,
				Message: "invalid json",
			})
		}

		result, err := nodeService.HandleHeartbeat(ctx, s, input)
		if err != nil {
			slog.Errorc(ctx, "could not handle heartbeat", err)
			return c.JSON(http.StatusInternalServerError, contracts.Error{
				Status:  http.StatusInternalServerError,
				Message: "something went wrong, could not handle heartbeat",
			})
		}

		return c.JSON(http.StatusOK, result)
	}
}
