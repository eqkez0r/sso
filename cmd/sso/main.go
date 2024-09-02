package main

import (
	"github.com/eqkez0r/sso/internal/app"
	"github.com/eqkez0r/sso/internal/config"
	"github.com/eqkez0r/sso/internal/logger/zap"
)

func main() {
	log := zap.New()

	cfg := config.New(log)

	log.Info("parsing cfg: ", cfg)

	a := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTl)

	a.GRPCSrv.Run()

}
