package eventgrpc

import (
	"context"
	"errors"
	"tn/internal/domain/models"
	"tn/internal/storage"
	"tn/internal/utils/converter"
	tokenmanager "tn/internal/utils/tokenManager"

	eventv1 "github.com/pvdiploma/diploma-protos/gen/go/event"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc/status"
)

type Event interface {
	AddEvent(ctx context.Context, event models.Event) (int64, error)
	UpdateEvent(ctx context.Context, event models.Event) (int64, error)
	DeleteEvent(ctx context.Context, eventID int64) (int64, error)
	GetEvent(ctx context.Context, eventID int64) (models.Event, error)
	GetAllEvents(ctx context.Context) ([]models.Event, error)
	GetPrevEvents(ctx context.Context) ([]models.Event, error)
}

type serverAPI struct {
	eventv1.UnimplementedEventServiceServer
	event Event
	tm    *tokenmanager.TokenManager
}

func Register(gRPC *grpc.Server, event Event, tm *tokenmanager.TokenManager) {
	eventv1.RegisterEventServiceServer(gRPC, &serverAPI{event: event, tm: tm})
}

// NOTE: mb it looks strange but i don't know how make it better
func (s *serverAPI) AuthMiddleware(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return -1, status.Error(codes.PermissionDenied, "invalid credentials")
	}

	token := md.Get("token")

	// refresh ...
	if fl, id := s.tm.IsOrganizer(token[0]); fl {
		return id, nil
	}

	return -1, status.Error(codes.PermissionDenied, "invalid credentials")
}

func (s *serverAPI) AddEvent(ctx context.Context, req *eventv1.AddEventRequest) (*eventv1.AddEventResponse, error) {

	ownerID, err := s.AuthMiddleware(ctx)
	if err != nil {
		return nil, err
	}

	id, err := s.event.AddEvent(ctx, models.Event{
		OwnerID:      ownerID,
		Name:         req.GetName(),
		Description:  req.GetDescription(),
		Country:      req.GetCountry(),
		City:         req.GetCity(),
		Place:        req.GetPlace(),
		Date:         req.GetDate().AsTime(),
		TicketAmount: req.GetTicketAmount(),
		Age:          req.GetAge(),
		Categories:   converter.ProtoCategoryToModels(req.GetCategories()),
	})

	if err != nil {
		if errors.Is(err, storage.ErrEventExists) {
			return nil, status.Error(codes.AlreadyExists, "event already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &eventv1.AddEventResponse{
		EventId: id,
	}, nil
}

func (s *serverAPI) UpdateEvent(ctx context.Context, req *eventv1.UpdateEventRequest) (*eventv1.UpdateEventResponse, error) {

	ownerID, err := s.AuthMiddleware(ctx)
	if err != nil {
		return nil, err
	}

	id, err := s.event.UpdateEvent(ctx, models.Event{
		OwnerID:      ownerID,
		ID:           req.GetEventId(),
		Name:         req.GetName(),
		Description:  req.GetDescription(),
		Country:      req.GetCountry(),
		City:         req.GetCity(),
		Place:        req.GetPlace(),
		Date:         req.GetDate().AsTime(),
		TicketAmount: req.GetTicketAmount(),
		Age:          req.GetAge(),
		Categories:   converter.ProtoCategoryToModels(req.GetCategories()),
	})

	if err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			return nil, status.Error(codes.NotFound, "event not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &eventv1.UpdateEventResponse{
		EventId: id,
	}, nil
}

func (s *serverAPI) DeleteEvent(ctx context.Context, req *eventv1.DeleteEventRequest) (*eventv1.DeleteEventResponse, error) {
	_, err := s.AuthMiddleware(ctx)
	if err != nil {
		return nil, err
	}

	id, err := s.event.DeleteEvent(ctx, req.GetEventId())
	if err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			return nil, status.Error(codes.NotFound, "event not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &eventv1.DeleteEventResponse{
		EventId: id,
	}, nil
}

func (s *serverAPI) GetEvent(ctx context.Context, req *eventv1.GetEventRequest) (*eventv1.GetEventResponse, error) {

	event, err := s.event.GetEvent(ctx, req.GetEventId())

	if err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			return nil, status.Error(codes.NotFound, "event not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &eventv1.GetEventResponse{
		Event: converter.ModelEventToProto(event),
	}, nil
}

func (s *serverAPI) GetAllEvents(ctx context.Context, req *eventv1.GetAllEventsRequest) (*eventv1.GetAllEventsResponse, error) {
	events, err := s.event.GetAllEvents(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			return nil, status.Error(codes.NotFound, "event not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	var protoEvents []*eventv1.Event
	for i := range events {
		protoEvents = append(protoEvents, converter.ModelEventToProto(events[i]))
	}
	return &eventv1.GetAllEventsResponse{
		Events: protoEvents,
	}, nil
}

func (s *serverAPI) GetPrevEvents(ctx context.Context, req *eventv1.GetPrevEventsRequest) (*eventv1.GetPrevEventsResponse, error) {
	events, err := s.event.GetPrevEvents(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			return nil, status.Error(codes.NotFound, "event not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	var protoEvents []*eventv1.Event
	for i := range events {
		protoEvents = append(protoEvents, converter.ModelEventToProto(events[i]))
	}
	return &eventv1.GetPrevEventsResponse{
		Events: protoEvents,
	}, nil

}
