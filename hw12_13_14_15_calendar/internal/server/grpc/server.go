package internalgrpc

import (
	"context"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/server"
	pb "github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/server/grpc/google"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage"
	grpcRun "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"time"
)

type GRPCServer struct {
	server     *server.Server
	httpServer *http.Server
	grpcServer *grpc.Server
	pb.UnimplementedStorageServer
}

func NewGRPCServer(ctx context.Context, log server.Logger, app server.Application, addr string) (*GRPCServer, error) {
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

	return server, nil
}

func (s *GRPCServer) Start(ctx context.Context) error {
	if err := s.server.App.Connect(ctx); err != nil {
		return err
	}

	return s.httpServer.ListenAndServe()
}

func (s *GRPCServer) Stop(ctx context.Context) error {
	if err := s.server.App.Close(ctx); err != nil {
		return err
	}

	return nil
}

func (s *GRPCServer) Add(ctx context.Context, googleEvent *pb.Event) (*emptypb.Empty, error) {
	//start, err := time.ParseInLocation(server.TimeFormat, googleEvent.GetStart().AsTime().String(), time.UTC)
	//if err != nil {
	//	s.server.Logger.Error(err)
	//	return nil, err
	//}
	//
	//end, err := time.ParseInLocation(server.TimeFormat, googleEvent.GetEnd().AsTime().String(), time.UTC)
	//if err != nil {
	//	s.server.Logger.Error(err)
	//	return nil, err
	//}

	event := storage.Event{
		Title:       googleEvent.GetTitle(),
		UserID:      googleEvent.GetUserID(),
		Description: googleEvent.GetDescription(),
		End:         googleEvent.GetEnd().AsTime(),
		Start:       googleEvent.GetStart().AsTime(),
	}

	return nil, s.server.App.Add(ctx, event)
}

func (s *GRPCServer) Delete(ctx context.Context, googleEvent *pb.Event) (*emptypb.Empty, error) {
	return nil, s.server.App.Delete(ctx, googleEvent.GetId())
}

func (s *GRPCServer) Update(ctx context.Context, googleEvent *pb.Event) (*emptypb.Empty, error) {
	start, err := time.ParseInLocation(server.TimeFormat, googleEvent.GetStart().AsTime().String(), time.UTC)
	if err != nil {
		s.server.Logger.Error(err)
		return nil, err
	}

	end, err := time.ParseInLocation(server.TimeFormat, googleEvent.GetEnd().AsTime().String(), time.UTC)
	if err != nil {
		s.server.Logger.Error(err)
		return nil, err
	}

	event := storage.Event{
		Title:       googleEvent.GetTitle(),
		UserID:      googleEvent.GetUserID(),
		Description: googleEvent.GetDescription(),
		End:         end,
		Start:       start,
	}

	return nil, s.server.App.Update(ctx, event)
}

func (s *GRPCServer) List(ctx context.Context, googleEvent *pb.Event) (*pb.Events, error) {
	var list pb.Events
	events := s.server.App.List(ctx, googleEvent.GetUserID())
	for _, event := range events {
		even := &pb.Event{
			Id:          event.ID,
			Title:       event.Title,
			UserID:      event.UserID,
			Description: event.Description,
			End:         timestamppb.New(event.End),
			Start:       timestamppb.New(event.Start),
		}

		list.Events = append(list.Events, even)
	}

	return &list, nil
}
