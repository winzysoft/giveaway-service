package main

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"giveaway-service/internal/app"
	"giveaway-service/internal/config"
	"giveaway-service/internal/lib/logger/handlers/slogpretty"
	"giveaway-service/internal/repository"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("starting application", slog.Any("config", cfg))

	ctx := context.Background()

	pool := repository.MustInitDB(ctx, log, cfg.Storage.Postgres)
	defer func(pool *pgxpool.Pool) {
		pool.Close()
	}(pool)

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	application := app.New(log, cfg, &psql, pool)
	go application.HttpServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	s := <-stop
	log.Info("stopping application", slog.String("signal", s.String()))
	application.HttpServer.Stop()
	log.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{Level: slog.LevelDebug},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
