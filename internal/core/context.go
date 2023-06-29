package core

import (
	"context"
	"github.com/labstack/echo/v4"
)

func UpgradeContext(c echo.Context) context.Context {
	ctx := c.Request().Context()
	return context.WithValue(ctx, "identity", c.Get("identity"))
}
