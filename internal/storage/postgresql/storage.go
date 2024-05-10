package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"tn/internal/domain/models"
	"tn/internal/storage"

	"gorm.io/driver/postgres"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(storagePath string) (*Storage, error) {

	dbSQL, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: dbSQL,
	}), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) DB() *gorm.DB {
	return s.db
}

func (s *Storage) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (s *Storage) SaveUser(ctx context.Context, login string, email string, pwdHash []byte, role int32) (int64, error) {
	user := models.User{
		Login:   login,
		Email:   email,
		PwdHash: pwdHash,
		Role:    role,
	}

	result := s.db.WithContext(ctx).Create(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return -1, storage.ErrUserExists
		}
		return -1, result.Error
	}
	return user.Id, nil

}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {

	var user models.User
	result := s.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return user, storage.ErrUserNotFound
		}
		return user, result.Error
	}
	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	user, err := s.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}

	return user.Role == 0, nil
}

func (s *Storage) IsOrginiser(ctx context.Context, userID int64) (bool, error) {
	user, err := s.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}

	return user.Role == 1, nil
}

func (s *Storage) IsDistributor(ctx context.Context, userID int64) (bool, error) {
	user, err := s.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}
	return user.Role == 2, nil
}

func (s *Storage) IsBuyer(ctx context.Context, userID int64) (bool, error) {
	user, err := s.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}
	return user.Role == 3, nil
}

func (s *Storage) FindByID(ctx context.Context, userID int64) (models.User, error) {
	var user models.User
	result := s.db.WithContext(ctx).Where("id = ?", userID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return user, storage.ErrUserNotFound
		}
		return user, result.Error
	}
	return user, nil
}

func (s *Storage) Update(ctx context.Context, user models.User) error {
	result := s.db.WithContext(ctx).Save(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return storage.ErrUserNotFound
		}
		return result.Error
	}
	return nil
}

func (s *Storage) App(ctx context.Context, appID int32) (models.App, error) {
	var app models.App
	result := s.db.WithContext(ctx).Where("id = ?", appID).First(&app)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return app, storage.ErrAppNotFound
		}
		return app, result.Error
	}
	return app, nil
}

func (s *Storage) SaveEvent(ctx context.Context, event models.Event) (int64, error) {
	result := s.db.WithContext(ctx).Omit("categories").Create(&event)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return -1, storage.ErrEventExists
		}
		return -1, result.Error
	}
	return event.ID, nil
}

func (s *Storage) SaveEventCategory(ctx context.Context, event models.EventCategory) (int64, error) {
	result := s.db.WithContext(ctx).Create(&event)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return -1, storage.ErrEventExists
		}
		return -1, result.Error
	}
	return event.ID, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event models.Event, omits []string) (int64, error) {
	result := s.db.WithContext(ctx).Where("id = ?", event.ID).Omit(omits...).Updates(&event)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return -1, storage.ErrEventNotFound
		}
		return -1, result.Error
	}
	return event.ID, nil
}

func (s *Storage) UpdateEventCategory(ctx context.Context, EventCategory models.EventCategory, omits ...string) (int64, error) {

	result := s.db.WithContext(ctx).Where("id = ?", EventCategory.ID).Omit(omits...).Updates(&EventCategory)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return -1, storage.ErrEventNotFound
		}
		return -1, result.Error
	}
	return EventCategory.ID, nil
}

func (s *Storage) DeleteEvent(ctx context.Context, eventID int64) (int64, error) {
	result := s.db.WithContext(ctx).Delete(&models.Event{}, eventID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return -1, storage.ErrEventNotFound
		}
		return -1, result.Error
	}
	return eventID, nil
}

func (s *Storage) DeleteEventCategory(ctx context.Context, eventID int64) (int64, error) {

	result := s.db.WithContext(ctx).Delete(&models.EventCategory{}, eventID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return -1, storage.ErrEventNotFound
		}
		return -1, result.Error
	}
	return eventID, nil

}

func (s *Storage) GetEvent(ctx context.Context, eventID int64) (models.Event, error) {
	var event models.Event
	result := s.db.WithContext(ctx).Where("id = ?", eventID).First(&event)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return event, storage.ErrEventNotFound
		}
		return event, result.Error
	}
	return event, nil
}

func (s *Storage) GetEventCategory(ctx context.Context, eventID int64) ([]models.EventCategory, error) {

	var eventCategory []models.EventCategory
	result := s.db.WithContext(ctx).Where("event_id = ?", eventID).Find(&eventCategory)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return eventCategory, storage.ErrEventNotFound
		}
		return eventCategory, result.Error
	}
	return eventCategory, nil
}

func (s *Storage) GetAllEvents(ctx context.Context) ([]models.Event, error) {

	var events []models.Event
	result := s.db.WithContext(ctx).Find(&events)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return events, storage.ErrEventNotFound
		}
		return events, result.Error
	}
	return events, nil
}

func (s *Storage) GetPrevEvents(ctx context.Context) ([]models.Event, error) {
	now := time.Now()
	var events []models.Event
	result := s.db.WithContext(ctx).Where("date < ?", now).Find(events)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, storage.ErrEventNotFound
		}
		return nil, result.Error
	}
	return events, nil
}

func (s *Storage) SaveTicket(ctx context.Context, ticket models.Ticket) (int64, error) {
	result := s.db.WithContext(ctx).Create(&ticket)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return -1, storage.ErrTicketExists
		}
		return -1, result.Error
	}
	return ticket.ID, nil
}

func (s *Storage) DeleteTicket(ctx context.Context, ticketID int64) (int64, error) {

	result := s.db.WithContext(ctx).Delete(&models.Ticket{}, ticketID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return -1, storage.ErrTicketNotFound
		}
		return -1, result.Error
	}
	return ticketID, nil
}

func (s *Storage) GetTicket(ctx context.Context, ticketID int64) (models.Ticket, error) {

	var ticket models.Ticket
	result := s.db.WithContext(ctx).Where("id = ?", ticketID).First(&ticket)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ticket, storage.ErrTicketNotFound
		}
		return ticket, result.Error
	}
	return ticket, nil
}

func (s *Storage) GetTicketByEmail(ctx context.Context, email string) (models.Ticket, error) {

	var ticket models.Ticket
	result := s.db.WithContext(ctx).Where("email = ?", email).First(&ticket)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ticket, storage.ErrTicketNotFound
		}
		return ticket, result.Error
	}
	return ticket, nil
}

func (s *Storage) ActivateTicket(ctx context.Context, ticketId int64) (int64, error) {
	result := s.db.WithContext(ctx).Model(&models.Ticket{}).Where("id = ?", ticketId).Update("is_activated", true)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return -1, storage.ErrTicketNotFound
		}
		return -1, result.Error
	}
	return ticketId, nil
}
