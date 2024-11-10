package main

import (
	"github.com/eqkez0r/sso/internal/app"
	"github.com/eqkez0r/sso/internal/config"
	"github.com/eqkez0r/sso/internal/logger/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log := zap.New()

	cfg := config.New(log)

	log.Info("parsing cfg: ", cfg)

	a, err := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTl)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	go a.GRPCSrv.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	sign := <-stop

	log.Info("received signal: ", sign)

	a.GRPCSrv.Stop()

	log.Info("shutting down gracefully")

}
