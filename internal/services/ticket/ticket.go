package ticket

import (
	"context"
	"errors"
	"log/slog"
	"tn/internal/domain/models"
)

type TicketStorage interface {
	SaveTicket(ctx context.Context, ticket models.Ticket) (int64, error)
	DeleteTicket(ctx context.Context, ticket models.Ticket) (int64, error)
	GetTicket(ctx context.Context, ticketID int64) (models.Ticket, error)
	GetTicketByEmail(ctx context.Context, email string) (models.Ticket, error)
}

var ErrInvalidTicketID = errors.New("invalid ticketID")

type TicketService struct {
	log           *slog.Logger
	TicketStorage *TicketStorage
}

func New(
	log *slog.Logger,
	ticketStorage *TicketStorage,
) *TicketService {
	return &TicketService{
		log:           log,
		TicketStorage: ticketStorage,
	}
}

func (s *TicketService) AddTicket(ctx context.Context, event_category_id int64, name string, surname string, patronymic string, email string) (int64, error) {
	//
}

func (s *TicketService) GetTicket(ctx context.Context, ticketID int64) (models.Ticket, error) {
	// просто взять билет
}

func (s *TicketService) DeleteTicket(ctx context.Context, ticketID int64) (int64, error) {
	// удаление
}

func (s *TicketService) ActivateTicket(ctx context.Context, ticketID int64) (int64, error) {
	// взять бмлет и пометить как активированный
}

func (s *TicketService) IsActivated(ctx context.Context, ticketID int64) (bool, error) {
	// get ticket и проверка на активированность
}
