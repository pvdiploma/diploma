package payment

import (
	"context"
	"log/slog"
	dealclient "tn/internal/clients/deal"
	eventclient "tn/internal/clients/event"
	ticketclient "tn/internal/clients/ticket"
	"tn/internal/domain/models"
)

type PaymentStorage interface {
	GetBalance(ctx context.Context, userID int64) (float64, error)
	UpdateBalance(ctx context.Context, userID int64, balance float64) (int64, error)
}

type PaymentService struct {
	log            *slog.Logger
	PaymentStorage PaymentStorage
	TicketClient   *ticketclient.Client
	EventClient    *eventclient.Client
	DealClient     *dealclient.Client
}

func New(log *slog.Logger, storage PaymentStorage, ticketClient *ticketclient.Client, eventClient *eventclient.Client, dealClient *dealclient.Client) *PaymentService {
	return &PaymentService{
		log:            log,
		PaymentStorage: storage,
		TicketClient:   ticketClient,
		EventClient:    eventClient,
		DealClient:     dealClient,
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

	widget, err := s.DealClient.GetDealWidget(ctx, purchaseTickets.WidgetID)
	if err != nil {
		s.log.Error("Failed to get deal widget", err)
		return []int64{}, err
	}

	deal, err := s.DealClient.GetDeal(ctx, widget.DealID)
	if err != nil {
		s.log.Error("Failed to get deal", err)
		return []int64{}, err
	}

	// forgot about discount
	discount := uint32(0)
	for _, t := range purchaseTickets.PurchaseTickets {
		ticket, err := s.TicketClient.CreateTicket(ctx,
			purchaseToken,
			deal.EventID,
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

func (s *PaymentService) SubmitPurchase(ctx context.Context, ticketsID []int64, widgetID int64) error {
	total := float64(0)

	for _, ticketID := range ticketsID {
		ticket, err := s.TicketClient.GetTicket(ctx, ticketID)
		if err != nil {
			s.log.Error("Failed to get ticket", err)
			return err
		}
		total += float64(ticket.Total)
	}

	widget, err := s.DealClient.GetDealWidget(ctx, widgetID)
	if err != nil {
		s.log.Error("Failed to get deal widget", err)
		return err
	}

	deal, err := s.DealClient.GetDeal(ctx, widget.DealID)
	if err != nil {
		s.log.Error("Failed to get deal", err)
		return err
	}

	orgBalance, err := s.PaymentStorage.GetBalance(ctx, deal.OrganizerID)
	if err != nil {
		s.log.Error("Failed to get balance", err)
		return err
	}

	distBalance, err := s.PaymentStorage.GetBalance(ctx, deal.DistributorID)
	if err != nil {
		s.log.Error("Failed to get balance", err)
		return err
	}

	orgBalance += total * (100 - float64(deal.Commission)) / 100
	distBalance += total * float64(deal.Commission) / 100

	_, err = s.PaymentStorage.UpdateBalance(ctx, deal.OrganizerID, orgBalance)
	if err != nil {
		s.log.Error("Failed to update balance", err)
		return err
	}

	_, err = s.PaymentStorage.UpdateBalance(ctx, deal.DistributorID, distBalance)
	if err != nil {
		s.log.Error("Failed to update balance", err)
		return err
	}
	return nil
}
