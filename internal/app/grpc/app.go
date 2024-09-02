package grpcapp

import (
	"fmt"
	authgrpc "github.com/eqkez0r/sso/internal/grpc/auth"
	"github.com/eqkez0r/sso/internal/logger"
	"google.golang.org/grpc"
	"net"
)

type App struct {
	log        logger.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log logger.Logger,
	port int,
) *App {
	gRPCServer := grpc.NewServer()
	authgrpc.RegisterServerAPI(gRPCServer)
	return &App{
		log:        log,
		port:       port,
		gRPCServer: gRPCServer,
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Infof("%s: starting gRPC server on port %d", op, a.port)
	if err = a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"
	a.log.Infof("%s: stopping", op)
	a.gRPCServer.GracefulStop()
}
