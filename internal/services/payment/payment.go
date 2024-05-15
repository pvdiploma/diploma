package payment

import (
	"context"
	"log/slog"
	eventclient "tn/internal/clients/event"
	ticketclient "tn/internal/clients/ticket"
)

type PaymentService struct {
	log          *slog.Logger
	TicketClient *ticketclient.Client
	EventClient  *eventclient.Client
}

func New(log *slog.Logger) *PaymentService {
	return &PaymentService{
		log: log,
	}
}

func (s *PaymentService) CheckAbilityToBuy(ctx context.Context) (bool, error) {
	//testing data
	test_categories := []int64{1, 2, 3, 4, 5}
	// it's mock
	for _, category := range test_categories {
		eventCategory, err := s.EventClient.GetEventCategory(ctx, category)
		if err != nil {
			s.log.Error("Failed to get event category", err)
			return false, err
		}
		if eventCategory.Amount < uint32(category) {
			s.log.Info("Not enough tickets", slog.Any("category", category))
			return false, nil
		}
	}

	// (later) Check bank ability to buy
	return true, nil
}
