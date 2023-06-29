package auth

import (
	"context"
	"github.com/eskpil/aarhus/internal/core"
	"github.com/eskpil/aarhus/internal/core/services/acl"
	usersService "github.com/eskpil/aarhus/internal/core/services/users"
	"github.com/eskpil/aarhus/internal/core/state"
	"github.com/eskpil/aarhus/internal/slog"
	"github.com/eskpil/aarhus/pkg/contracts"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func nodeMiddleware(ctx context.Context, s *state.State, c echo.Context) error {
	nodeToken := c.Request().Header.Get("X-Aarhus-Node-Token")

	node, err := acl.VerifyNodeIdentity(ctx, s, nodeToken)
	if err != nil {
		slog.Errorc(ctx, "could not get node", err)
		return echo.NewHTTPError(http.StatusUnauthorized, contracts.Error{Status: http.StatusUnauthorized, Message: "invalid node"})
	}

	iden := new(core.Identity)

	iden.Type = core.IdentityTypeNode
	iden.Node = node

	c.Set("identity", iden)

	return nil
}

func userMiddleware(ctx context.Context, s *state.State, c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		slog.Errorc(ctx, "could not get session", err)
		return echo.NewHTTPError(http.StatusInternalServerError, contracts.Error{Status: http.StatusUnauthorized, Message: "something went wrong"})
	}

	uid := sess.Values["uid"].(string)

	user, err := usersService.FindById(ctx, s, uid)
	if err != nil {
		slog.Errorc(ctx, "could not find user", err)
		return echo.NewHTTPError(http.StatusInternalServerError, contracts.Error{
			Status:  http.StatusInternalServerError,
			Message: "something went wrong",
		})
	}

	iden := new(core.Identity)

	iden.Type = core.IdentityTypeUser
	iden.User = user

	c.Set("identity", iden)

	return nil
}

func fastPathMiddleware(ctx context.Context, s *state.State, c echo.Context) error {
	value := c.Request().Header.Get("X-Fast-Path")

	// TODO: Seed a devonly account which has access to everything
	if value == "linusen.a@gmail.com" {
		user, err := usersService.FindByEmail(ctx, s, "linusen.a@gmail.com")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, contracts.Error{
				Status:  http.StatusUnauthorized,
				Message: "could not get user",
			})
		}

		iden := new(core.Identity)

		iden.Type = core.IdentityTypeUser
		iden.User = user

		c.Set("identity", iden)
	} else {
		return echo.NewHTTPError(http.StatusUnauthorized, contracts.Error{
			Status:  http.StatusUnauthorized,
			Message: "invalid fast path user",
		})
	}

	return nil
}

func Middleware(s *state.State) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
			defer cancel()

			if viper.GetBool("DEVELOPMENT") && c.Request().Header.Get("X-Fast-Path") != "" {
				if err := fastPathMiddleware(ctx, s, c); err != nil {
					return err
				}

				return next(c)
			}

			if c.Request().Header.Get("X-Aarhus-Node-Token") != "" {
				if err := nodeMiddleware(ctx, s, c); err != nil {
					return err
				}

				return next(c)
			}

			if c.Request().Header.Get("Cookie") != "" {
				if err := userMiddleware(ctx, s, c); err != nil {
					return err
				}

				return next(c)
			}

			return echo.NewHTTPError(http.StatusBadRequest, contracts.Error{
				Status:  http.StatusBadRequest,
				Message: "missing authorization method",
			})
		}
	}
}
