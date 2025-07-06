package app

import (
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	httpapp "giveaway-service/internal/app/http"
	"giveaway-service/internal/config"
	giveawayhandler "giveaway-service/internal/handler/giveaway"
	giveawayrepository "giveaway-service/internal/repository/giveaway"
	giveawayservice "giveaway-service/internal/service/giveaway"
	"log/slog"
)

type App struct {
	HttpServer *httpapp.App
}

func New(
	log *slog.Logger,
	cfg *config.Config,
	psql *squirrel.StatementBuilderType,
	pool *pgxpool.Pool,
) *App {
	repository := giveawayrepository.New(log, psql)
	service := giveawayservice.New(log, pool, repository, repository, repository, repository)
	handler := giveawayhandler.New(service)

	httpApp := httpapp.New(
		log,
		cfg.Server.Port,
		*handler,
	)

	return &App{
		HttpServer: httpApp,
	}
}
