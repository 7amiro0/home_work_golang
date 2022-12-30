package internalgrpc

import (
	"context"
	"fmt"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/server"
	pb "github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/server/grpc/google"
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

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

func getListEvents(events storage.SliceEvents) *pb.Events {
	var list []*pb.Event

	for _, event := range events.Events {
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

		list = append(list, even)
	}

	return &pb.Events{Events: list}
}

func (s *GRPCServer) ListUpcoming(ctx context.Context, until *pb.Until) (*pb.Events, error) {
	var duration time.Duration

	switch until.GetUntil() {
	case 0:
		duration = storage.Day
	case 1:
		duration = storage.Week
	case 2:
		duration = storage.Month
	}

	events, err := s.server.App.ListUpcoming(ctx, until.GetEvent().GetUser().GetName(), duration)
	if err != nil {
		s.server.Logger.Error("[ERR] While listing storage: ", err)
		return nil, err
	}

	return getListEvents(events), nil
}

func (s *GRPCServer) List(ctx context.Context, pbEvent *pb.Event) (*pb.Events, error) {
	events, err := s.server.App.List(ctx, pbEvent.GetUser().GetName())
	if err != nil {
		s.server.Logger.Error("[ERR] While listing storage: ", err)
		return nil, err
	}

	return getListEvents(events), nil
}
