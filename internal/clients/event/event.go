package event

// Implementation of client for use in external services

import (
	"context"
	"log/slog"
	"time"
	"tn/internal/domain/models"

	sl "tn/pkg/logger"

	"tn/internal/utils/converter"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	eventv1 "github.com/pvdiploma/diploma-protos/gen/go/event"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api eventv1.EventServiceClient
	log *slog.Logger
}

func NewClient(
	ctx context.Context,
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded), //!!!  Be careful with that
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	cc, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)

	if err != nil {
		log.Error("failed to dial: %v", err)
		return nil, err
	}

	return &Client{
		api: eventv1.NewEventServiceClient(cc),
		log: log,
	}, nil
}

func (s *Client) GetEvent(ctx context.Context, eventID int64) (models.Event, error) {
	var event models.Event
	resp, err := s.api.GetEvent(ctx, &eventv1.GetEventRequest{EventId: eventID})
	if err != nil {
		s.log.Error("Failed to get event", sl.Err(err))
		return models.Event{}, err
	}
	event = converter.ProtoEventToModel(resp.Event) // check resp.GetEvent()
	return event, nil
}

func (s *Client) GetEventByCategoryId(ctx context.Context, eventCategoryID int64) (models.Event, error) {
	var event models.Event
	resp, err := s.api.GetEventByCategoryId(ctx, &eventv1.GetEventByCategoryIdRequest{EventCategoryId: eventCategoryID})
	if err != nil {
		s.log.Error("Failed to get event", sl.Err(err))
		return models.Event{}, err
	}
	event = converter.ProtoEventToModel(resp.Event)
	return event, nil
}

func InterceptorLogger(s *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, level grpclog.Level, msg string, fields ...any) {
		s.Log(ctx, slog.Level(level), msg, fields...)
	})
}
