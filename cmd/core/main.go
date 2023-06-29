package main

import (
	"context"
	"net/http"

	"github.com/eskpil/aarhus/internal/core/middleware/auth"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/spf13/viper"

	"github.com/eskpil/aarhus/internal/core/handlers"
	"github.com/eskpil/aarhus/internal/core/state"

	"github.com/eskpil/aarhus/internal/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	router := echo.New()

	viper.AutomaticEnv()

	router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			slog.Info("request", slog.String("uri", v.URI), slog.Int("status", v.Status))
			return nil
		},
	}))

	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://labstack.com", "https://labstack.net", viper.GetString("FRONTEND")},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	s, err := state.New(context.Background())
	if err != nil {
		slog.Fatal("could not initialize state", err)
	}

	router.Use(session.Middleware(sessions.NewCookieStore([]byte(viper.GetString("SESSION_SECRET")))))

	authGroup := router.Group("/v1/auth/")
	{
		authGroup.GET("", handlers.HandleAuth(s))
		authGroup.GET("callback/", handlers.HandleAuthCallback(s))
	}

	restGroup := router.Group("/v1/")
	{
		restGroup.Use(auth.Middleware(s))

		restGroup.POST("nodes/", handlers.HandleNodeCreate(s))
		restGroup.GET("nodes/", handlers.HandleGetAllNodes(s))
		restGroup.GET("nodes/:nodeId", handlers.HandleGetNode(s))

		restGroup.POST("servers/", handlers.HandleServerCreate(s))
		restGroup.GET("servers/", handlers.HandleGetAllServers(s))
		restGroup.GET("servers/:serverId", handlers.HandleGetServer(s))

		restGroup.GET("tickets/:ticketId", handlers.HandleGetTicket(s))
		restGroup.POST("tickets/", handlers.HandleCreateTicket(s))

		restGroup.POST("heartbeat/", handlers.HandleHeartbeatPost(s))

		restGroup.GET("@me/", handlers.HandleGetMe(s))

		restGroup.PUT("@me/tasks/:taskId/status/", handlers.HandlePutNodeTaskStatus(s))
	}

	router.HideBanner = true
	router.HidePort = true

	slog.Info("starting server", slog.Int("port", 8080))
	if err := router.Start(":8080"); err != nil {
		slog.Fatal("could not start router", err)
	}
}
