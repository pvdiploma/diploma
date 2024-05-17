package payment

import (
	"context"
	"log/slog"
	eventclient "tn/internal/clients/event"
	ticketclient "tn/internal/clients/ticket"
	"tn/internal/domain/models"
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

func (s *PaymentService) CheckAbilityToBuy(ctx context.Context, purchaseTickets []models.PurchaseTickets) (bool, error) {

	for _, category := range purchaseTickets {
		eventCategory, err := s.EventClient.GetEventCategory(ctx, category.ID)
		if err != nil {
			s.log.Error("Failed to get event category", err)
			return false, err
		}
		if eventCategory.Amount < uint32(category.Amount) {
			s.log.Info("Not enough tickets", slog.Any("category", category))
			return false, nil
		}
	}

	// (later) Check bank ability to buy
	return true, nil
}

func (s *PaymentService) CreateTickets(ctx context.Context, purchaseTickets models.PurchaseInfo, purchaseToken string) ([]int64, error) {
	var ticketsID []int64
	// test eventID
	eventID := int64(1)
	// forgot about discount
	discount := uint32(0)
	for _, t := range purchaseTickets.PurchaseTickets {
		ticket, err := s.TicketClient.CreateTicket(ctx,
			purchaseToken,
			eventID,
			t.ID,
			purchaseTickets.Name,
			purchaseTickets.Surname,
			purchaseTickets.Patronymic,
			discount,
			purchaseTickets.Email,
		)
		if err != nil {
			s.log.Error("Failed to create tickets", err)
			return []int64{}, err
		}

		ticketsID = append(ticketsID, ticket.Id)
	}

	return ticketsID, nil
}
