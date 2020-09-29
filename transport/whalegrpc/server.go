// Package whalegrpc defines the grpc implementation of the transport.Server transport.
package whalegrpc

import (
	"context"
	"net"

	"github.com/farwydi/cleanwhale/config"
	"github.com/farwydi/cleanwhale/transport"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// NewGRPCServer Creates an grpc transport.
// Configures the main Addr parameter with cfg.
// registerService is called to register user owner grpc service.
// logger is used for startup notification and shutting.
func NewGRPCServer(cfg config.GRPCConfig, logger *zap.Logger, registerService func(srv *grpc.Server)) transport.Server {
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(logger),
		)),
	)

	// Allows to register user service
	registerService(srv)

	reflection.Register(srv)

	return &grpcServer{addr: cfg.Addr, srv: srv, logger: logger}
}

type grpcServer struct {
	addr   string
	srv    *grpc.Server
	logger *zap.Logger
}

// Run Listen tcp server and also expects ctx.Done for graceful shutdown.
func (t *grpcServer) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", t.addr)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()

		t.logger.Info("Shutting down grpc server...")

		t.srv.GracefulStop()
	}()
	t.logger.Info("Starting grpc server",
		zap.String("addr", t.addr))
	return t.srv.Serve(lis)
}
