package ticket

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	eventclient "tn/internal/clients/event"
	"tn/internal/domain/models"
	"tn/internal/storage"

	sl "tn/pkg/logger"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TicketStorage interface {
	SaveTicket(ctx context.Context, ticket models.Ticket) (int64, error)
	DeleteTicket(ctx context.Context, ticket int64) (int64, error)
	GetTicket(ctx context.Context, ticketID int64) (models.Ticket, error)
	GetTicketByEmail(ctx context.Context, email string) (models.Ticket, error)
}

var ErrInvalidTicketID = errors.New("invalid ticketID")

type TicketService struct {
	log           *slog.Logger
	TicketStorage TicketStorage
	EventClient   *eventclient.Client
}

func New(
	log *slog.Logger,
	ticketStorage TicketStorage,
	eventClient *eventclient.Client,
) *TicketService {
	return &TicketService{
		log:           log,
		TicketStorage: ticketStorage,
		EventClient:   eventClient,
	}
}

type Image struct {
	Image     []byte
	ImagePath string
}

func GenerateImage(ticket models.Ticket, event models.Event, eventCategory models.EventCategory) (Image, error) {
	return Image{
		Image:     []byte("test"),
		ImagePath: "test_path",
	}, nil
}

func (s *TicketService) AddTicket(ctx context.Context, eventCategoryID int64, name string, surname string, patronymic string, discount uint32, email string) (int64, error) {
	//НУЖНА ЛИ ТУТ ТРАНЗАКЦИЯ????

	event, err := s.EventClient.GetEventByCategoryId(ctx, eventCategoryID)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return -1, storage.ErrEventNotFound
		}
		return -1, err
	}

	eventCategory, err := exctractEventCategory(eventCategoryID, event)

	if err != nil {
		return -1, err
	}

	ticket := models.Ticket{
		EventCategoryID: eventCategoryID,
		Name:            name,
		Surname:         surname,
		Patronymic:      patronymic,
		Email:           email,
		Discount:        discount,
		Total:           eventCategory.Price - (eventCategory.Price * discount / 100),
		QRCode:          nil,
		IsActivated:     false,
	}
	img, err := GenerateImage(ticket, event, eventCategory)

	if err != nil {
		return -1, err
	}

	ticket.ImageBytes = img.Image
	ticket.ImagePath = img.ImagePath

	id, err := s.TicketStorage.SaveTicket(ctx, ticket)
	if err != nil {
		s.log.Error("Failed to save ticket", sl.Err(err))
		return -1, nil
	}
	return id, nil
}

func (s *TicketService) GetTicket(ctx context.Context, ticketID int64) (models.Ticket, error) {
	ticket, err := s.TicketStorage.GetTicket(ctx, ticketID)
	if err != nil {
		s.log.Error("Failed to get ticket", sl.Err(err))
		return models.Ticket{}, err
	}
	return ticket, nil
}

func (s *TicketService) DeleteTicket(ctx context.Context, ticketID int64) (int64, error) {
	id, err := s.TicketStorage.DeleteTicket(ctx, ticketID)
	if err != nil {
		s.log.Error("Failed to delete ticket", sl.Err(err))
		return -1, err
	}

	return id, nil
}

func (s *TicketService) ActivateTicket(ctx context.Context, ticketID int64) (int64, error) {
	// взять бмлет и пометить как активированный
	return -1, nil
}

func (s *TicketService) IsActivated(ctx context.Context, ticketID int64) (bool, error) {
	ticket, err := s.TicketStorage.GetTicket(ctx, ticketID)
	if err != nil {
		s.log.Error("Failed to get ticket", sl.Err(err))
		return false, err
	}

	return ticket.IsActivated, nil
}

func exctractEventCategory(eventCategoryId int64, event models.Event) (models.EventCategory, error) {

	for _, eventCategory := range event.Categories {
		if eventCategory.ID == eventCategoryId {
			return eventCategory, nil
		}
	}
	return models.EventCategory{}, fmt.Errorf("event category with id %d not found", eventCategoryId)
}
