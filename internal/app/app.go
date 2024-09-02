package app

import (
	grpcapp "github.com/eqkez0r/sso/internal/app/grpc"
	"github.com/eqkez0r/sso/internal/logger"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log logger.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {

	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
