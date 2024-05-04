package event

import (
	"context"
	"errors"
	"log/slog"
	"tn/internal/domain/models"
)

// work with database layer
type EventStorage interface {
	SaveEvent(ctx context.Context, event models.Event) (int64, error)
	SaveEventCategory(ctx context.Context, event models.EventCategory) (int64, error)

	UpdateEvent(ctx context.Context, event models.Event) (int64, error)
	UpdateEventCategory(ctx context.Context, event models.EventCategory) (int64, error)

	DeleteEvent(ctx context.Context, eventID int64) (int64, error)
	DeleteEventCategory(ctx context.Context, eventID int64) (int64, error)

	GetEvent(ctx context.Context, eventID int64) (models.Event, error)
	GetEventCategory(ctx context.Context, eventID int64) (models.EventCategory, error)

	GetAllEvents(ctx context.Context) ([]models.Event, error)
	GetAllEventsCategory(ctx context.Context, eventID int64) ([]models.EventCategory, error)

	GetPrevEvents(ctx context.Context) ([]models.Event, error)
	GetPrevEventsCategory(ctx context.Context, eventID int64) ([]models.EventCategory, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidEventID     = errors.New("invalid eventID")
)

type EventService struct {
	log          *slog.Logger
	EventStorage EventStorage
}

func New(log *slog.Logger,
	eventStorage EventStorage,
) *EventService {

	return &EventService{
		log:          log,
		EventStorage: eventStorage,
	}
}

func (s *EventService) AddEvent(ctx context.Context, event models.Event) (int64, error) {
	// logic
	return s.EventStorage.SaveEvent(ctx, event)
}

func (s *EventService) UpdateEvent(ctx context.Context, event models.Event) (int64, error) {
	// logic
	return s.EventStorage.UpdateEvent(ctx, event)
}

func (s *EventService) DeleteEvent(ctx context.Context, eventID int64) (int64, error) {
	// logic
	return s.EventStorage.DeleteEvent(ctx, eventID)
}

func (s *EventService) GetEvent(ctx context.Context, eventID int64) (models.Event, error) {

	return s.EventStorage.GetEvent(ctx, eventID)

}

func (s *EventService) GetAllEvents(ctx context.Context) ([]models.Event, error) {

	return s.EventStorage.GetAllEvents(ctx)
}

func (s *EventService) GetPrevEvents(ctx context.Context) ([]models.Event, error) {

	return s.EventStorage.GetPrevEvents(ctx)
}
