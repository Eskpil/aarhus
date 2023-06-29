package handlers

import (
	"context"
	"fmt"
	"github.com/eskpil/aarhus/internal/core"
	nodeService "github.com/eskpil/aarhus/internal/core/services/node"
	"github.com/eskpil/aarhus/internal/core/state"
	"github.com/eskpil/aarhus/internal/slog"
	"github.com/eskpil/aarhus/pkg/contracts"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func HandleGetMe(s *state.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(core.UpgradeContext(c), 10*time.Second)
		defer cancel()

		iden := ctx.Value("identity").(*core.Identity)

		switch iden.Type {
		case core.IdentityTypeUser:
			{
				return c.JSON(http.StatusOK, iden.User)
			}
		case core.IdentityTypeNode:
			{
				return c.JSON(http.StatusOK, iden.Node)
			}
		}

		slog.Errorc(ctx, "unknown identity type", fmt.Errorf("?????"), slog.Int("type", int(iden.Type)))
		return c.JSON(http.StatusInternalServerError, contracts.Error{
			Status:  http.StatusInternalServerError,
			Message: "something went wrong",
		})
	}
}

func HandlePutNodeTaskStatus(s *state.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(core.UpgradeContext(c), 10*time.Second)
		defer cancel()

		identity := ctx.Value("identity").(*core.Identity)
		if identity.Type != core.IdentityTypeNode {
			slog.Errorc(ctx, "failed to handle node status change", fmt.Errorf("user tried to access node only endpoint"))
			return c.JSON(http.StatusUnauthorized, contracts.Error{
				Status:  http.StatusUnauthorized,
				Message: "node only",
			})
		}

		body := new(contracts.UpdateTaskStatus)
		if err := c.Bind(&body); err != nil {
			slog.Errorc(ctx, "could not parse body", err)
			return c.JSON(http.StatusBadRequest, contracts.Error{
				Status:  http.StatusBadRequest,
				Message: "invalid json",
			})
		}

		taskId := c.Param("taskId")

		if err := nodeService.UpdateTaskStatus(ctx, s, identity.Node.Id, taskId, body.Status); err != nil {
			slog.Errorc(ctx, "could not update task status", err)
			return c.JSON(http.StatusInternalServerError, contracts.Error{
				Status:  http.StatusInternalServerError,
				Message: "could not update task status",
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"status": "updated"})
	}
}
