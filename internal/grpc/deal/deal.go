package dealgrpc

import (
	"context"
	"errors"
	"tn/internal/domain/models"
	"tn/internal/storage"
	"tn/internal/utils/converter"

	dealv1 "github.com/pvdiploma/diploma-protos/gen/go/deal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Deal interface {
	OfferDeal(ctx context.Context, deal models.Deal) (int64, error)
	AcceptDeal(ctx context.Context, dealID int64) (int64, error)
	RejectDeal(ctx context.Context, dealID int64) (int64, error)
	GetSentDeals(ctx context.Context, senderID int64) ([]models.Deal, error)
	GetProposedDeals(ctx context.Context, recipientID int64) ([]models.Deal, error)
	GetDealsByStatus(ctx context.Context, userID int64, status models.DealStatus) ([]models.Deal, error)
	GetDeal(ctx context.Context, dealID int64) (models.Deal, error)
	GetDealWidget(ctx context.Context, widgetID int64) (models.Widget, error)
}

type serverAPI struct {
	dealv1.UnimplementedDealsServiceServer
	deal Deal
}

func Register(gRPC *grpc.Server, deal Deal) {
	dealv1.RegisterDealsServiceServer(gRPC, &serverAPI{deal: deal})
}

func (s *serverAPI) OfferDeal(ctx context.Context, req *dealv1.OfferDealRequest) (*dealv1.OfferDealResponse, error) {

	id, err := s.deal.OfferDeal(ctx, models.Deal{
		SenderID:      req.GetSenderId(),
		RecipientID:   req.GetRecipientId(),
		OrganizerID:   req.GetOrganizerId(),
		DistributorID: req.GetDistributorId(),
		Commission:    req.GetCommission(),
		EventID:       req.GetEventId(),
	})
	if err != nil {
		if errors.Is(err, storage.ErrDealExists) {
			return nil, status.Error(codes.AlreadyExists, "deal already exists")
		}
		return nil, status.Error(codes.Internal, "failed to offer deal")
	}

	return &dealv1.OfferDealResponse{
		Id: id,
	}, nil
}

func (s *serverAPI) AcceptDeal(ctx context.Context, req *dealv1.AcceptDealRequest) (*dealv1.AcceptDealResponse, error) {

	id, err := s.deal.AcceptDeal(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, storage.ErrDealNotFound) {
			return nil, status.Error(codes.NotFound, "deal not found")
		}
		return nil, status.Error(codes.Internal, "failed to accept deal")
	}

	return &dealv1.AcceptDealResponse{
		Id: id,
	}, nil
}

func (s *serverAPI) RejectDeal(ctx context.Context, req *dealv1.RejectDealRequest) (*dealv1.RejectDealResponse, error) {

	id, err := s.deal.RejectDeal(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, storage.ErrDealNotFound) {
			return nil, status.Error(codes.NotFound, "deal not found")
		}
		return nil, status.Error(codes.Internal, "failed to reject deal")
	}

	return &dealv1.RejectDealResponse{
		Id: id,
	}, nil

}

func (s *serverAPI) GetSentDeals(ctx context.Context, req *dealv1.GetSentDealsRequest) (*dealv1.GetSentDealsResponse, error) {

	deals, err := s.deal.GetSentDeals(ctx, req.GetSenderId())

	if err != nil {
		if errors.Is(err, storage.ErrDealNotFound) {
			return nil, status.Error(codes.NotFound, "deals not found")
		}
		return nil, status.Error(codes.Internal, "failed to get sent deals")
	}

	return &dealv1.GetSentDealsResponse{
		Deals: converter.ModelDealsToProto(deals),
	}, nil
}

func (s *serverAPI) GetProposedDeals(ctx context.Context, req *dealv1.GetProposedDealsRequest) (*dealv1.GetProposedDealsResponse, error) {

	deals, err := s.deal.GetProposedDeals(ctx, req.GetRecipientId())

	if err != nil {
		if errors.Is(err, storage.ErrDealNotFound) {
			return nil, status.Error(codes.NotFound, "deals not found")
		}
		return nil, status.Error(codes.Internal, "failed to get sent deals")
	}

	return &dealv1.GetProposedDealsResponse{
		Deals: converter.ModelDealsToProto(deals),
	}, nil
}

func (s *serverAPI) GetDealsByStatus(ctx context.Context, req *dealv1.GetDealsByStatusRequest) (*dealv1.GetDealsByStatusResponse, error) {
	deals, err := s.deal.GetDealsByStatus(ctx, req.GetUserId(), converter.ProtoDealStatusToModels(req.GetStatus()))
	if err != nil {
		if errors.Is(err, storage.ErrDealNotFound) {
			return nil, status.Error(codes.NotFound, "deals not found")
		}
		return nil, status.Error(codes.Internal, "failed to get sent deals")
	}

	return &dealv1.GetDealsByStatusResponse{
		Deals: converter.ModelDealsToProto(deals),
	}, nil

}

func (s *serverAPI) GetDeal(ctx context.Context, req *dealv1.GetDealRequest) (*dealv1.GetDealResponse, error) {

	deal, err := s.deal.GetDeal(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, storage.ErrDealNotFound) {
			return nil, status.Error(codes.NotFound, "deal not found")
		}
		return nil, status.Error(codes.Internal, "failed to get deal")
	}

	return &dealv1.GetDealResponse{
		Deal: converter.ModelDealToProto(deal),
	}, nil
}

func (s *serverAPI) GetDealWidget(ctx context.Context, req *dealv1.GetDealWidgetRequest) (*dealv1.GetDealWidgetResponse, error) {

	widget, err := s.deal.GetDealWidget(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, storage.ErrDealWidgetNotFound) {
			return nil, status.Error(codes.NotFound, "deal widget not found")
		}
		return nil, status.Error(codes.Internal, "failed to get deal widget")
	}

	return &dealv1.GetDealWidgetResponse{
		Widget: converter.ModelWidgetToProto(widget),
	}, nil
}
