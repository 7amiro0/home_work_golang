package internalhttp

import (
	"context"
	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/app"
	pb "github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/server/grps/api"
	grpcGetwey "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"

	"net"
)

type Server struct {
	host, port string
	ctx        context.Context
	app        Application
	muxServer  *grpcGetwey.ServeMux
	listener   net.Listener
	httpServer *http.Server
	pb.UnimplementedStorageServer
}

type Logger interface {
	app.Logger
}

type Application interface {
	app.Storage
}

var l Logger

func NewServer(ctx context.Context, logger Logger, app Application, host, port string) *Server {
	l = logger
	server := grpcGetwey.NewServeMux()

	return &Server{
		host: host,
		port: port,
		ctx:  ctx,
		app:  app,
		muxServer: server,
	}
}

func (s *Server) Start(ctx context.Context) error {
	err := s.app.Connect(ctx)
	if err != nil {
		return err
	}

	router := s.getRouter()
	s.muxServer.
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.app.Close(ctx); err != nil {
		return err
	}
	//return s.httpServer.Shutdown(ctx)
	return s.listener.Close()
}
