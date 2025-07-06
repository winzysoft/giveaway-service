package httpapp

import (
	"context"
	"fmt"
	"giveaway-service/internal/handler/giveaway"
	"log/slog"
	"net/http"
)

type App struct {
	ctx  context.Context // нужен ли context объекту приложения
	log  *slog.Logger
	port int
	//httpServer *http.Server
	mux *http.ServeMux
}

func New(
	log *slog.Logger,
	port int,
	giveawayHandler giveaway.Handler,
) *App {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /giveaway/{id}", giveawayHandler.HandleFindGiveaway)
	mux.HandleFunc("POST /giveaway", giveawayHandler.HandleSaveGiveaway)

	return &App{
		log:  log,
		port: port,
		mux:  mux,
	}
}

// MustRun runs gRPC server and panics if any error occurs.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "httpapp.Run"

	err := http.ListenAndServe(fmt.Sprintf(":%d", a.port), a.mux)
	if err != nil {
		a.log.Error(fmt.Sprintf("%s %s", op, err))
		return err
	}

	return nil
}

// Stop stops gRPC server.
func (a *App) Stop() {
	const op = "httpapp.Stop"
	a.log.With(slog.String("op", op))

	a.log.Info("stopping gRPC server", slog.Int("port", a.port))
	//if err := a.httpServer.Shutdown(a.ctx); err != nil {
	//	a.log.Error(fmt.Sprintf("%s %s", op, err))
	//}
}
