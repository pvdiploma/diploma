package event

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"tn/internal/domain/models"

	sl "tn/pkg/logger"

	"gorm.io/gorm"
)

// work with database layer
type EventStorage interface {
	SaveEvent(ctx context.Context, event models.Event, tx *gorm.DB) (int64, error)
	SaveEventCategory(ctx context.Context, event models.EventCategory, tx *gorm.DB) (int64, error)

	UpdateEvent(ctx context.Context, tx *gorm.DB, event models.Event, omits []string) (int64, error)
	UpdateEventCategory(ctx context.Context, tx *gorm.DB, event models.EventCategory, omits ...string) (int64, error)

	DeleteEvent(ctx context.Context, tx *gorm.DB, eventID int64) (int64, error)
	DeleteEventCategory(ctx context.Context, tx *gorm.DB, eventID int64) (int64, error)

	GetEvent(ctx context.Context, eventID int64) (models.Event, error)
	GetEventCategory(ctx context.Context, eventID int64) ([]models.EventCategory, error)
	GetCategory(ctx context.Context, eventCategoryID int64) (models.EventCategory, error)
	GetEventIDByCategoryID(ctx context.Context, eventCategoryID int64) (int64, error)

	GetAllEvents(ctx context.Context) ([]models.Event, error)
	// GetAllEventsCategory(ctx context.Context, eventID int64) ([]models.EventCategory, error) do i need this?

	GetPrevEvents(ctx context.Context) ([]models.Event, error)
	// GetPrevEventsCategory(ctx context.Context, eventID int64) ([]models.EventCategory, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidEventID     = errors.New("invalid eventID")
)

type EventService struct {
	log          *slog.Logger
	EventStorage EventStorage
	db           *gorm.DB
}

func New(log *slog.Logger,
	eventStorage EventStorage,
	db *gorm.DB,
) *EventService {
	return &EventService{
		log:          log,
		EventStorage: eventStorage,
		db:           db,
	}
}

func (s *EventService) AddEvent(ctx context.Context, event models.Event) (int64, error) {

	tx := s.db.Begin()
	if tx.Error != nil {
		s.log.Error("Failed to create transaction", sl.Err(tx.Error))
		return -1, tx.Error
	}
	eventID, err := s.EventStorage.SaveEvent(ctx, event, tx)

	if err != nil {
		tx.Rollback()
		s.log.Error("Failed to save event", sl.Err(err))
		return -1, err
	}

	for _, category := range event.Categories {

		s.log.Info("Saving event category", slog.Int64("eventID", eventID), slog.Any("category", category))
		category.EventID = eventID
		_, err = s.EventStorage.SaveEventCategory(ctx, category, tx)

		if err != nil {
			tx.Rollback()
			s.log.Error("Failed to save event category", sl.Err(err))
			return -1, err
		}

	}

	err = tx.Commit().Error

	if err != nil {
		s.log.Error("Failed to commit transaction", sl.Err(err))
		return -1, err
	}

	return eventID, nil
}

func (s *EventService) UpdateEvent(ctx context.Context, event models.Event) (int64, error) {

	tx := s.db.Begin()
	if tx.Error != nil {
		s.log.Error("Failed to create transaction", sl.Err(tx.Error))
		return -1, tx.Error
	}

	eventOmits := GetEventOmitFields(event)
	_, err := s.EventStorage.UpdateEvent(ctx, tx, event, eventOmits)
	if err != nil {
		tx.Rollback()
		s.log.Error("Failed to update event", sl.Err(err))
		return -1, err
	}

	for i := range event.Categories {
		categoryOmits := GetEventCategoryOmitFields(event.Categories[i])
		_, err := s.EventStorage.UpdateEventCategory(ctx, tx, event.Categories[i], categoryOmits...)
		if err != nil {
			tx.Rollback()
			s.log.Error("Failed to update event category", sl.Err(err))
			return -1, err
		}
	}
	err = tx.Commit().Error

	if err != nil {
		s.log.Error("Failed to commit transaction", sl.Err(err))
		return -1, err
	}

	return event.ID, nil

}

func (s *EventService) DeleteEvent(ctx context.Context, eventID int64) (int64, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		s.log.Error("Failed to create transaction", sl.Err(tx.Error))
		return -1, tx.Error
	}

	_, err := s.EventStorage.DeleteEventCategory(ctx, tx, eventID)

	if err != nil {
		tx.Rollback()
		s.log.Error("Failed to delete event", sl.Err(err))
		return -1, err
	}

	_, err = s.EventStorage.DeleteEvent(ctx, tx, eventID)
	if err != nil {
		tx.Rollback()
		s.log.Error("Failed to delete event", sl.Err(err))
		return -1, err
	}

	err = tx.Commit().Error
	if err != nil {
		s.log.Error("Failed to commit transaction", sl.Err(err))
		return -1, err
	}

	return eventID, nil
}

func (s *EventService) GetEvent(ctx context.Context, eventID int64) (models.Event, error) {

	event, err := s.EventStorage.GetEvent(ctx, eventID)
	if err != nil {
		s.log.Error("Failed to get event", sl.Err(err))
	}

	event.Categories, err = s.EventStorage.GetEventCategory(ctx, eventID)
	fmt.Println(event.Categories)
	if err != nil {
		s.log.Error("Failed to get event category", sl.Err(err))
		return event, err
	}

	return event, nil
}

func (s *EventService) GetEventByCategoryId(ctx context.Context, eventCategoryID int64) (models.Event, error) {

	eventID, err := s.EventStorage.GetEventIDByCategoryID(ctx, eventCategoryID)
	if err != nil {
		s.log.Error("Failed to get event id", sl.Err(err))
		return models.Event{}, err
	}
	event, err := s.GetEvent(ctx, eventID)
	if err != nil {
		s.log.Error("Failed to get event by id", sl.Err(err))
		return models.Event{}, err
	}
	return event, nil
}

func (s *EventService) GetEventCategory(ctx context.Context, eventCategoryID int64) (models.EventCategory, error) {
	eventCategory, err := s.EventStorage.GetCategory(ctx, eventCategoryID)
	if err != nil {
		s.log.Error("Failed to get event category", sl.Err(err))
		return models.EventCategory{}, err
	}
	return eventCategory, nil
}

func (s *EventService) GetAllEvents(ctx context.Context) ([]models.Event, error) {

	events, err := s.EventStorage.GetAllEvents(ctx)
	if err != nil {
		s.log.Error("Failed to get all events", sl.Err(err))
	}

	//TODO: OPTIMIZE IT
	for i := range events {
		events[i].Categories, err = s.EventStorage.GetEventCategory(ctx, events[i].ID)
		if err != nil {
			s.log.Error("Failed to get event category", sl.Err(err))
			return events, err
		}
	}
	return events, nil
}

func (s *EventService) GetPrevEvents(ctx context.Context) ([]models.Event, error) {

	events, err := s.EventStorage.GetPrevEvents(ctx)
	if err != nil {
		s.log.Error("Failed to get prev events", sl.Err(err))
	}

	//TODO: OPTIMIZE IT
	for i := range events {
		events[i].Categories, err = s.EventStorage.GetEventCategory(ctx, events[i].ID)
		if err != nil {
			s.log.Error("Failed to get event category", sl.Err(err))
			return events, err
		}
	}
	return events, nil
}
