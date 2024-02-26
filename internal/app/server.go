package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/EwvwGeN/medods_assignment/internal/config"
	v1 "github.com/EwvwGeN/medods_assignment/internal/http/v1"
	"github.com/gorilla/mux"
)

type server struct {
	log *slog.Logger
	cfg config.Config
	auth Auth
	router *mux.Router
}

type Auth interface {
	RegisterUser()
	LoginUser()
	RefreshToken()
}

func ServerNewInstance(ctx context.Context, cfg config.Config, log *slog.Logger, auth Auth) *server {
	return &server{
		log: log,
		cfg: cfg,
		auth: auth,
		router: mux.NewRouter(),
	}
}

func (s *server) RunServer(ctx context.Context) (errCloseCh chan error) {
	s.log.Info("starting server")
	s.configureRouter()
	errCloseCh = make(chan error)
	srv := &http.Server{
		Handler: s.router,
		Addr:    fmt.Sprintf("%s:%s", s.cfg.HttpConfig.Host, s.cfg.HttpConfig.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	s.log.Info("starting listening", slog.String("addres", srv.Addr))
	go func() {
		for {
			select {
			case <-ctx.Done():
				s.log.Info("Graceful shutdown server")
				errCloseCh <- srv.Shutdown(ctx)
				return
			}
		}
	}()
	srv.ListenAndServe()
	return
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/api/healthchecker", v1.Healthcheck(s.cfg.HttpConfig.PingTimeout))
}