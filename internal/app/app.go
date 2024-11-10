package app

import (
	grpcapp "github.com/eqkez0r/sso/internal/app/grpc"
	"github.com/eqkez0r/sso/internal/logger"
	"github.com/eqkez0r/sso/internal/services/auth"
	"github.com/eqkez0r/sso/storage/postgres"

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
) (*App, error) {
	storage, err := postgres.New(storagePath)
	if err != nil {
		return nil, err
	}

	authservice := auth.New(log, storage, storage, storage, tokenTTL)
	grpcApp := grpcapp.New(log, grpcPort, authservice)

	return &App{
		GRPCSrv: grpcApp,
	}, nil
}
