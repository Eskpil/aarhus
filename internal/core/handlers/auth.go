package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eskpil/aarhus/internal/core/services/acl"
	usersService "github.com/eskpil/aarhus/internal/core/services/users"
	"github.com/eskpil/aarhus/internal/core/state"
	"github.com/eskpil/aarhus/internal/slog"
	"github.com/eskpil/aarhus/pkg/contracts"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
	"time"
)

func HandleAuth(s *state.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := url.Values{
			"client_id":     []string{viper.GetString("DISCORD_CLIENT_ID")},
			"response_type": []string{"code"},
			"scope":         []string{"identify email"},
			"redirect_uri":  []string{"http://localhost:8080/v1/auth/callback/"},
			"prompt":        []string{"consent"},
		}

		uri := url.URL{
			Scheme:   "https",
			Host:     "discord.com",
			Path:     "/oauth2/authorize",
			RawQuery: params.Encode(),
		}

		return c.Redirect(http.StatusTemporaryRedirect, uri.String())
	}
}

func HandleAuthCallback(s *state.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
		defer cancel()

		sess, err := session.Get("session", c)
		if err != nil {
			slog.Errorc(ctx, "could not get session", err)
			return c.Redirect(http.StatusTemporaryRedirect, "frontend.error")
		}

		client := resty.New()

		resp, err := client.R().
			SetFormData(map[string]string{
				"client_id":     viper.GetString("DISCORD_CLIENT_ID"),
				"client_secret": viper.GetString("DISCORD_CLIENT_SECRET"),
				"grant_type":    "authorization_code",
				"code":          c.QueryParam("code"),
				"redirect_uri":  "http://localhost:8080/v1/auth/callback/",
			}).
			Post("https://discord.com/api/oauth2/token")
		if err != nil {
			slog.Errorc(ctx, "could not get token", err)
			return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/error?type=error", viper.GetString("FRONTEND")))
		}

		var tokens struct {
			AccessToken string `json:"access_token"`
			TokenType   string `json:"token_type"`
		}

		if err := json.Unmarshal(resp.Body(), &tokens); err != nil {
			slog.Errorc(ctx, "could not unmarshal tokens", err)
			return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/error?type=error", viper.GetString("FRONTEND")))
		}

		resp, err = client.R().
			SetHeader("Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken)).
			Get("https://discord.com/api/v10/users/@me")

		if err != nil {
			slog.Errorc(ctx, "could not get user", err)
			return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/error?type=error", viper.GetString("FRONTEND")))
		}

		user := new(contracts.User)
		if err := json.Unmarshal(resp.Body(), &user); err != nil {
			slog.Errorc(ctx, "could not unmarshal user", err)
			return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/error?type=error", viper.GetString("FRONTEND")))
		}

		if !user.Verified {
			slog.Errorc(ctx, "user not verified", err)
			return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/error?type=notverified", viper.GetString("FRONTEND")))
		}

		ok, err := acl.CheckUserAccess(ctx, s, user)
		if err != nil {
			slog.Errorc(ctx, "could not check user access", err)
			return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/error?type=error", viper.GetString("FRONTEND")))
		}

		if !ok {
			return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/error?type=noaccess", viper.GetString("FRONTEND")))
		}

		if err := usersService.Replace(ctx, s, user); err != nil {
			slog.Errorc(ctx, "could not replace user", err)
			return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/error?type=error", viper.GetString("FRONTEND")))
		}

		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
		}

		sess.Values["uid"] = user.Id
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			slog.Errorc(ctx, "could not save session", err)
			return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/error?type=error", viper.GetString("FRONTEND")))
		}

		return c.Redirect(http.StatusTemporaryRedirect, viper.GetString("FRONTEND"))
	}
}
