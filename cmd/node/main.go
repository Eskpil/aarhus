package main

import (
	"context"
	"github.com/eskpil/aarhus/internal/node"
	"github.com/eskpil/aarhus/internal/node/handler"
	"github.com/eskpil/aarhus/internal/node/middleware/auth"
	"github.com/eskpil/aarhus/internal/slog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"sync"
	"time"
)

func main() {
	viper.AutomaticEnv()

	wg := new(sync.WaitGroup)

	s, err := node.New()
	if err != nil {
		slog.Fatal("could not create state", err)
	}

	router := echo.New()

	router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			slog.Info("request", slog.String("uri", v.URI), slog.Int("status", v.Status))
			return nil
		},
	}))

	router.Use(auth.Middleware(s))

	router.GET("/v1/socket/", handler.HandleSocket(s))

	router.HideBanner = true
	router.HidePort = true

	slog.Info("starting")

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		if err := router.Start(":8000"); err != nil {
			slog.Fatal("could not start webserver", err)
		}
	}(wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		if err := s.Start(); err != nil {
			slog.Fatal("could not start state", err)
		}
	}(wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for range time.Tick(30 * time.Second) {
			slog.Info("doing heartbeat")
			if err := s.Heartbeat(context.Background()); err != nil {
				slog.Error("could not handle heartbeat", err)
			}
		}
	}(wg)

	slog.Info("doing heartbeat")
	if err := s.Heartbeat(context.Background()); err != nil {
		slog.Error("could not handle heartbeat", err)
	}

	wg.Wait()
}
