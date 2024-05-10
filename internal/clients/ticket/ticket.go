package ticket

import (
	"context"
	"log/slog"
	"time"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	ticketv1 "github.com/pvdiploma/diploma-protos/gen/go/ticket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api ticketv1.TicketServiceClient
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
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
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
		api: ticketv1.NewTicketServiceClient(cc),
		log: log,
	}, nil

}

func (c *Client) DeleteTicket(ctx context.Context, ticketID int64) (*ticketv1.DeleteTicketResponse, error) {
	id, err := c.api.DeleteTicket(context.Background(), &ticketv1.DeleteTicketRequest{
		Id: ticketID,
	})

	if err != nil {
		c.log.Error("failed to delete ticket: %v", err)
		return nil, err
	}
	return &ticketv1.DeleteTicketResponse{
		Id: id.Id, // id.GetId() ????
	}, nil
}

// TODO: put it in separate file
func InterceptorLogger(s *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, level grpclog.Level, msg string, fields ...any) {
		s.Log(ctx, slog.Level(level), msg, fields...)
	})
}
