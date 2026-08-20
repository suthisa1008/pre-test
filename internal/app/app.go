package app

import (
	"context"
	"net/http"

	"gorm.io/gorm"

	"github.com/demo/product-api/internal/config"
	httpapi "github.com/demo/product-api/internal/handler/http"
	"github.com/demo/product-api/internal/repository/postgres"
	"github.com/demo/product-api/internal/service"
)

type App struct {
	Server *http.Server
	DB     *gorm.DB
}

func New(_ context.Context, cfg config.Config) (*App, error) {
	db, err := postgres.Open(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	repo := postgres.NewProductRepository(db)
	svc := service.NewProductService(repo)
	handler := httpapi.NewProductHandler(svc)

	return &App{
		DB: db,
		Server: &http.Server{
			Addr:    cfg.HTTPAddr,
			Handler: httpapi.NewRouter(handler),
		},
	}, nil
}

func (a *App) Close() {
	_ = postgres.Close(a.DB)
}
