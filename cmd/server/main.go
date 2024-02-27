package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/EwvwGeN/medods_assignment/internal/app"
	c "github.com/EwvwGeN/medods_assignment/internal/config"
	"github.com/EwvwGeN/medods_assignment/internal/jwt"
	l "github.com/EwvwGeN/medods_assignment/internal/logger"
	"github.com/EwvwGeN/medods_assignment/internal/service"
	"github.com/EwvwGeN/medods_assignment/internal/storage"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "", "path to config file")
}
func main() {
	flag.Parse()
	cfg, err := c.LoadConfig(configPath)
	if err != nil {
		panic(fmt.Sprintf("cant load config from path %s: %s", configPath, err.Error()))
	}

	logger := l.SetupLogger(cfg.LogLevel)

	logger.Info("config loaded")
	logger.Debug("config data", slog.Any("cfg", cfg))
	
	mainCtx, cancel := context.WithCancel(context.Background())

	jwtManager := jwt.NewJwtManager(cfg.JwtSecret)

	mongoDB, err := storage.NewMongoProvider(mainCtx, cfg.MongoConfig)
	if err != nil {
		panic("cant get mongo provider")
	}

	auth := service.NewAuth(mainCtx, logger, mongoDB, jwtManager, cfg.TokenTTL, cfg.RefreshTTL)

	app := app.ServerNewInstance(mainCtx, *cfg, logger, auth)
	app.RunServer(mainCtx)

	stopChecker := make(chan os.Signal, 1)
	signal.Notify(stopChecker, syscall.SIGTERM, syscall.SIGINT)
	<- stopChecker
	logger.Info("start stopping service")
	cancel()
	logger.Info("service stopped")
}