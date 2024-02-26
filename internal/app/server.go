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
	RegisterUser(email string) (uuid string, err error)
	CreateTokenPair(uuid string) (token, refresh string, err error)
	RefreshToken(oldRefresh string) (newToken, newRefresh string, err error)
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
	go srv.ListenAndServe()
	return
}

func (s *server) configureRouter() {
	s.router.HandleFunc(
		"/api/healthchecker",
		v1.Healthcheck(s.cfg.HttpConfig.PingTimeout)).
	Methods(http.MethodGet)

	s.router.HandleFunc(
		"/api/register",
		s.RegisterUser()).
	Methods(http.MethodPost)

	s.router.HandleFunc(
		"/api/createTokenPair/{uuid:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$}",
		s.CreateTokenPair()).
	Methods(http.MethodGet)

	s.router.HandleFunc(
		"/api/refreshToken",
		s.RefrashToken()).
	Methods(http.MethodPost)
}