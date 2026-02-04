package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/andreyxaxa/Comment-Tree/config"
	"github.com/andreyxaxa/Comment-Tree/internal/controller/restapi"
	"github.com/andreyxaxa/Comment-Tree/internal/repo/persistent"
	"github.com/andreyxaxa/Comment-Tree/internal/usecase/comment"
	"github.com/andreyxaxa/Comment-Tree/pkg/httpserver"
	"github.com/andreyxaxa/Comment-Tree/pkg/logger"
	"github.com/andreyxaxa/Comment-Tree/pkg/postgres"
)

func Run(cfg *config.Config) {
	// Logger
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Use-Case
	commentUseCase := comment.New(
		persistent.New(pg),
	)

	// HTTP Server
	httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP.Port))
	restapi.NewRouter(httpServer.App, cfg, commentUseCase, l)

	// Start Server
	httpServer.Start()

	// Waiting Signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %v", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %v", err))
	}
}
