package deal

import (
	"context"
	"errors"
	"log/slog"
	"tn/internal/domain/models"
	"tn/internal/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DealStorage interface {
	CreateDeal(ctx context.Context, deal models.Deal) (int64, error)
	UpdateDealStatus(ctx context.Context, dealID int64, status models.DealStatus) (int64, error)
	GetSentDeals(ctx context.Context, userID int64) ([]models.Deal, error)
	GetProposedDeals(ctx context.Context, userID int64) ([]models.Deal, error)
	GetDealsByStatus(ctx context.Context, userID int64, status models.DealStatus) ([]models.Deal, error)
}

type WidgetStorage interface {
}

var (
	ErrInvalidDealID = errors.New("invalid dealID")
)

type DealService struct {
	log           *slog.Logger
	DealStorage   DealStorage
	WidgetStorage WidgetStorage
}

func New(log *slog.Logger,
	dealStorage DealStorage,
	widgetStorage WidgetStorage,

) *DealService {
	return &DealService{
		log:           log,
		DealStorage:   dealStorage,
		WidgetStorage: widgetStorage,
	}
}

func (s *DealService) OfferDeal(ctx context.Context, deal models.Deal) (int64, error) {
	id, err := s.DealStorage.CreateDeal(ctx, deal)
	if err != nil {
		if errors.Is(err, storage.ErrDealExists) {
			return -1, status.Error(codes.AlreadyExists, "deal already exists")
		}
		return -1, status.Error(codes.Internal, "failed to offer deal")
	}

	return id, nil
}

func (s *DealService) AcceptDeal(ctx context.Context, dealID int64) (int64, error) {
	id, err := s.DealStorage.UpdateDealStatus(ctx, dealID, models.Accepted)
	if err != nil {
		if errors.Is(err, storage.ErrDealNotFound) {
			return -1, status.Error(codes.NotFound, "deal not found")
		}
		return -1, status.Error(codes.Internal, "failed to accept deal")
	}

	return id, nil
}

func (s *DealService) RejectDeal(ctx context.Context, dealID int64) (int64, error) {
	id, err := s.DealStorage.UpdateDealStatus(ctx, dealID, models.Rejected)
	if err != nil {
		if errors.Is(err, storage.ErrDealNotFound) {
			return -1, status.Error(codes.NotFound, "deal not found")
		}
		return -1, status.Error(codes.Internal, "failed to reject deal")
	}

	return id, nil
}

func (s *DealService) GetSentDeals(ctx context.Context, senderID int64) ([]models.Deal, error) {
	deals, err := s.DealStorage.GetSentDeals(ctx, senderID)
	if err != nil {
		if errors.Is(err, storage.ErrDealNotFound) {
			return deals, status.Error(codes.NotFound, "deals not found")
		}
		return deals, status.Error(codes.Internal, "failed to get sent deals")
	}

	return deals, nil
}

func (s *DealService) GetProposedDeals(ctx context.Context, recipientID int64) ([]models.Deal, error) {
	deals, err := s.DealStorage.GetProposedDeals(ctx, recipientID)
	if err != nil {
		if errors.Is(err, storage.ErrDealNotFound) {
			return deals, status.Error(codes.NotFound, "deals not found")
		}
		return deals, status.Error(codes.Internal, "failed to get proposed deals")
	}

	return deals, nil
}

func (s *DealService) GetDealsByStatus(ctx context.Context, userID int64, dealStatus models.DealStatus) ([]models.Deal, error) {
	deals, err := s.DealStorage.GetDealsByStatus(ctx, userID, dealStatus)
	if err != nil {
		if errors.Is(err, storage.ErrDealNotFound) {
			return deals, status.Error(codes.NotFound, "deals not found")
		}
		return deals, status.Error(codes.Internal, "failed to get deals by status")
	}

	return deals, nil
}
