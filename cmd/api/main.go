// @title Product API
// @version 1.0
// @description Product REST API with PostgreSQL, Clean Architecture, and partial PATCH.
// @host localhost:8080
// @BasePath /
// @schemes http
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/demo/product-api/docs"
	"github.com/demo/product-api/internal/app"
	"github.com/demo/product-api/internal/config"
)

func main() {
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	application, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatalf("start app: %v", err)
	}
	defer application.Close()

	go func() {
		log.Printf("listening on %s", cfg.HTTPAddr)
		if err := application.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTO)
	defer cancel()
	if err := application.Server.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
