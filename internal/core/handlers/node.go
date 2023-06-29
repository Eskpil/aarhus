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

func HandleNodeCreate(s *state.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(core.UpgradeContext(c), 10*time.Second)
		defer cancel()

		input := new(contracts.NodeCreateInput)
		if err := c.Bind(input); err != nil {
			slog.Errorc(ctx, "could not bind input", err)
			return c.JSON(http.StatusBadRequest, contracts.Error{
				Status:  http.StatusBadRequest,
				Message: "invalid json",
			})
		}

		node, err := nodeService.Create(ctx, s, input)
		if err != nil {
			slog.Errorc(ctx, "could not create node", err)
			return c.JSON(http.StatusInternalServerError, contracts.Error{
				Status:  http.StatusBadRequest,
				Message: "something went wrong",
			})
		}

		return c.JSON(http.StatusOK, node)
	}
}

func HandleGetAllNodes(s *state.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(core.UpgradeContext(c), 10*time.Second)
		defer cancel()

		nodes, err := nodeService.GetAll(ctx, s)
		if err != nil {
			slog.Errorc(ctx, "could not get all nodes", err)
			return c.JSON(http.StatusInternalServerError, contracts.Error{
				Status:  http.StatusInternalServerError,
				Message: "something went wrong",
			})
		}

		return c.JSON(http.StatusOK, nodes)
	}
}

func HandleGetNode(s *state.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(core.UpgradeContext(c), 10*time.Second)
		defer cancel()

		nodeId := c.Param("nodeId")

		node, err := nodeService.GetById(ctx, s, nodeId)
		if err != nil {
			slog.Errorc(ctx, "could not get node", err)
			return c.JSON(http.StatusInternalServerError, contracts.Error{
				Status:  http.StatusInternalServerError,
				Message: "could not get node",
			})
		}

		return c.JSON(http.StatusOK, node)
	}
}
