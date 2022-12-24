package internalgrpc

import (
	"context"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/server"
	pb "github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/server/grpc/google"
	grpcRun "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"net/http"
)

type GRPCServer struct {
	server     *server.Server
	httpServer *http.Server
	grpcServer *grpc.Server
	pb.UnimplementedStorageServer
}

func NewGRPCServer(ctx context.Context, log server.Logger, app server.Application, addr string) *GRPCServer {
	server := &GRPCServer{
		server:     server.New(ctx, log, app, addr),
		grpcServer: grpc.NewServer(),
	}

	mux := grpcRun.NewServeMux()

	if err := pb.RegisterStorageHandlerServer(ctx, mux, server); err != nil {
		server.server.Logger.Fatal(err)
	}

	pb.RegisterStorageServer(server.grpcServer, server)

	server.httpServer = &http.Server{
		Addr:    server.server.Addr,
		Handler: server.server.Middle.LoggingMiddleware(mux, server.server.Logger),
	}

	return server
}

func (s *GRPCServer) Start(ctx context.Context) error {
	if err := s.server.App.Connect(ctx); err != nil {
		s.server.Logger.Error("[ERR] App don`t connect: ", err)
		return err
	}

	return s.httpServer.ListenAndServe()
}

func (s *GRPCServer) Stop(ctx context.Context) error {
	return s.server.App.Close(ctx)
}
