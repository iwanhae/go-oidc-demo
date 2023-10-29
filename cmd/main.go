package main

import (
	"context"
	"net/http"

	"github.com/caarlos0/env/v10"
	"github.com/iwanhae/oidc-go-demo/pkg/errors"
	"github.com/iwanhae/oidc-go-demo/pkg/oidc"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type config struct {
	OIDC oidc.Config
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	oidcSvc, err := oidc.NewOIDCService(context.Background(), &cfg.OIDC)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	e.GET("/", func(c echo.Context) error {
		return c.HTML(200, "<a href='/login'>Login</a>")
	})
	e.GET("/login", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, oidcSvc.GetRedirectURL())
	})
	e.GET("/oatuh/redirect", func(c echo.Context) error {
		code := c.QueryParam("code")
		ctx := c.Request().Context()
		if len(code) == 0 {
			return c.NoContent(http.StatusBadRequest)
		}

		tkn, err := oidcSvc.Exchange(ctx, code)
		if err != nil {
			return errors.Wrap(err, "oidc: failed to exchange code")
		}

		userinfo, err := oidcSvc.FetchUserInfo(ctx, tkn)
		if err != nil {
			return errors.Wrap(err, "oidc: failed to fetch user info")
		}

		return c.JSON(http.StatusOK, userinfo)
	})

	e.Start(":8080")
}
