package handler

import (
	"github.com/eskpil/aarhus/internal/node"
	"github.com/labstack/echo/v4"
)

func HandleSocket(s *node.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		ws, err := s.Upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}

		_ = ws

		return nil
	}
}
