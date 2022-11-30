package internalgrpc

import (
	"context"
	"fmt"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/server"
	pb "github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/server/grpc/google"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage"
	grpcRun "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
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
		s.server.Logger.Error("[ERR] App don`t connect: ", err)
		return err
	}

	return s.httpServer.ListenAndServe()
}

func (s *GRPCServer) Stop(ctx context.Context) error {
	return s.server.App.Close(ctx)
}

func (s *GRPCServer) Add(ctx context.Context, pbEvent *pb.Event) (*emptypb.Empty, error) {
	start := pbEvent.GetStart().AsTime()
	end := pbEvent.GetEnd().AsTime()

	if !start.Before(end) {
		return nil, fmt.Errorf(server.ErrEndBeforeStart, end, start)
	}

	event := storage.Event{
		Title: pbEvent.GetTitle(),
		User: storage.User{
			Name: pbEvent.GetUser().GetName(),
			ID:   pbEvent.GetUser().GetId(),
		},
		Description: pbEvent.GetDescription(),
		Notify:      pbEvent.GetNotify(),
		End:         end,
		Start:       start,
	}

	return nil, s.server.App.Add(ctx, &event)
}

func (s *GRPCServer) Delete(ctx context.Context, pbEvent *pb.Event) (*emptypb.Empty, error) {
	return nil, s.server.App.Delete(ctx, pbEvent.GetId())
}

func (s *GRPCServer) Update(ctx context.Context, pbEvent *pb.Event) (*emptypb.Empty, error) {
	start := pbEvent.GetStart().AsTime()
	end := pbEvent.GetEnd().AsTime()

	if !start.Before(end) {
		return nil, fmt.Errorf(server.ErrEndBeforeStart, end, start)
	}

	event := storage.Event{
		ID:          pbEvent.GetId(),
		Title:       pbEvent.GetTitle(),
		Description: pbEvent.GetDescription(),
		Notify:      pbEvent.GetNotify(),
		End:         end,
		Start:       start,
	}

	return nil, s.server.App.Update(ctx, &event)
}

func (s *GRPCServer) List(ctx context.Context, pbEvent *pb.Event) (*pb.Events, error) {
	var list pb.Events
	events, err := s.server.App.List(ctx, pbEvent.GetUser().GetName())
	if err != nil {
		s.server.Logger.Error("[ERR] While list storage: ", err)
		return nil, err
	}

	for _, event := range events {
		even := &pb.Event{
			Id:    event.ID,
			Title: event.Title,
			User: &pb.User{
				Name: event.User.Name,
				Id:   event.User.ID,
			},
			Description: event.Description,
			Notify:      event.Notify,
			End:         timestamppb.New(event.End),
			Start:       timestamppb.New(event.Start),
		}

		list.Events = append(list.Events, even)
	}

	return &list, nil
}
