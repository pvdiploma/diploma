package event

import (
	"context"
	"log/slog"
	"time"
	"tn/internal/domain/models"

	"tn/internal/utils/converter"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	dealv1 "github.com/pvdiploma/diploma-protos/gen/go/deal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api dealv1.DealsServiceClient
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
		api: dealv1.NewDealsServiceClient(cc),
		log: log,
	}, nil
}

func (c *Client) GetDeal(ctx context.Context, dealID int64) (models.Deal, error) {

	resp, err := c.api.GetDeal(ctx, &dealv1.GetDealRequest{Id: dealID})
	if err != nil {
		c.log.Error("Failed to get deal", err)
		return models.Deal{}, err
	}

	return converter.ProtoDealToModel(resp.Deal), nil
}

func (c *Client) GetDealWidget(ctx context.Context, widgetID int64) (models.Widget, error) {

	resp, err := c.api.GetDealWidget(ctx, &dealv1.GetDealWidgetRequest{Id: widgetID})
	if err != nil {
		c.log.Error("Failed to get widget", err)

		return models.Widget{}, err
	}

	return converter.ProtoWidgetToModel(resp.Widget), nil
}

func InterceptorLogger(s *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, level grpclog.Level, msg string, fields ...any) {
		s.Log(ctx, slog.Level(level), msg, fields...)
	})
}
