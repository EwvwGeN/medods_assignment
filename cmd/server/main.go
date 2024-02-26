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
	l "github.com/EwvwGeN/medods_assignment/internal/logger"
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

	app := app.ServerNewInstance(mainCtx, *cfg, logger, nil)
	app.RunServer(mainCtx)


	stopChecker := make(chan os.Signal, 1)
	signal.Notify(stopChecker, syscall.SIGTERM, syscall.SIGINT)
	<- stopChecker
	logger.Info("start stopping service")
	cancel()
	logger.Info("service stopped")
}