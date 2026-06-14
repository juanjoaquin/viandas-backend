package main

import (
	"context"
	"fmt"

	"github.com/juanjoaquin/viandas-backend/database"
	"github.com/juanjoaquin/viandas-backend/internal/api"
	"github.com/juanjoaquin/viandas-backend/internal/repository"
	"github.com/juanjoaquin/viandas-backend/internal/service"
	"github.com/juanjoaquin/viandas-backend/settings"
	"github.com/labstack/echo/v5"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			context.Background,
			settings.New,
			database.New,
			repository.New,
			service.New,
			api.New,
			echo.New,
		),
		fx.Invoke(setLifeCycle),
	)

	app.Run()
}

func setLifeCycle(lc fx.Lifecycle, a *api.API, s *settings.Settings, e *echo.Echo) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			address := fmt.Sprintf(":%s", s.Port)
			go a.Start(e, address)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
