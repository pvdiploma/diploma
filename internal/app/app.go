package app

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type App struct {
	log  *slog.Logger
	grpc *grpc.Server
	port int
}

func NewApp(log *slog.Logger, grpc *grpc.Server, port int) *App {
	return &App{
		log:  log,
		grpc: grpc,
		port: port,
	}
}

func (a *App) Run() error {
	log := a.log.With("port", a.port)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return err
	}

	log.Info("Starting grpc auth server", slog.String("addr", l.Addr().String()))

	if err := a.grpc.Serve(l); err != nil {
		return err
	}

	return nil
}

func (a *App) Stop() error {
	a.log.Info("Stopping grpc auth server", slog.Int("addr", a.port))
	a.grpc.GracefulStop()
	return nil
}
